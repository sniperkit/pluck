package parse

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/slotix/dataflowkit/logger"
	"github.com/slotix/dataflowkit/storage"
	"github.com/spf13/viper"
)

var storageType storage.Type

// Config provides basic configuration
type Config struct {
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// HTMLServer represents the web service that serves up HTML
type HTMLServer struct {
	server *http.Server
	wg     sync.WaitGroup
}

// Start func launches Parsing service
func Start(cfg Config) *HTMLServer {
	ctx := context.Background()
	_, cancel := context.WithCancel(ctx)
	defer cancel()
	logger := log.NewLogger(false)

	var svc Service
	svc = ParseService{}
	//svc = StatsMiddleware("18")(svc)
	if !viper.GetBool("SKIP_STORAGE_MW") {
		var err error
		storageType, err = storage.TypeString(viper.GetString("STORAGE_TYPE"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		svc = StorageMiddleware(storage.NewStore(storageType))(svc)
	}
	svc = LoggingMiddleware(logger)(svc)

	endpoints := Endpoints{
		ParseEndpoint: MakeParseEndpoint(svc),
	}

	r := NewHttpHandler(ctx, endpoints, logger)

	// Create the HTML Server
	htmlServer := HTMLServer{
		server: &http.Server{
			Addr:           cfg.Host,
			Handler:        r,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
			MaxHeaderBytes: 1 << 20,
		},
	}
	// Add to the WaitGroup for the listener goroutine
	htmlServer.wg.Add(1)

	go func() {
		fmt.Printf("Starting Parse Server %s\n", htmlServer.server.Addr)
		htmlServer.server.ListenAndServe()
		htmlServer.wg.Done()
	}()
	return &htmlServer
}

// Stop turns off the HTML Server
func (htmlServer *HTMLServer) Stop() error {
	// Create a context to attempt a graceful 5 second shutdown.
	const timeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fmt.Printf("\nFetch Server : Service stopping\n")

	// Attempt the graceful shutdown by closing the listener
	// and completing all inflight requests
	if err := htmlServer.server.Shutdown(ctx); err != nil {
		// Looks like we timed out on the graceful shutdown. Force close.
		if err := htmlServer.server.Close(); err != nil {
			fmt.Printf("\nFetch Server : Service stopping : Error=%v\n", err)
			return err
		}
	}
	// Wait for the listener to report that it is closed.
	htmlServer.wg.Wait()
	fmt.Printf("\nFetch Server : Stopped\n")
	return nil
}