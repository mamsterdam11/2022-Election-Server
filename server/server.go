package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/mamsterdam11/Snowflake-News/news"
)

const (
	serverAddr   = ":8282"
	shutdownTime = 10 * time.Second // max time waiting for server to shutdown

	recentNewsCount = 3 // number of recent news stories to provide
)

type WebServer struct {
	srv       *http.Server
	collector *news.NewsCollector
}

func NewWebServer(nc *news.NewsCollector) *WebServer {
	return &WebServer{
		collector: nc,
		srv: &http.Server{
			Addr: serverAddr,
		},
	}
}

// Serve inits an HTTP server, listening at serverAddr for requests
// on the /news endpoint. The server attempts to gracefully shut down
// upon context cancellation
func (w *WebServer) Serve(ctx context.Context) {
	mux := http.NewServeMux()
	mux.Handle("/news", w.handleNews())
	w.srv.Handler = mux

	go func() {
		if err := w.srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe failed: %v", err)
		}
	}()
	log.Printf("Server listening on %s", serverAddr)

	<-ctx.Done()

	w.shutdown()
	log.Printf("Server shutdown gracefully")
}

func (w *WebServer) handleNews() http.HandlerFunc {
	recentNews := w.collector.RecentNews(recentNewsCount)
	formattedNews := news.FormatNews(recentNews)

	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(formattedNews))
	}
}

func (w *WebServer) shutdown() {
	ctxShutdown, cancel := context.WithTimeout(context.Background(), shutdownTime)
	defer cancel()

	if err := w.srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}
}
