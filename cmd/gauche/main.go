package main

import (
	"net/http"
	"os"

	"github.com/tjvc/gauche/internal/handler"
	"github.com/tjvc/gauche/internal/middleware"
	"github.com/tjvc/gauche/internal/store"
	"golang.org/x/exp/slog"
)

type application struct {
	store  *store.Store
	logger *slog.Logger
}

func mainHandler(application application) http.Handler {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			if r.Method == http.MethodGet {
				handler.Index(w, application.store)
				return
			}

			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		key := r.URL.Path[1:]

		switch r.Method {
		case http.MethodGet:
			handler.Get(w, key, application.store)
		case http.MethodPut:
			handler.Put(w, key, r, application.store)
		case http.MethodDelete:
			handler.Delete(w, key, application.store)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux := http.NewServeMux()
	mux.Handle("/", middleware.Log(application.logger, middleware.Recover(h)))
	return mux
}

func main() {
	store := store.New()

	application := application{
		store:  &store,
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	http.ListenAndServe(":8080", mainHandler(application))
}
