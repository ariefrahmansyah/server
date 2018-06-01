package server

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/ariefrahmansyah/server/router"
	"github.com/cockroachdb/cmux"
	"golang.org/x/net/netutil"
)

// Options for the web server.
type Options struct {
	ListenAddress  string
	MaxConnections int

	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	version *Version
}

// Server serves various HTTP endpoints.
type Server struct {
	logger  *log.Logger
	options *Options
	router  *router.Router
}

// New initializes a new web server.
func New(logger *log.Logger, options *Options) *Server {
	if logger == nil {
		logger = log.New(ioutil.Discard, "", 0)
	}

	server := &Server{
		logger:  logger,
		options: options,
	}

	router := router.New()
	router.Get("/ping", server.Ping)
	router.Get("/version", server.Version)

	server.router = router

	return server
}

// Run serves the web server.
func (s *Server) Run(ctx context.Context) error {
	// Create the main listener
	listener, err := net.Listen("tcp", s.options.ListenAddress)
	if err != nil {
		return err
	}
	listener = netutil.LimitListener(listener, s.options.MaxConnections)

	// Listner multiplexer
	listenerMux := cmux.New(listener)
	httpListener := listenerMux.Match(cmux.HTTP1Fast())

	// TODO: grpc listener

	// HTTP server
	httpServer := &http.Server{
		Handler:     s.router,
		ReadTimeout: s.options.ReadTimeout,
		ErrorLog:    s.logger,
	}

	// Start listening
	errCh := make(chan error)
	go func() {
		errCh <- httpServer.Serve(httpListener)
	}()
	go func() {
		errCh <- listenerMux.Serve()
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		httpServer.Shutdown(ctx)
		return nil
	}
}

// Router of web server.
func (s *Server) Router() *router.Router {
	return s.router
}

// Ping writes pong.
func (s *Server) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
