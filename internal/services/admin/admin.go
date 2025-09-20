package admin

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

type Admin struct {
	wg       *sync.WaitGroup
	port     uint16
	once     sync.Once
	adminAPI *API
}

func New(port uint16) *Admin {
	return &Admin{
		port:     port,
		adminAPI: api(),
	}
}

func (a *Admin) Start(ctx context.Context, wg *sync.WaitGroup) error {
	var startErr error
	a.once.Do(func() {
		a.wg = wg
		a.wg.Add(1)
		go func() {
			startErr = a.start(ctx)
			a.wg.Done()
		}()
	})
	return startErr
}

func (a *Admin) start(ctx context.Context) error {
	r := chi.NewRouter()

	// Mount the admin API routes
	r.Mount("/", a.adminAPI.Routes())

	// Keep the existing health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	srv := &http.Server{
		Addr:    ":" + fmt.Sprintf("%d", a.port), // Use the port from the Admin struct
		Handler: r,
	}

	go func() {
		<-ctx.Done()
		slog.Info("admin shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("admin shutdown error", "err", err)
		} else {
			slog.Info("admin shutdown complete")
		}
	}()

	slog.Info("admin API starting on :" + fmt.Sprintf("%d", a.port))
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("admin API server error", "err", err)
	}
	return err
}

// Stop shuts down the Admin service
func (a *Admin) Stop() {
	// Implement shutdown logic here
}
