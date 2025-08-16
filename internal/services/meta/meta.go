package meta

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

type Meta struct {
	wg   *sync.WaitGroup
	port uint16
	once sync.Once
}

func New(port uint16) *Meta {
	return &Meta{port: port}
}

func (m *Meta) Start(ctx context.Context, wg *sync.WaitGroup) error {
	var startErr error
	m.once.Do(func() {
		m.wg = wg
		m.wg.Add(1)
		go func() {
			startErr = m.start(ctx)
			m.wg.Done()
		}()
	})
	return startErr
}

func (m *Meta) start(ctx context.Context) error {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	srv := &http.Server{
		Addr:    ":" + fmt.Sprintf("%d", m.port),
		Handler: r,
	}

	go func() {
		<-ctx.Done()
		slog.Info("meta shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("meta shutdown error", "err", err)
		} else {
			slog.Info("meta shutdown complete")
		}
	}()

	slog.Info("meta API starting on :" + fmt.Sprintf("%d", m.port))
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("meta API server error", "err", err)
	}
	return err
}

// Stop shuts down the Meta service
func (m *Meta) Stop() {
	// Implement shutdown logic here
}
