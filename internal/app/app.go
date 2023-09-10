package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/size12/planning-poker/internal/config"
	"github.com/size12/planning-poker/internal/entity"
	"github.com/size12/planning-poker/internal/handlers"
	"github.com/size12/planning-poker/internal/service"
)

type App struct {
	cfg      *config.Config
	server   *http.Server
	handlers *handlers.Handlers

	shutdownPool context.CancelFunc
}

func NewApp(cfg *config.Config) (*App, error) {
	return &App{cfg: cfg}, nil
}

func (app *App) Run() {
	pool, err := service.NewRoomPool()
	if err != nil {
		log.Fatalf("Failed get new room pool: %v\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	app.shutdownPool = cancel

	go pool.DeleteInactive(ctx, app.cfg.InactiveTimeout)

	h, err := handlers.NewHandlers(pool)
	if err != nil {
		log.Fatalf("Failed get new handlers: %v\n", err)
	}

	app.handlers = h

	router := chi.NewRouter()

	router.Get("/createRoom", app.CreateRoom)

	server := &http.Server{
		Addr:    app.cfg.RunAddress,
		Handler: router,
	}

	app.server = server

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed run server on %s: %v\n", app.cfg.RunAddress, err)
		}
	}()
}

func (app *App) Shutdown() {
	// shutdown delete inactive rooms service.
	app.shutdownPool()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := app.server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed shutdown server gracefully: %v\n", err)
	}
	log.Println("Server shutdown successfully")
}

func (app *App) CreateRoom(w http.ResponseWriter, r *http.Request) {
	room := &entity.Room{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read request body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, room)
	if err != nil {
		http.Error(w, "can't unmarshal request body", http.StatusBadRequest)
		return
	}

	room, err = app.handlers.CreateRoom(room.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("can't create room: %v", err), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(room)
	if err != nil {
		http.Error(w, "can't marshal room", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(response)
	if err != nil {
		log.Printf("Failed write data to client: %v\n", err)
		return
	}
}
