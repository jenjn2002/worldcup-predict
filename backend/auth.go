package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userContextKey contextKey = "user"

type AuthClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

func generateToken(secret string, user User) (string, error) {
	claims := AuthClaims{
		UserID:   user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func parseToken(secret, tokenString string) (*AuthClaims, error) {
	claims := &AuthClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		if err == nil {
			err = jwt.ErrTokenInvalidClaims
		}
		return nil, err
	}
	return claims, nil
}

// authMiddleware requires a valid JWT and injects claims into the request context.
func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := s.optionalClaims(r)
		if err != nil || claims == nil {
			writeError(w, http.StatusUnauthorized, "missing or invalid authorization token")
			return
		}
		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next(w, r.WithContext(ctx))
	}
}

// adminMiddleware requires a valid JWT belonging to an admin user.
func (s *Server) adminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return s.authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		claims := claimsFromContext(r)
		if claims == nil || !claims.IsAdmin {
			writeError(w, http.StatusForbidden, "admin privileges required")
			return
		}
		next(w, r)
	})
}

// optionalClaims parses the bearer token if present, but does not error if absent.
func (s *Server) optionalClaims(r *http.Request) (*AuthClaims, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return nil, nil
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, jwt.ErrTokenMalformed
	}
	return parseToken(s.jwtSecret, strings.TrimSpace(parts[1]))
}

func claimsFromContext(r *http.Request) *AuthClaims {
	claims, ok := r.Context().Value(userContextKey).(*AuthClaims)
	if !ok {
		return nil
	}
	return claims
}
