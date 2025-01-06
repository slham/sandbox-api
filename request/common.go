package request

import "context"

const (
	reqCtx = "request-context"
)

type RequestContext struct {
	Stop         bool
	RequestID    string
	UserID       string
	ClientUserID string
	Roles        []string
}

func WithRequestContext(ctx context.Context, rc *RequestContext) context.Context {
	return context.WithValue(ctx, reqCtx, rc)
}

func GetRequestContext(ctx context.Context) *RequestContext {
	if rc, ok := ctx.Value(reqCtx).(*RequestContext); ok {
		return rc
	}
	return nil
}

func SetStop(ctx context.Context) context.Context {
	rc := GetRequestContext(ctx)
	if rc == nil {
		rc = &RequestContext{Stop: true}
	}

	rc.Stop = true

	return WithRequestContext(ctx, rc)
}

func GetStop(ctx context.Context) bool {
	rc := GetRequestContext(ctx)
	if rc == nil {
		return false
	}
	return rc.Stop
}
