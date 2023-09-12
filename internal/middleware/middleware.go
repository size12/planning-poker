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

//func Room(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		id, err := uuid.Parse(chi.URLParam(r, "id"))
//		if err != nil {
//			http.Error(w, "failed read room id from url", http.StatusBadRequest)
//			return
//		}
//
//		ctx := context.WithValue(r.Context(), "room_id", id)
//		next.ServeHTTP(w, r.WithContext(ctx))
//	})
//}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("poker_player_id")
		if err != nil {
			cookie, err = generateCookie()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, cookie)
			log.Printf("set new cookie %v\n", cookie)
		}

		id, err := uuid.Parse(cookie.Value)

		if err != nil {
			log.Printf("Failed parse user cookie: %v\n", err)
			cookie, err = generateCookie()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, cookie)
			log.Printf("set new cookie %v\n", cookie)

			id, err = uuid.Parse(cookie.Value)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		ctx := context.WithValue(r.Context(), "player_id", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
