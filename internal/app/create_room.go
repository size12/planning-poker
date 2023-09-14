package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/olahol/melody"
	"github.com/size12/planning-poker/internal/app/websocket"
	"github.com/size12/planning-poker/internal/service/room"
)

// CreateRoom is handler for POST /rooms/create.
func (app *App) CreateRoom(writer http.ResponseWriter, request *http.Request) {
	r := &room.Room{}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Printf("Failed read request body: %v\n", err)
		return
	}
	defer request.Body.Close()

	err = json.Unmarshal(body, r)
	if err != nil {
		http.Error(writer, "can't unmarshal request body", http.StatusBadRequest)
		return
	}

	r, err = app.handlers.CreateRoom(r.Name)
	if err != nil {
		http.Error(writer, fmt.Sprintf("can't create request: %v", err), http.StatusBadRequest)
		return
	}

	adminID, ok := request.Context().Value("player_id").(uuid.UUID)
	if !ok {
		log.Printf("Waited uuid.UUID value from context, got: %v\n", adminID)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.SetAdmin(adminID)

	m := melody.New()

	m.HandleConnect(websocket.WsConnected(r, m))
	m.HandleDisconnect(websocket.WsDisconnected(r, m))
	m.HandleMessage(websocket.WsOnMessage(r, m))

	app.wsPool.Store(r.ID, m)

	response, err := json.Marshal(r)
	if err != nil {
		log.Printf("Failed marshal response: %v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	_, err = writer.Write(response)
	if err != nil {
		log.Printf("Failed write data to client: %v\n", err)
		return
	}

	log.Printf("Successfully created room with id %s\n", r.ID)
}
