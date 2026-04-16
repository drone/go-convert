package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

// Server wraps an http.Server with structured logging and graceful shutdown.
type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

// NewServer builds a Server that listens on the given port and applies all
// middleware (recovery → requestID → logging) before the conversion handlers.
func NewServer(port, maxBatchSize int, maxYAMLBytes int64, logger *slog.Logger) *Server {
	h := NewHandler(maxBatchSize, maxYAMLBytes)
	mux := buildMux(h)
	wrapped := chain(mux,
		recoveryMiddleware(logger),
		requestIDMiddleware(),
		loggingMiddleware(logger),
	)
	return &Server{
		logger: logger,
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      wrapped,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

// Start listens and serves HTTP, blocking until the server stops.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return fmt.Errorf("failed to bind %s: %w", s.httpServer.Addr, err)
	}
	s.logger.Info("server listening", "addr", s.httpServer.Addr)
	return s.httpServer.Serve(ln)
}

// Stop gracefully shuts down the server using the provided context deadline.
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// buildMux registers all routes. Go 1.22+ ServeMux supports "METHOD /path" patterns.
func buildMux(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.Healthz)
	mux.HandleFunc("POST /api/v1/convert/pipeline", h.ConvertPipeline)
	mux.HandleFunc("POST /api/v1/convert/template", h.ConvertTemplate)
	mux.HandleFunc("POST /api/v1/convert/input-set", h.ConvertInputSet)
	mux.HandleFunc("POST /api/v1/convert/batch", h.ConvertBatch)
	mux.HandleFunc("POST /api/v1/checksum", h.ComputeChecksum)
	return mux
}

// ---- middleware helpers ----

type middlewareFunc func(http.Handler) http.Handler

// chain wraps h with each middleware, outermost first.
// e.g. chain(h, a, b, c) produces: a(b(c(h)))
func chain(h http.Handler, mw ...middlewareFunc) http.Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		h = mw[i](h)
	}
	return h
}

// recoveryMiddleware catches panics, logs them, and returns HTTP 500.
func recoveryMiddleware(logger *slog.Logger) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered",
						"panic", fmt.Sprintf("%v", rec),
						"method", r.Method,
						"path", r.URL.Path,
					)
					writeError(w, http.StatusInternalServerError,
						"INTERNAL_ERROR", "an unexpected error occurred", nil)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

type contextKey string

const requestIDCtxKey contextKey = "requestID"

// requestIDMiddleware generates a random request ID, sets it on the response
// header, and stores it in the request context.
func requestIDMiddleware() middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := newRequestID()
			w.Header().Set("X-Request-ID", id)
			ctx := context.WithValue(r.Context(), requestIDCtxKey, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func newRequestID() string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return fmt.Sprintf("%x", b)
}

// responseWriter wraps http.ResponseWriter to capture the written status code.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.wroteHeader = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

// loggingMiddleware logs each request with method, path, status, and latency.
func loggingMiddleware(logger *slog.Logger) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rw, r)
			logger.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.status,
				"latency_ms", time.Since(start).Milliseconds(),
				"request_id", r.Context().Value(requestIDCtxKey),
			)
		})
	}
}
