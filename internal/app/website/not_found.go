package website

import (
	"io"
	"log"
	"net/http"
	"os"
)

var NotFoundPath = "website/404.html"

func (site *Website) NotFound(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotFound)

	file, err := os.Open(NotFoundPath)
	if err != nil {
		log.Printf("Failed open not found page: %v\n", err)
		return
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Failed read not found page: %v\n", err)
		return
	}

	writer.Header().Set("Content-Type", "text/html")

	_, err = writer.Write(bytes)
	if err != nil {
		log.Printf("Failed write to client: %v\n", err)
	}

	return
}
