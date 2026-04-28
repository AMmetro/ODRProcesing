package core_http_middleware

import (
	"net/http"
	"time"

	core_logger "github.com/AMmetro/ODRProcesing/internal/core/logger"
	core_http_response "github.com/AMmetro/ODRProcesing/internal/core/transport/http/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const requestIdHeader = "X-Request-ID"

func RequestID() Middleware {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			requestId := r.Header.Get(requestIdHeader)
			if requestId == "" {
				requestId = uuid.NewString()
			}

			r.Header.Set(requestIdHeader, requestId)
			w.Header().Set(requestIdHeader, requestId) /*  */

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIdHeader)

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			// ctx := context.WithValue(r.Context(), core_logger.LoggerContextKey, l)

			ctx := core_logger.ToContext(r.Context(), l)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			rw := core_http_response.NewResponseWriter(w)
			before := time.Now()
			log.Debug(
				">>> incoming HTTP request",
				zap.String("http_method", r.Method),
				zap.Time("time", before.UTC()),
			)
			next.ServeHTTP(rw, r)
			log.Debug(
				"<<< done HTTP request",
				zap.Int("Status_cod", rw.GetStatusCode()),
				zap.Duration("latency", time.Now().Sub(before)),
			)
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(p,
						"during handle HTTP request got unexpected panic",
					)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
