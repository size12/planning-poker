package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	chim "github.com/go-chi/chi/v5/middleware"
	"github.com/size12/planning-poker/internal/app/website"
	"github.com/size12/planning-poker/internal/config"
	"github.com/size12/planning-poker/internal/handlers"
	"github.com/size12/planning-poker/internal/middleware"
)

type App struct {
	cfg      *config.Config
	server   *http.Server
	handlers *handlers.Handlers
	router   chi.Router

	wsPool *sync.Map

	website *website.Website
}

func NewApp(cfg *config.Config) (*App, error) {
	h, err := handlers.NewHandlers()
	if err != nil {
		return nil, err
	}

	router := chi.NewRouter()

	server := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: router,
	}

	site, err := website.New(cfg.BaseURL)
	if err != nil {
		return nil, err
	}

	return &App{
		handlers: h,
		cfg:      cfg,
		server:   server,
		wsPool:   &sync.Map{},
		router:   router,
		website:  site}, nil
}

func (app *App) Run() {
	site := app.website

	router := app.router
	router.Use(chim.Recoverer)
	router.Use(middleware.Auth)

	router.Route("/rooms", func(r chi.Router) {
		r.Get("/create", site.CreateRoom)
		r.Post("/create", app.CreateRoom)

		r.Get("/status", app.GetRoom)
		r.Get("/{id:[a-z0-9-]+}", func(writer http.ResponseWriter, request *http.Request) {
			site.Room(writer, request, app.handlers)
		})
	})

	router.HandleFunc("/ws/{id:[a-z0-9-]+}", app.handleWS)

	router.NotFound(site.NotFound)

	go func() {
		err := app.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed run server on %s: %v\n", app.cfg.RunAddress, err)
		}
	}()
}

func (app *App) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := app.server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed shutdown server gracefully: %v\n", err)
	}
	log.Println("Server shutdown gracefully")
}
