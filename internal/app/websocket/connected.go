package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/olahol/melody"
	"github.com/size12/planning-poker/internal/entity/access"
	"github.com/size12/planning-poker/internal/service/room"
)

func WsConnected(r *room.Room, m *melody.Melody) func(*melody.Session) {
	return func(session *melody.Session) {
		object, ok := session.Get("player_id")
		if !ok {
			log.Printf("Someone without player_id connected from room %s\n", r.ID)
			return
		}

		ID, ok := object.(uuid.UUID)
		if !ok {
			log.Printf("Someone without valid player_id connected to room %s\n", r.ID)
			return
		}

		log.Printf("User %s connected to room %s\n", ID, r.ID)

		object, ok = session.Get("player_cookie")
		if !ok {
			log.Printf("Someone without player_cookie connected from room %s\n", r.ID)
			return
		}

		cookieID, ok := object.(uuid.UUID)
		if !ok {
			log.Printf("Someone without valid player_cookie connected to room %s\n", r.ID)
			return
		}

		msg := &Response{}

		if r.IsAdmin(cookieID) {
			msg.Code = http.StatusOK
			msg.Message = string(access.Admin) + ":" + ID.String()
		} else {
			msg.Code = http.StatusForbidden
			msg.Message = string(access.User) + ":" + ID.String()
		}

		bytes, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Failed marshal response %v\n", err)
			return
		}

		err = session.Write(bytes)
		if err != nil {
			log.Printf("Failed write to client %v\n", err)
			return
		}
	}
}
