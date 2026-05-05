package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	pb "github.com/drone/go-convert/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Server wraps both an HTTP and gRPC server with structured logging and graceful shutdown.
type Server struct {
	httpServer *http.Server
	grpcServer *grpc.Server
	grpcPort   int
	logger     *slog.Logger
}

// NewServer builds a Server that listens on the given HTTP port and gRPC port,
// applying middleware (recovery → requestID → logging) before the HTTP handlers.
func NewServer(port, grpcPort, maxBatchSize int, maxYAMLBytes int64, logger *slog.Logger) *Server {
	h := NewHandler(maxBatchSize, maxYAMLBytes)
	mux := buildMux(h)
	wrapped := chain(mux,
		recoveryMiddleware(logger),
		requestIDMiddleware(),
		loggingMiddleware(logger),
	)

	gs := grpc.NewServer(grpc.UnaryInterceptor(grpcLoggingInterceptor(logger)))
	pb.RegisterGoConvertServiceServer(gs, &GRPCHandler{})

	return &Server{
		logger:     logger,
		grpcPort:   grpcPort,
		grpcServer: gs,
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      wrapped,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

// StartGRPC starts the gRPC server in a goroutine; returns the listener error (if any).
func (s *Server) StartGRPC() error {
	addr := fmt.Sprintf(":%d", s.grpcPort)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to bind gRPC %s: %w", addr, err)
	}
	s.logger.Info("gRPC server listening", "addr", addr)
	return s.grpcServer.Serve(ln)
}

// Start listens and serves HTTP, blocking until the server stops.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return fmt.Errorf("failed to bind %s: %w", s.httpServer.Addr, err)
	}
	s.logger.Info("HTTP server listening", "addr", s.httpServer.Addr)
	return s.httpServer.Serve(ln)
}

// Stop gracefully shuts down both the gRPC and HTTP servers.
func (s *Server) Stop(ctx context.Context) error {
	s.grpcServer.GracefulStop()
	return s.httpServer.Shutdown(ctx)
}

// buildMux registers all routes. Compatible with Go 1.19+.
func buildMux(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", methodFilter("GET", h.Healthz))
	mux.HandleFunc("/api/v1/convert/pipeline", methodFilter("POST", h.ConvertPipeline))
	mux.HandleFunc("/api/v1/convert/template", methodFilter("POST", h.ConvertTemplate))
	mux.HandleFunc("/api/v1/convert/input-set", methodFilter("POST", h.ConvertInputSet))
	mux.HandleFunc("/api/v1/convert/trigger", methodFilter("POST", h.ConvertTrigger))
	mux.HandleFunc("/api/v1/convert/batch", methodFilter("POST", h.ConvertBatch))
	mux.HandleFunc("/api/v1/convert/expression", methodFilter("POST", h.ConvertExpression))
	mux.HandleFunc("/api/v1/checksum", methodFilter("POST", h.ComputeChecksum))
	return mux
}

// methodFilter wraps a handler to only accept requests with the specified HTTP method.
// Returns 405 Method Not Allowed for other methods.
func methodFilter(method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
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

// grpcLoggingInterceptor logs each gRPC call with method, status code, and latency.
func grpcLoggingInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		st, _ := status.FromError(err)
		logger.Info("grpc request",
			"method", info.FullMethod,
			"code", st.Code().String(),
			"latency_ms", time.Since(start).Milliseconds(),
		)
		return resp, err
	}
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
