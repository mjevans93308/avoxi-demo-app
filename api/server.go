package api

import (
	"context"
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
	Router   *gin.Engine
	Outbound *Outbound
}

func (a *App) Initialize(testing bool) {
	a.Router = gin.Default()
	if viper.GetString(config.Environment) == config.TestEnv {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	a.initializeRoutes(testing)
	a.Outbound = initOutboundApi()
	logger.Info("App initialized")
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
	logger.Info("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// gracefully drain connections if applicable
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("Server Shutdown Failed:%+v", err)
	}
	logger.Info("Server Exited Properly")
}

// initializeRoutes attaches all routes and subrouters to our App instance
func (a *App) initializeRoutes(testing bool) {

	// create router with Alive endpoint for uptime checking
	a.Router.GET(config.Alive, a.aliveHandler)
	a.Router.POST(config.Inform, a.informHandler)
	a.Router.GET(config.Teapot, a.teapotHandler)

	// created basic auth secured subrouter for business logic endpoints
	// gin supports basic auth cred pairs that if not matched will fail the request with a 401

	if testing {
		noauth := a.Router.Group(config.Api_Group + config.V1_Group)
		noauth.POST(config.CheckIPLocation, a.CheckGeoLocation)
	} else {
		authenticated := a.Router.Group(config.Api_Group+config.V1_Group, gin.BasicAuth(gin.Accounts{
			viper.GetString(config.Basic_Auth_Username): viper.GetString(config.Basic_Auth_Password),
		}))

		authenticated.POST(config.CheckIPLocation, a.CheckGeoLocation)
	}

}
