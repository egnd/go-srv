package ctxtools

import "context"

func AddReqID(ctx context.Context, reqID string) context.Context {
	if reqID == "" {
		return ctx
	}

	return context.WithValue(ctx, ctxReqID, reqID)
}

func GetReqID(ctx context.Context) string {
	return ctx.Value(ctxReqID).(string)
}
