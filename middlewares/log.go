package middlewares

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
	"github.com/slham/sandbox-api/request"
)

type Level string
type ctxKey string

const (
	DEBUG      Level  = "DEBUG"
	INFO       Level  = "INFO"
	WARN       Level  = "WARN"
	ERROR      Level  = "ERROR"
	slogFields ctxKey = "slog_fields"
	requestID  string = "request_id"
)

var levels = map[Level]slog.Level{
	DEBUG: slog.LevelDebug,
	INFO:  slog.LevelInfo,
	WARN:  slog.LevelWarn,
	ERROR: slog.LevelError,
}

type ContextHandler struct {
	slog.Handler
}

// Handle adds contextual attributes to the Record before calling the underlying
// handler
func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}

	return h.Handler.Handle(ctx, r)
}

// AppendCtx adds an slog attribute to the provided context so that it will be
// included in any Record created with such context
func AppendCtx(parent context.Context, attr slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
		v = append(v, attr)
		return context.WithValue(parent, slogFields, v)
	}

	v := []slog.Attr{}
	v = append(v, attr)
	return context.WithValue(parent, slogFields, v)
}

func Initialize(lvl Level) bool {
	if level, ok := levels[lvl]; ok {
		var programLevel = new(slog.LevelVar)
		programLevel.Set(level)

		h := &ContextHandler{slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})}
		logger := slog.New(h)
		slog.SetDefault(logger)
		slog.Info("logging level set", "level", lvl)
		return ok
	} else {
		slog.Error("invalid logging level", "level", lvl)
		return ok
	}
}

// Initializes transaction logging
func LoggingInbound(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := initContext(r)
		r = r.WithContext(ctx)

		slog.DebugContext(ctx, "inbound", "headers", fmt.Sprintf("%v", r.Header))
		slog.DebugContext(ctx, "inbound", "remote", r.RemoteAddr, "method", r.Method, "url", r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

func LoggingOutbound(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		slog.DebugContext(ctx, "outbound", "headers", fmt.Sprintf("%v", w.Header()))
		h.ServeHTTP(w, r)
	})
}

func initContext(r *http.Request) context.Context {
	ctx := r.Context()
	vars := mux.Vars(r)
	clientUserID := vars["user_id"]
	if clientUserID != "" {
		ctx = AppendCtx(ctx, slog.String("client_user_id", clientUserID))
	}
	reqID := fmt.Sprintf("%s_%s", "req", ksuid.New())
	ctx = AppendCtx(ctx, slog.String(requestID, reqID))
	rc := request.GetRequestContext(ctx)
	if rc == nil {
		rc = &request.RequestContext{}
	}
	rc.RequestID = reqID
	rc.ClientUserID = clientUserID
	ctx = request.WithRequestContext(ctx, rc)

	return ctx
}
