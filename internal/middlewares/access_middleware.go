// Copyright © ivanlobanov. All rights reserved.
package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/cantylv/service-happy-birthday/internal/utils/myconstants"
	"github.com/cantylv/service-happy-birthday/internal/utils/recorder"
	"github.com/satori/uuid"
	"go.uber.org/zap"
)

type AccessLogStart struct {
	UserAgent      string
	RealIp         string
	ContentLength  string
	URI            string
	Method         string
	StartTimeHuman string
	RequestId      string
	Logger         *zap.Logger
}

type AccessLogEnd struct {
	LatencyMs      int64
	ResponseSize   string // in bytes
	ResponseStatus int
	EndTimeHuman   string
	RequestId      string
	Logger         *zap.Logger
}

// formatTime
// Returns format time. For formatting we should use this datetimetz `Mon Jan 2 15:04:05 MST 2006`.
func formatTime(t time.Time) string {
	return t.Format("02.01.2006 15:04:05 UTC-07") // for me, it's more readable format
}

// Access
// Middleware that logs the start and end of request handling.
func Access(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := zap.Must(zap.NewProduction())
		requestId := uuid.NewV4().String()
		ctx := context.WithValue(context.Background(), myconstants.RequestId, requestId)
		r = r.WithContext(ctx)

		rec := recorder.NewResponseWriter(w)

		timeNow := time.Now()
		// nginx will proxy headers, like "User-Agent", "X-Real-IP", "Content-Length"
		startLog := AccessLogStart{
			UserAgent:      r.Header.Get("User-Agent"),
			RealIp:         r.Header.Get("X-Real-IP"),
			ContentLength:  r.Header.Get("Content-Length"),
			URI:            r.RequestURI,
			Method:         r.Method,
			StartTimeHuman: formatTime(timeNow),
			RequestId:      requestId,
			Logger:         logger,
		}
		LogInitRequest(startLog)

		h.ServeHTTP(rec, r)

		timeEnd := time.Now()
		endLog := AccessLogEnd{
			LatencyMs:      timeEnd.Sub(timeNow).Milliseconds(),
			ResponseSize:   w.Header().Get("Content-Length"),
			ResponseStatus: rec.StatusCode,
			EndTimeHuman:   formatTime(timeEnd),
			RequestId:      requestId,
			Logger:         logger,
		}
		LogEndRequest(endLog)
	})
}

// LogInitRequest
// Logs user-agent, real-ip and etc..
func LogInitRequest(data AccessLogStart) {
	data.Logger.Info("request-id "+data.RequestId,
		zap.String("user-agent", data.UserAgent),
		zap.String("real-ip", data.RealIp),
		zap.String("content-length", data.ContentLength),
		zap.String("uri", data.URI),
		zap.String("method", data.Method),
		zap.String("start-time-human", data.StartTimeHuman),
	)
}

// LogEndRequest
// Logs latency in ms, response size and etc..
func LogEndRequest(data AccessLogEnd) {
	data.Logger.Info("request-id "+data.RequestId,
		zap.Int64("Latensy-MS", data.LatencyMs),
		zap.String("Response-Size", data.ResponseSize),
		zap.Int("Response-Status", data.ResponseStatus),
		zap.String("End-Time-Human", data.EndTimeHuman),
	)
}
