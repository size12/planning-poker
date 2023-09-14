package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func generateCookie() (*http.Cookie, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:  "poker_player_id",
		Value: id.String(),
		Path:  "/",
	}

	return cookie, nil
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("poker_player_id")
		if err != nil {
			cookie, err = generateCookie()
			if err != nil {
				log.Printf("Failed generate cookie: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, cookie)
		}

		id, err := uuid.Parse(cookie.Value)

		if err != nil {
			log.Printf("Failed parse user cookie: %v\n", err)
			cookie, err = generateCookie()
			if err != nil {
				log.Printf("Failed generate cookie: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, cookie)

			id, err = uuid.Parse(cookie.Value)
			if err != nil {
				log.Printf("Failed parse just setted cookie: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		ctx := context.WithValue(r.Context(), "player_id", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
