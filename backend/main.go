package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	databaseURL := getEnv("DATABASE_URL", "postgres://wcuser:wcpass@localhost:5432/worldcup?sslmode=disable")
	jwtSecret := getEnv("JWT_SECRET", "dev_secret_change_me")
	port := getEnv("PORT", "8080")

	db := connectDB(databaseURL)
	defer db.Close()

	server := &Server{db: db, jwtSecret: jwtSecret}

	ensureAdminUser(db)

	mux := http.NewServeMux()

	// Public
	mux.HandleFunc("POST /api/register", server.handleRegister)
	mux.HandleFunc("POST /api/login", server.handleLogin)
	mux.HandleFunc("GET /api/matches", server.handleListMatches)
	mux.HandleFunc("GET /api/teams", server.handleListTeams)
	mux.HandleFunc("GET /api/leaderboard", server.handleLeaderboard)

	// Authenticated
	mux.HandleFunc("GET /api/me", server.authMiddleware(server.handleMe))
	mux.HandleFunc("POST /api/predictions", server.authMiddleware(server.handleSubmitPrediction))
	mux.HandleFunc("GET /api/predictions/me", server.authMiddleware(server.handleMyPredictions))

	// Admin
	mux.HandleFunc("POST /api/admin/matches", server.adminMiddleware(server.handleCreateMatch))
	mux.HandleFunc("PUT /api/admin/matches/{id}/result", server.adminMiddleware(server.handleSetMatchResult))

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	handler := corsMiddleware(mux)

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// corsMiddleware allows the Next.js frontend (running on a different
// origin/port) to call the API.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ensureAdminUser creates a default admin account (from env vars) on first
// startup, so the operator can immediately set match results.
func ensureAdminUser(db *sql.DB) {
	username := getEnv("ADMIN_USERNAME", "admin")
	password := getEnv("ADMIN_PASSWORD", "")
	if password == "" {
		return
	}

	var exists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`, username).Scan(&exists)
	if err != nil {
		log.Printf("could not check for admin user: %v", err)
		return
	}
	if exists {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("could not hash admin password: %v", err)
		return
	}

	_, err = db.Exec(
		`INSERT INTO users (username, email, password_hash, is_admin) VALUES ($1, $2, $3, TRUE)`,
		username, username+"@worldcup.local", string(hash),
	)
	if err != nil {
		log.Printf("could not create admin user: %v", err)
		return
	}
	log.Printf("created default admin user %q", username)
}
