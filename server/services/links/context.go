package links

import "context"

const (
	ctxKey linksCtxKey = "links_ctx_key"
)

type linksCtxKey string

// NewContext places links to context
func NewContext(ctx context.Context, links interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := links.(*Service); !ok {
		links = newLinks(ctx)
	}

	return context.WithValue(ctx, ctxKey, links)
}

// FromContext returns links form context
func FromContext(ctx context.Context) *Service {
	if links, ok := ctx.Value(ctxKey).(*Service); ok {
		return links
	}

	return newLinks(ctx)
}
