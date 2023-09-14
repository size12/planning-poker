package app

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/size12/planning-poker/internal/service/room"
)

func (app *App) GetRoom(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Printf("Failed read request body: %v\n", err)
		return
	}
	defer request.Body.Close()

	r := &room.Room{}

	err = json.Unmarshal(body, r)
	if err != nil {
		http.Error(writer, "can't unmarshal request body", http.StatusBadRequest)
		return
	}

	r, err = app.handlers.RoomByID(r.ID)
	if err != nil {
		http.Error(writer, "failed find room with such id", http.StatusNotFound)
		return
	}

	r.RLock()
	defer r.RUnlock()

	response, err := json.Marshal(r)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed marshal room %s: %v\n", r.ID, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	_, err = writer.Write(response)
	if err != nil {
		log.Printf("Failed write data to client: %v\n", err)
		return
	}
}
