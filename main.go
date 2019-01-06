package main

// v0.0.0
import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"

	"github.com/nimbo-stratuz/bikeshare-directions/api"
	"github.com/nimbo-stratuz/bikeshare-directions/service"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		if r.Header.Get("X-Request-ID") != "" {
			ctx = context.WithValue(ctx, middleware.RequestIDKey, r.Header.Get("X-Request-ID"))
		} else {
			ctx = context.WithValue(ctx, middleware.RequestIDKey, uuid.New().String())
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Routes for /v1
func Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(
		render.SetContentType(render.ContentTypeJSON),

		requestIDMiddleware,
		middleware.Logger,
	)

	r.Mount("/", api.Routes())

	return r
}

func main() {

	if env, err := service.Config.Get("env"); err != nil || env == "prod" {
		log.SetLevel(log.InfoLevel)
		log.SetFormatter(&log.JSONFormatter{})
	}

	// Make sure application quits gracefully
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGTERM)
	signal.Notify(exit, syscall.SIGINT)
	go func() {
		<-exit
		service.Close()
		os.Exit(0)
	}()

	r := Routes()

	port, err := service.Config.Get("server", "port")
	if err != nil {
		log.Println("PORT not specified")
	}
	log.Fatal(http.ListenAndServe(":"+port, r))
}
