package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/mjevans93308/avoxi-demo-app/config"
	"github.com/mjevans93308/avoxi-demo-app/utils"
)

var logger = utils.InitLogger()

// App structured as standalone object to better facilitate testing
type App struct {
	Router *gin.Engine
}

func (a *App) Initialize() {
	logger.Info("App initialized")
	a.Router = gin.Default()
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
			logger.Errorf("Error running server: %s\n", err.Error())
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
		logger.Errorf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}

func (a *App) initializeRoutes() {

	// create router with Alive endpoint for uptime checking
	a.Router.GET(config.Alive, a.aliveHandler)

	// created basic auth secured subrouter for business logic endpoints
	// gin supports basic auth cred pairs, that if not matched will fail the request with a 401
	authenticated := a.Router.Group(config.Api_Group+config.V1_Group, gin.BasicAuth(gin.Accounts{
		viper.GetString(config.Basic_Auth_Username): viper.GetString(config.Basic_Auth_Password),
	}))

	authenticated.POST(config.CheckIPLocation, a.CheckGeoLocation)
}
