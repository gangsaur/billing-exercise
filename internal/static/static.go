package static

type contextKey string

const (
	RequestIdKey = contextKey("X-Request-ID")
)
