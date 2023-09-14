package websocket

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/olahol/melody"
	"github.com/size12/planning-poker/internal/entity/update"
	"github.com/size12/planning-poker/internal/service/room"
)

func WsOnMessage(r *room.Room, m *melody.Melody) func(*melody.Session, []byte) {
	return func(session *melody.Session, msg []byte) {
		u := &update.Update{}

		err := json.Unmarshal(msg, u)
		if err != nil {
			response := &Response{
				Code:    http.StatusBadRequest,
				Message: "Failed unmarshal json: " + err.Error(),
			}

			bytes, err := json.Marshal(response)
			if err != nil {
				log.Printf("Failed marshal response message: %v\n", err)
				return
			}

			err = session.Write(bytes)
			if err != nil {
				log.Printf("Failed write to client: %v\n", err)
				return
			}
		}

		log.Printf("Got websocket message: %v\n", u)

		err = r.Update(u)
		if err != nil {

			response := &Response{
				Message: err.Error(),
			}

			switch {
			case errors.Is(err, room.ErrAccessDenied):
				{
					response.Code = http.StatusForbidden
				}
			case errors.Is(err, room.ErrPlayerNotFound):
				{
					response.Code = http.StatusNotFound
				}
			default:
				response.Code = http.StatusInternalServerError
			}

			bytes, err := json.Marshal(response)
			if err != nil {
				log.Printf("Failed marshal response message: %v\n", err)
				return
			}

			err = session.Write(bytes)
			if err != nil {
				log.Printf("Failed write to client: %v\n", err)
				return
			}
		}

		r.RLock()
		defer r.RUnlock()

		bytes, err := json.Marshal(r)
		if err != nil {
			log.Printf("Failed marshal room: %v\n", err)
			return
		}

		err = m.Broadcast(bytes)
		if err != nil {
			log.Printf("Failed broadcast: %v\n", err)
			return
		}

	}
}
