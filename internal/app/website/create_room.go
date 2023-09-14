package website

import (
	"log"
	"net/http"
)

var CreateRoomFilename = "create_room.html"

// CreateRoom loads page where you can put room name and create it.
// GET /rooms/create.
func (site *Website) CreateRoom(writer http.ResponseWriter, request *http.Request) {
	err := site.files.Lookup(CreateRoomFilename).Execute(writer, nil)
	if err != nil {
		log.Printf("Failed execute create room template: %v\n", err)
		return
	}
}
