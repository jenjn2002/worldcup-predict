package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	db        *sql.DB
	jwtSecret string
}

// ---------------------------------------------------------------
// Auth
// ---------------------------------------------------------------

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Username == "" || req.Email == "" || len(req.Password) < 6 {
		writeError(w, http.StatusBadRequest, "username, email are required and password must be at least 6 characters")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	var user User
	err = s.db.QueryRow(
		`INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)
		 RETURNING id, username, email, is_admin, created_at`,
		req.Username, req.Email, string(hash),
	).Scan(&user.ID, &user.Username, &user.Email, &user.IsAdmin, &user.CreatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			writeError(w, http.StatusConflict, "username or email already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	token, err := generateToken(s.jwtSecret, user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not generate token")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Username = strings.TrimSpace(req.Username)

	var user User
	var passwordHash string
	err := s.db.QueryRow(
		`SELECT id, username, email, password_hash, is_admin, created_at
		 FROM users WHERE username = $1 OR email = $1`,
		req.Username,
	).Scan(&user.ID, &user.Username, &user.Email, &passwordHash, &user.IsAdmin, &user.CreatedAt)

	if err == sql.ErrNoRows {
		writeError(w, http.StatusUnauthorized, "invalid username or password")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)) != nil {
		writeError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	token, err := generateToken(s.jwtSecret, user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not generate token")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r)
	var user User
	err := s.db.QueryRow(
		`SELECT id, username, email, is_admin, created_at FROM users WHERE id = $1`,
		claims.UserID,
	).Scan(&user.ID, &user.Username, &user.Email, &user.IsAdmin, &user.CreatedAt)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// ---------------------------------------------------------------
// Matches
// ---------------------------------------------------------------

func (s *Server) handleListMatches(w http.ResponseWriter, r *http.Request) {
	claims, _ := s.optionalClaims(r)

	rows, err := s.db.Query(`
		SELECT m.id, m.stage, m.group_name, m.home_team_id, m.away_team_id,
		       ht.name, at.name, m.match_time, m.home_score, m.away_score, m.status
		FROM matches m
		JOIN teams ht ON ht.id = m.home_team_id
		JOIN teams at ON at.id = m.away_team_id
		ORDER BY m.match_time ASC, m.id ASC
	`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load matches")
		return
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var m Match
		if err := rows.Scan(&m.ID, &m.Stage, &m.GroupName, &m.HomeTeamID, &m.AwayTeamID,
			&m.HomeTeam, &m.AwayTeam, &m.MatchTime, &m.HomeScore, &m.AwayScore, &m.Status); err != nil {
			writeError(w, http.StatusInternalServerError, "could not read matches")
			return
		}
		matches = append(matches, m)
	}

	// attach current user's predictions, if logged in
	if claims != nil && len(matches) > 0 {
		predRows, err := s.db.Query(
			`SELECT match_id, pred_home, pred_away, points FROM predictions WHERE user_id = $1`,
			claims.UserID,
		)
		if err == nil {
			defer predRows.Close()
			predMap := map[int][3]int{}
			for predRows.Next() {
				var matchID, ph, pa, pts int
				if err := predRows.Scan(&matchID, &ph, &pa, &pts); err == nil {
					predMap[matchID] = [3]int{ph, pa, pts}
				}
			}
			for i := range matches {
				if pred, ok := predMap[matches[i].ID]; ok {
					ph, pa, pts := pred[0], pred[1], pred[2]
					matches[i].MyPredHome = &ph
					matches[i].MyPredAway = &pa
					matches[i].MyPoints = &pts
				}
			}
		}
	}

	writeJSON(w, http.StatusOK, matches)
}

// ---------------------------------------------------------------
// Predictions
// ---------------------------------------------------------------

type predictionRequest struct {
	MatchID  int `json:"match_id"`
	PredHome int `json:"pred_home"`
	PredAway int `json:"pred_away"`
}

func (s *Server) handleSubmitPrediction(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r)

	var req predictionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.PredHome < 0 || req.PredAway < 0 {
		writeError(w, http.StatusBadRequest, "scores must be zero or positive")
		return
	}

	var matchTime time.Time
	var status string
	err := s.db.QueryRow(`SELECT match_time, status FROM matches WHERE id = $1`, req.MatchID).
		Scan(&matchTime, &status)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "match not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}

	if status != "scheduled" || time.Now().After(matchTime) {
		writeError(w, http.StatusBadRequest, "predictions are closed for this match (already started or finished)")
		return
	}

	_, err = s.db.Exec(`
		INSERT INTO predictions (user_id, match_id, pred_home, pred_away, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (user_id, match_id)
		DO UPDATE SET pred_home = EXCLUDED.pred_home, pred_away = EXCLUDED.pred_away, updated_at = NOW()
	`, claims.UserID, req.MatchID, req.PredHome, req.PredAway)

	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not save prediction")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleMyPredictions(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r)

	rows, err := s.db.Query(
		`SELECT id, user_id, match_id, pred_home, pred_away, points, updated_at
		 FROM predictions WHERE user_id = $1`,
		claims.UserID,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load predictions")
		return
	}
	defer rows.Close()

	preds := []Prediction{}
	for rows.Next() {
		var p Prediction
		if err := rows.Scan(&p.ID, &p.UserID, &p.MatchID, &p.PredHome, &p.PredAway, &p.Points, &p.UpdatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "could not read predictions")
			return
		}
		preds = append(preds, p)
	}
	writeJSON(w, http.StatusOK, preds)
}

// ---------------------------------------------------------------
// Leaderboard
// ---------------------------------------------------------------

func (s *Server) handleLeaderboard(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query(`
		SELECT u.id, u.username,
		       COALESCE(SUM(p.points), 0) AS total_points,
		       COALESCE(SUM(CASE WHEN p.points = 3 THEN 1 ELSE 0 END), 0) AS exact_count,
		       COUNT(p.id) AS predicted_count
		FROM users u
		LEFT JOIN predictions p ON p.user_id = u.id
		GROUP BY u.id, u.username
		ORDER BY total_points DESC, exact_count DESC, u.username ASC
	`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load leaderboard")
		return
	}
	defer rows.Close()

	entries := []LeaderboardEntry{}
	for rows.Next() {
		var e LeaderboardEntry
		if err := rows.Scan(&e.UserID, &e.Username, &e.TotalPoints, &e.ExactCount, &e.PredictedCount); err != nil {
			writeError(w, http.StatusInternalServerError, "could not read leaderboard")
			return
		}
		entries = append(entries, e)
	}
	writeJSON(w, http.StatusOK, entries)
}

// ---------------------------------------------------------------
// Admin
// ---------------------------------------------------------------

type matchResultRequest struct {
	HomeScore int `json:"home_score"`
	AwayScore int `json:"away_score"`
}

// handleSetMatchResult sets the final score for a match and scores all predictions:
//   - exact score match  -> 3 points
//   - correct outcome (win/draw/loss) -> 1 point
//   - otherwise -> 0 points
func (s *Server) handleSetMatchResult(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	matchID, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid match id")
		return
	}

	var req matchResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.HomeScore < 0 || req.AwayScore < 0 {
		writeError(w, http.StatusBadRequest, "scores must be zero or positive")
		return
	}

	tx, err := s.db.Begin()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		`UPDATE matches SET home_score = $1, away_score = $2, status = 'finished' WHERE id = $3`,
		req.HomeScore, req.AwayScore, matchID,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not update match")
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		writeError(w, http.StatusNotFound, "match not found")
		return
	}

	actualOutcome := outcome(req.HomeScore, req.AwayScore)

	rows, err := tx.Query(`SELECT id, pred_home, pred_away FROM predictions WHERE match_id = $1`, matchID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load predictions")
		return
	}

	type predRow struct {
		id, ph, pa int
	}
	var preds []predRow
	for rows.Next() {
		var p predRow
		if err := rows.Scan(&p.id, &p.ph, &p.pa); err != nil {
			rows.Close()
			writeError(w, http.StatusInternalServerError, "could not read predictions")
			return
		}
		preds = append(preds, p)
	}
	rows.Close()

	for _, p := range preds {
		points := 0
		if p.ph == req.HomeScore && p.pa == req.AwayScore {
			points = 3
		} else if outcome(p.ph, p.pa) == actualOutcome {
			points = 1
		}
		if _, err := tx.Exec(`UPDATE predictions SET points = $1 WHERE id = $2`, points, p.id); err != nil {
			writeError(w, http.StatusInternalServerError, "could not update prediction points")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		writeError(w, http.StatusInternalServerError, "could not save changes")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// outcome returns 1 for home win, 0 for draw, -1 for away win.
func outcome(home, away int) int {
	switch {
	case home > away:
		return 1
	case home < away:
		return -1
	default:
		return 0
	}
}

type adminMatchRequest struct {
	Stage      string `json:"stage"`
	GroupName  string `json:"group_name"`
	HomeTeamID int    `json:"home_team_id"`
	AwayTeamID int    `json:"away_team_id"`
	MatchTime  string `json:"match_time"` // RFC3339
}

func (s *Server) handleCreateMatch(w http.ResponseWriter, r *http.Request) {
	var req adminMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	matchTime, err := time.Parse(time.RFC3339, req.MatchTime)
	if err != nil {
		writeError(w, http.StatusBadRequest, "match_time must be in RFC3339 format, e.g. 2026-06-15T18:00:00Z")
		return
	}
	if req.Stage == "" {
		req.Stage = "group"
	}

	var id int
	err = s.db.QueryRow(
		`INSERT INTO matches (stage, group_name, home_team_id, away_team_id, match_time)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		req.Stage, req.GroupName, req.HomeTeamID, req.AwayTeamID, matchTime,
	).Scan(&id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create match")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]int{"id": id})
}

func (s *Server) handleListTeams(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query(`SELECT id, name, COALESCE(group_name, '') FROM teams ORDER BY group_name, name`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load teams")
		return
	}
	defer rows.Close()

	teams := []Team{}
	for rows.Next() {
		var t Team
		if err := rows.Scan(&t.ID, &t.Name, &t.GroupName); err != nil {
			writeError(w, http.StatusInternalServerError, "could not read teams")
			return
		}
		teams = append(teams, t)
	}
	writeJSON(w, http.StatusOK, teams)
}
