package main

import (
	"net/http"

	"github.com/tjvc/gauche/internal/handler"
	"github.com/tjvc/gauche/internal/logging"
	"github.com/tjvc/gauche/internal/middleware"
	"github.com/tjvc/gauche/internal/store"
)

type application struct {
	store  *store.Store
	logger logging.Logger
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
	logger := logging.JSONLogger{}

	application := application{
		store:  &store,
		logger: logger,
	}

	http.ListenAndServe(":8080", mainHandler(application))
}
