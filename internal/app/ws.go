package app

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/olahol/melody"
)

func (app *App) handleWS(writer http.ResponseWriter, request *http.Request) {
	roomID := chi.URLParam(request, "id")
	id, err := uuid.Parse(roomID)
	if err != nil {
		http.Error(writer, "failed parse room_id", http.StatusBadRequest)
		return
	}

	playerIDCookie, ok := request.Context().Value("player_id").(uuid.UUID)
	if !ok {
		http.Error(writer, "failed get player_id", http.StatusBadRequest)
		return
	}

	object, ok := app.wsPool.Load(id)
	if !ok {
		app.website.NotFound(writer, request)
		return
	}

	playerID, err := uuid.NewUUID()
	if err != nil {
		log.Printf("Failed generate uuid: %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
	}

	ws, ok := object.(*melody.Melody)
	if !ok {
		log.Printf("Wanted *melody.Melody object in *sync.Map, but got %T\n", ws)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	keys := make(map[string]interface{})
	keys["room_id"] = id
	keys["player_id"] = playerID           // only for current room.
	keys["player_cookie"] = playerIDCookie // for admin access.

	err = ws.HandleRequestWithKeys(writer, request, keys)
	if err != nil {
		log.Printf("Failed handle websocket connection: %v\n", err)
		return
	}
}
