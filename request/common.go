package request

import "context"

const (
	reqCtx = "request-context"
)

type RequestContext struct {
	RequestID string
	UserID    string
	Roles     []string
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
