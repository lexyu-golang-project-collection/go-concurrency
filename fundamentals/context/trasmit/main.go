package main

import (
	"context"
	"fmt"
)

type ctxKey string

const (
	appIDKey     = ctxKey("app_id")
	traceIDKey   = ctxKey("trace_id")
	sessionIDKey = ctxKey("session_id")
	userIDKey    = ctxKey("user_id")
)

// RequestContext 封裝請求相關資訊
type RequestContext struct {
	AppID     string
	TraceID   string
	SessionID string
	UserID    string
}

// Context tree
func buildContextTree(baseCtx context.Context, req RequestContext) context.Context {
	// 1 layer
	appCtx := context.WithValue(baseCtx, appIDKey, req.AppID)

	// 2 layer
	traceCtx := context.WithValue(appCtx, traceIDKey, req.TraceID)

	// 3 layer
	sessionCtx := context.WithValue(traceCtx, sessionIDKey, req.SessionID)

	// 4 layer
	userCtx := context.WithValue(sessionCtx, userIDKey, req.UserID)

	return userCtx
}

// multi-layer process logic
func processRequest(ctx context.Context) {
	if appID := ctx.Value(appIDKey); appID != nil {
		fmt.Printf("App Layer: Processing request for app %v\n", appID)
	}

	if traceID := ctx.Value(traceIDKey); traceID != nil {
		fmt.Printf("Trace Layer: Request traced with ID %v\n", traceID)
	}

	if sessionID := ctx.Value(sessionIDKey); sessionID != nil {
		fmt.Printf("Session Layer: Active session %v\n", sessionID)
	}

	if userID := ctx.Value(userIDKey); userID != nil {
		fmt.Printf("User Layer: Processing for user %v\n", userID)
	}
}

func main() {
	baseCtx := context.Background()

	// 模擬請求資訊
	req := RequestContext{
		AppID:     "app-001",
		TraceID:   "trace-123",
		SessionID: "session-456",
		UserID:    "user-789",
	}

	// context tree
	ctx := buildContextTree(baseCtx, req)

	processRequest(ctx)
}
