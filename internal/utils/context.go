package utils

import (
	"context"
	"net/url"
)

const (
	CtxRequestID     = "XRequestID"
	CtxRequestURL    = "url"
	CtxRequestMethod = "method"
)

func ExtractAttrs(ctx context.Context) []any {
	var attrs []any
	if v, ok := ctx.Value(CtxRequestID).(string); ok {
		attrs = append(attrs, CtxRequestID, v)
	}
	if val := ctx.Value(CtxRequestURL); val != nil {
		switch v := val.(type) {
		case string:
			attrs = append(attrs, CtxRequestURL, v)
		case *url.URL:
			attrs = append(attrs, CtxRequestURL, v.String())
		}
	}
	if v, ok := ctx.Value(CtxRequestMethod).(string); ok {
		attrs = append(attrs, CtxRequestMethod, v)
	}

	return attrs
}
