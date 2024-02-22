package rest

import (
	"awds/commons"
	"awds/logic"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type RESTAdapter struct {
	config     *commons.Config
	address    string
	router     *gin.Engine
	httpServer *http.Server
	logic      *logic.Logic
}

// Start starts RESTAdapter
func Start(config *commons.Config, logik *logic.Logic) (*RESTAdapter, error) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"function": "Start",
	})

	addr := fmt.Sprintf(":%d", config.RestPort)
	router := gin.Default()
	router.Use(cors.New(
		cors.Config{
			AllowOrigins: []string{"http://localhost:5173", "http://localhost:5174", "http://127.0.0.1:5173", "http://127.0.0.1:5174",
				"http://155.230.36.27:5173", "http://155.230.36.27:5174", "http://155.230.36.27:4140", "http://155.230.36.27:4141"},
			AllowMethods: []string{"POST", "GET", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
			// allow headers
			AllowCredentials: true,
			// allow credentials
			MaxAge: 24 * time.Hour,
		}))

	adapter := &RESTAdapter{
		config:  config,
		address: addr,
		router:  router,
		httpServer: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		logic: logik,
	}

	// setup HTTP request router
	adapter.setupRouter()

	fmt.Printf("Starting REST service at %s\n", adapter.address)
	logger.Infof("Starting REST service at %s\n", adapter.address)
	// listen and serve in background
	go func() {
		err := adapter.httpServer.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	return adapter, nil
}

// Stop stops RESTAdapter
func (adapter *RESTAdapter) Stop() error {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "Stop",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := adapter.httpServer.Shutdown(ctx)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
