package web

import (
	"embed"
	"io/fs"
	"log/slog"
	"mime"
	"net/http"
	"path"
	"strings"
)

// uiDist embeds the built UI assets from the Vite build output.
// Vite is configured to output directly into internal/web/dist.
//
//go:embed dist/**
var uiDist embed.FS

// UIFileSystem returns an fs.FS rooted at ui/dist for serving via http.FS.
func UIFileSystem() fs.FS {
	sub, err := fs.Sub(uiDist, "dist")
	if err != nil {
		// If this happens, it means the dist directory wasn't built or paths moved.
		// Log and fall back to the root which will 404 on access.
		slog.Error("failed to sub FS for UI dist", "err", err)
		return uiDist
	}
	return sub
}

// SPAHandler serves files from the embedded UI dist directory and falls back to index.html
// for unknown routes to support client-side routing.
func SPAHandler() http.Handler {
	distFS := http.FS(UIFileSystem())
	fileServer := http.FileServer(distFS)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Clean the path and ensure it does not attempt to escape the dist root.
		p := path.Clean(r.URL.Path)
		if p == "/" {
			// Serve index.html explicitly
			serveIndex(w, r)
			return
		}

		// Attempt to open the requested file
		f, err := distFS.Open(strings.TrimPrefix(p, "/"))
		if err == nil {
			// If it exists and is a file, serve it directly
			_ = f.Close()
			// Set a sensible Content-Type if possible
			if ctype := mime.TypeByExtension(path.Ext(p)); ctype != "" {
				w.Header().Set("Content-Type", ctype)
			}
			fileServer.ServeHTTP(w, r)
			return
		}

		// Fallback: serve index.html for SPA routes
		serveIndex(w, r)
	})
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	data, err := fs.ReadFile(UIFileSystem(), "index.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
