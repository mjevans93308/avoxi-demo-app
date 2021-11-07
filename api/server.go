package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/mjevans93308/avoxi-demo-app/config"
)

// App structured as standalone object to better facilitate testing
type App struct {
	Router *mux.Router
}

func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(address string) {

	srv := &http.Server{
		Handler:           a.Router,
		Addr:              address,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}

	// create go chan to notify if server stops
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// run server inside goroutine
	// not strictly necessary for a demo, but safer for err detection and more realistic for real world app
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error running server: %s\n", err.Error())
		}
	}()

	// check channel for if we've stopped server
	<-done
	log.Print("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// gracefully drain connections if applicable
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc(config.HOME, a.homeHandler).Methods("GET")

	s := a.Router.PathPrefix(config.API_GROUP + config.V1_GROUP).Subrouter()
	s.HandleFunc(config.ALIVE, a.aliveHandler)

	// register additional routes here
}

func (a *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello, world!"))
}

func (a *App) aliveHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("It's...ALIVE!!!"))
}

// func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
// 	response, _ := json.Marshal(payload)

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	w.Write(response)
// }
