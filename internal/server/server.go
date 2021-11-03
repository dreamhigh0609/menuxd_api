package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	v1 "github.com/super-radmir/menuxd_api/internal/server/v1"
)

type config struct {
	port  string
	debug bool
}

var conf config
var fileLog *os.File

// SetConfig set server configuration.
func SetConfig(port string, debug bool) {
	conf.port = port
	conf.debug = debug
}

var basePath = os.Getenv("XD_BASE_PATH")

func getRoutes() (http.Handler, error) {
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
			"PATCH",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r := chi.NewRouter()
	r.Use(cors.Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(setLogger(conf.debug))
	r.Use(middleware.Recoverer)

	// Configurations
	r.Handle("/public/*", http.StripPrefix(
		"/public", http.FileServer(http.Dir(basePath+"public")),
	))
	v1Routes, err := v1.NewAPI()
	if err != nil {
		return nil, err
	}

	r.Mount("/api/v1", v1Routes)
	r.Handle(
		"/docs/*",
		http.StripPrefix("/docs/", http.FileServer(http.Dir(basePath+"docs"))),
	)
	r.Handle("/*", http.FileServer(http.Dir(basePath+"static")))

	return r, nil
}

// New inicialize a new server with configuration.
func New() (*http.Server, error) {
	r, err := getRoutes()
	if err != nil {
		return nil, err
	}

	srv := &http.Server{
		Addr:         ":" + conf.port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return srv, nil
}

func setLogger(isDebug bool) func(next http.Handler) http.Handler {
	if isDebug {
		return middleware.Logger
	}

	var err error
	fileLog, err = os.OpenFile("logs.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return middleware.Logger
	}

	logger := middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(fileLog, "", log.LstdFlags), NoColor: true})

	return logger
}

func Close() error {
	if fileLog != nil {
		return fileLog.Close()
	}

	return nil
}
