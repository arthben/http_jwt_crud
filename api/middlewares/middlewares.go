package middlewares

import (
	"bytes"
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/arthben/http_jwt_crud/internal/logging"
	"github.com/google/uuid"
)

const KeyRequestID = "request_id"

type Middleware func(logger *slog.Logger, next http.Handler) http.HandlerFunc

func MiddlewareChain(middleware ...Middleware) Middleware {
	return func(logger *slog.Logger, next http.Handler) http.HandlerFunc {
		// middleware read in reverse mode.
		for i := len(middleware) - 1; i >= 0; i-- {
			next = middleware[i](logger, next)
		}

		return next.ServeHTTP
	}
}

func ContextRequestID(logger *slog.Logger, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// use header X-Request-ID from client. otherwise, generate using UUID4
		reqID := r.Header.Get("X-Request-ID")
		if len(reqID) == 0 {
			reqID = uuid.New().String()
			r.Header.Set("X-Request-ID", reqID)
		}

		embededRequestCtx := context.WithValue(r.Context(), KeyRequestID, reqID)

		// embed to object http.Request
		next.ServeHTTP(w, r.WithContext(embededRequestCtx))
	}
}

func CORS(logger *slog.Logger, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header
		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Credentials", "true")
		header.Set("Access-Control-Allow-Headers", "*")
		header.Set("Access-Control-Allow-Methods", "*")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func RequestLogger(logger *slog.Logger, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		reqID, _ := r.Context().Value(KeyRequestID).(string)

		reqLoggerInfo := logger.With(
			slog.String("request_id", reqID),
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.String("query", r.URL.Query().Encode()),
		)

		// read body
		reqBody, err := readBody(r)
		if err != nil {
			reqLoggerInfo.Error(
				"Error Reading Request Body",
				slog.String("error", err.Error()),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// tap response from http.ResponseWriter
		newWriter := wrapResponseWriter(w)

		// embed logger slog to request when needed on handlers
		embededLoggerCtx := logging.Attach(reqLoggerInfo, r.Context())
		r = r.WithContext(embededLoggerCtx)

		next.ServeHTTP(newWriter, r)

		latency := float64(time.Since(t).Seconds())
		// write log
		go func() {
			status := newWriter.ResponseStatus
			respBody := newWriter.ResponseBody.String()
			reqLoggerInfo.Info(
				"",
				slog.Float64("latency", latency),
				slog.Int("status", status),
				slog.String("request", string(reqBody)),
				slog.String("response", strings.Replace(respBody, "\n", "", 1)),
			)
		}()
	}
}

func readBody(r *http.Request) ([]byte, error) {
	// baca isi body untuk bisa disimpan di log.
	// setelah baca isi log, rewrite lagi body nya agar saat next.ServerHTTP body ada isinya
	byteBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Gagal Baca Payload Requst %s\n", err.Error())
		return nil, err
	}

	bodyToAttach := bytes.NewBuffer([]byte{})
	bodyCopy := bytes.NewBuffer([]byte{})
	payloadCopy := io.MultiWriter(bodyCopy, bodyToAttach)

	_, err = payloadCopy.Write(byteBody)
	if err != nil {
		log.Printf("Gagal Copy Payload %s\n", err.Error())
		return nil, err
	}

	r.Body = io.NopCloser(bodyToAttach)

	return bodyCopy.Bytes(), nil
}

type httpResponseWrapper struct {
	wrapped        http.ResponseWriter
	ResponseStatus int
	ResponseBody   *bytes.Buffer
}

// Kita perlu ini supaya bisa dapet copy-an body & status code nya
func wrapResponseWriter(w http.ResponseWriter) *httpResponseWrapper {
	return &httpResponseWrapper{
		wrapped:      w,
		ResponseBody: bytes.NewBuffer([]byte{}),
	}
}

// Header implements http.ResponseWriter.
func (w *httpResponseWrapper) Header() http.Header {
	return w.wrapped.Header()
}

// Write implements http.ResponseWriter.
func (w *httpResponseWrapper) Write(p []byte) (int, error) {
	writeTargets := io.MultiWriter(w.wrapped, w.ResponseBody)
	return writeTargets.Write(p)
}

// WriteHeader implements http.ResponseWriter.
func (w *httpResponseWrapper) WriteHeader(statusCode int) {
	w.ResponseStatus = statusCode
	w.wrapped.WriteHeader(statusCode)
}
