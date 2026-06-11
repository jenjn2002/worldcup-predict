package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

type Team struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	GroupName string `json:"group_name"`
}

type Match struct {
	ID         int       `json:"id"`
	Stage      string    `json:"stage"`
	GroupName  string    `json:"group_name"`
	HomeTeamID int       `json:"home_team_id"`
	AwayTeamID int       `json:"away_team_id"`
	HomeTeam   string    `json:"home_team"`
	AwayTeam   string    `json:"away_team"`
	MatchTime  time.Time `json:"match_time"`
	HomeScore  *int      `json:"home_score"`
	AwayScore  *int      `json:"away_score"`
	Status     string    `json:"status"`

	// Current user's prediction for this match (if authenticated & exists)
	MyPredHome *int `json:"my_pred_home,omitempty"`
	MyPredAway *int `json:"my_pred_away,omitempty"`
	MyPoints   *int `json:"my_points,omitempty"`
}

type Prediction struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	MatchID   int       `json:"match_id"`
	PredHome  int       `json:"pred_home"`
	PredAway  int       `json:"pred_away"`
	Points    int       `json:"points"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LeaderboardEntry struct {
	UserID         int    `json:"user_id"`
	Username       string `json:"username"`
	TotalPoints    int    `json:"total_points"`
	ExactCount     int    `json:"exact_count"`
	PredictedCount int    `json:"predicted_count"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
