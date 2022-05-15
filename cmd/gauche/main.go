package main

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/tjvc/gauche/internal/handler"
	"github.com/tjvc/gauche/internal/logging"
	"github.com/tjvc/gauche/internal/middleware"
	"github.com/tjvc/gauche/internal/store"
)

type application struct {
	store    *store.Store
	logger   logging.Logger
	maxMemMB int
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

		application.store.Set("dummy", []byte("123"))

		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		usedMemMB := float64(m.TotalAlloc) / 1024 / 1024
		maxMemMB := float64(application.maxMemMB)
		fmt.Printf("usedMemMB: %f\n", usedMemMB)
		fmt.Printf("maxMemMB: %f\n", maxMemMB)
		fmt.Printf("TotalAlloc: %d\n", m.TotalAlloc)
		fmt.Printf("Alloc: %d\n", m.Alloc)
		fmt.Printf("Store size: %d\n", application.store.Size)
		fmt.Printf("Ratio: %f\n", float64(m.Alloc)/float64(application.store.Size))

		key := r.URL.Path[1:]

		switch r.Method {
		case http.MethodGet:
			handler.Get(w, key, application.store)
		case http.MethodPut:
			if usedMemMB > maxMemMB {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

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
	maxMemMB := 1024

	application := application{
		store:    &store,
		logger:   logger,
		maxMemMB: maxMemMB,
	}

	http.ListenAndServe(":8080", mainHandler(application))
}
