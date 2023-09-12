package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/olahol/melody"
	"github.com/size12/planning-poker/internal/config"
	"github.com/size12/planning-poker/internal/entity"
	"github.com/size12/planning-poker/internal/entity/update"
	"github.com/size12/planning-poker/internal/entity/voting"
	"github.com/size12/planning-poker/internal/handlers"
	"github.com/size12/planning-poker/internal/middleware"
	"github.com/size12/planning-poker/internal/service"
)

type App struct {
	cfg      *config.Config
	server   *http.Server
	handlers *handlers.Handlers
	router   chi.Router

	wsPool *sync.Map
	files  *template.Template
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

	temp, err := template.ParseFiles("website/404.html", "website/create_room.html", "website/room.html")
	if err != nil {
		log.Fatalf("Failed parse templates: %v\n", err)
	}

	return &App{handlers: h, cfg: cfg, server: server, wsPool: &sync.Map{}, router: router, files: temp}, nil
}

func (app *App) Run() {

	router := app.router
	router.Use(middleware.Auth)

	router.Route("/rooms", func(r chi.Router) {
		r.Get("/create", app.CreateRoomWindow)
		r.Post("/create", app.CreateRoom)
		r.Get("/status", app.GetRoom)
		r.Get("/{id:[a-z0-9-]+}", app.Room)
	})

	router.HandleFunc("/ws/{id:[a-z0-9-]+}", app.wsRoom)

	router.NotFound(app.NotFound)

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
	log.Println("Server shutdown successfully")
}

func (app *App) NotFound(w http.ResponseWriter, r *http.Request) {
	path := "website/404.html"
	w.WriteHeader(http.StatusNotFound)
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Failed open not found page: %v\n", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Failed read not found page: %v\n", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("Failed write to client: %v\n", err)
	}

	return
}

func (app *App) CreateRoom(w http.ResponseWriter, r *http.Request) {
	room := &service.Room{}

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

	adminID := r.Context().Value("player_id").(uuid.UUID)
	room.AdminID = adminID

	m := melody.New()

	m.HandleConnect(wsConnected(room, m))
	m.HandleDisconnect(wsDisconnected(room, m))
	m.HandleMessage(wsMessage(room, m))

	app.wsPool.Store(room.ID, m)

	response, err := json.Marshal(room)
	if err != nil {
		http.Error(w, "can't marshal room", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(response)
	if err != nil {
		log.Printf("Failed write data to client: %v\n", err)
		return
	}
}

func (app *App) CreateRoomWindow(w http.ResponseWriter, r *http.Request) {

	err := app.files.Lookup("create_room.html").Execute(w, nil)
	if err != nil {
		log.Printf("Failed parse room file: %v\n", err)
		app.NotFound(w, r)
		return
	}
}

func (app *App) wsRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
	id, err := uuid.Parse(roomID)
	if err != nil {
		http.Error(w, "failed parse room_id", http.StatusBadRequest)
		return
	}

	playerID, ok := r.Context().Value("player_id").(uuid.UUID)
	if !ok {
		http.Error(w, "failed get player_id", http.StatusBadRequest)
		return
	}

	object, ok := app.wsPool.Load(id)
	if !ok {
		app.NotFound(w, r)
		return
	}

	ws, ok := object.(*melody.Melody)
	if !ok {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	keys := make(map[string]interface{})
	keys["room_id"] = id
	keys["player_id"] = playerID

	err = ws.HandleRequestWithKeys(w, r, keys)
	if err != nil {
		log.Printf("Failed handle websocket connection: %v\n", err)
		return
	}
}

func (app *App) Room(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
	id, err := uuid.Parse(roomID)
	if err != nil {
		app.NotFound(w, r)
		return
	}

	room, err := app.handlers.GetRoom(id)
	if err != nil {
		app.NotFound(w, r)
		return
	}

	playerID, ok := r.Context().Value("player_id").(uuid.UUID)
	if !ok {
		log.Printf("Failed get player_id from r.Context()\n")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = room.PlayerByID(playerID)
	log.Println("conn", playerID, err)
	if err == nil {
		w.WriteHeader(http.StatusConflict)
		_, err = w.Write([]byte("<div>please return to first tab</div><script>alert(\"please return to first tab\")</script>"))
		if err != nil {
			log.Printf("Failed write response to client: %v\n", err)
			return
		}
		return
	}

	temp := struct {
		Url      string
		RoomID   string
		PlayerID string
		IsAdmin  bool
	}{
		app.cfg.BaseURL,
		roomID,
		playerID.String(),
		playerID == room.AdminID,
	}

	err = app.files.Lookup("room.html").Execute(w, temp)
	if err != nil {
		log.Printf("Failed parse room file: %v\n", err)
		app.NotFound(w, r)
		return
	}
}

func (app *App) GetRoom(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed read request body: %v\n", err)
		return
	}
	defer r.Body.Close()

	info := &service.RoomInfo{}

	err = json.Unmarshal(body, info)
	if err != nil {
		http.Error(w, "Failed unmarshal request body", http.StatusBadRequest)
		return
	}

	room, err := app.handlers.GetRoom(info.ID)
	if err != nil {
		http.Error(w, "failed find room with such id", http.StatusNotFound)
		return
	}

	response, err := json.Marshal(room)
	if err != nil {
		http.Error(w, "failed marshal room", http.StatusInternalServerError)
		log.Printf("Failed marshal room: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(response)
	if err != nil {
		log.Printf("Failed write data to client: %v\n", err)
		return
	}
}

func wsConnected(room *service.Room, m *melody.Melody) func(*melody.Session) {
	return func(session *melody.Session) {
		//fmt.Printf("connected: %+v\n", session)
	}
}

func wsMessage(room *service.Room, m *melody.Melody) func(*melody.Session, []byte) {
	return func(session *melody.Session, msg []byte) {
		u := &update.Update{}

		err := json.Unmarshal(msg, u)
		if err != nil {
			err = session.Write([]byte("400"))
			if err != nil {
				log.Printf("Failed write to client: %v\n", err)
				return
			}
		}

		id := u.PlayerID

		log.Println(u)

		switch u.Type {
		case update.Connected:
			player, err := entity.NewPlayer(id)
			if err != nil {
				log.Printf("Failed create new player: %v\n", err)
				err = session.Write([]byte("500"))
				if err != nil {
					log.Printf("Failed write to client: %v\n", err)
					return
				}
			}

			err = room.AddPlayer(player)
			if err != nil {
				err = session.Write([]byte("500"))
				if err != nil {
					log.Printf("Failed write to client: %v\n", err)
					return
				}
			}
		case update.Reveal:
			{
				if id != room.AdminID {
					_ = session.Write([]byte("403"))
					return
				}

				err = room.SetStatus(voting.StatusRevealed)
				if err != nil {
					//log.Printf("Failed set status: %v\n", err)
					return
				}
			}
		case update.StartNewVoting:
			{
				if id != room.AdminID {
					_ = session.Write([]byte("403"))
					return
				}

				err = room.ClearVotes()
				if err != nil {
					log.Printf("Failed clear votes: %v\n", err)
					return
				}
				err = room.SetStatus(voting.StatusVoting)
				if err != nil {
					log.Printf("Failed set status: %v\n", err)
					return
				}
			}
		case update.Vote:
			{
				err = room.Vote(id, u.Vote)
				if err != nil {
					log.Printf("Failed vote: %v\n", err)
					return
				}
			}
		}

		bytes, err := json.Marshal(room)
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

func wsDisconnected(room *service.Room, m *melody.Melody) func(*melody.Session) {
	return func(session *melody.Session) {
		//fmt.Printf("disconnected: %+v\n", session)

		object, ok := session.Get("player_id")
		if !ok {
			log.Println("Someone without user_id disconnected")
			return
		}

		id, ok := object.(uuid.UUID)
		if !ok {
			log.Println("Someone without valid user_id disconnected")
			return
		}

		err := room.RemovePlayer(id)
		if err != nil {
			log.Printf("Failed remove player from room: %v\n", err)
			return
		}

		bytes, err := json.Marshal(room)
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

//ws := melody.New()
//
//app.ws = ws
//
//ws.HandleMessage(func(s *melody.Session, msg []byte) {
//	err = ws.Broadcast(msg)
//	if err != nil {
//		log.Printf("Failed broadcast message: %v\n", err)
//	}
//	fmt.Println(string(msg))
//	fmt.Printf("%+v\n", s)
//})
//
//ws.HandleDisconnect(func(s *melody.Session) {
//	value, ok := s.Get("room_id")
//	if !ok {
//		log.Println("Connection doesn't have room_id")
//		return
//	}
//
//	id, ok := value.(uuid.UUID)
//	if !ok {
//		log.Println("Connection doesn't have room_id")
//		return
//	}
//
//	room, err := app.handlers.GetRoom(id)
//	if err != nil {
//		log.Println("Someone disconnected from non-existed room")
//		return
//	}
//
//	value, ok = s.Get("player_id")
//	if !ok {
//		log.Println("Connection doesn't have player_id")
//		return
//	}
//
//	id, ok = value.(uuid.UUID)
//	if !ok {
//		log.Println("Connection doesn't have valid player_id")
//		return
//	}
//
//	err = room.RemovePlayer(id)
//	if err != nil {
//		log.Println("Can't find player in room")
//		return
//	}
//
//})
