package website

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/size12/planning-poker/internal/handlers"
)

var RoomFilename = "room.html"

// Room loads page where you can vote and see other players.
// GET rooms/{id:[a-z0-9-]+}.
func (site *Website) Room(writer http.ResponseWriter, request *http.Request, h *handlers.Handlers) {
	id := chi.URLParam(request, "id")

	roomID, err := uuid.Parse(id)

	if err != nil {
		site.NotFound(writer, request)
		return
	}

	_, err = h.RoomByID(roomID)
	switch {
	case errors.Is(err, handlers.ErrRoomNotFound):
		{
			site.NotFound(writer, request)
			return
		}
	case err != nil:
		{
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err = site.files.Lookup(RoomFilename).Execute(writer, site.url)
	if err != nil {
		log.Printf("Failed execute room template: %v\n", err)
		return
	}
}
