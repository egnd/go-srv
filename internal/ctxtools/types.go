package ctxtools

type ctxKey int

const (
	ctxLogger ctxKey = iota
	ctxReqID
)
