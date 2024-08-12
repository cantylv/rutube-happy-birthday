package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/satori/uuid"
	"go.uber.org/zap"
)

type AccessLogStart struct {
	UserAgent      string
	RealIp         string
	ContentLength  uint32
	URI            string
	Method         string
	StartTimeHuman string
	RequestId      string
}

type AccessLogEnd struct {
	LatencyMs      uint32
	EndTimeHuman   string
	RequestTimeMs  uint32
	RequestId      string
	ResponseStatus int
}

func Access(h http.Handler, r *http.Request) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := zap.Must(zap.NewProduction()).Sugar()
		requestId := uuid.NewV4().String()
		timeNow := time.Now().UTC()
		LogInitRequest(r, logger, timeNow, requestId)

		unauthId, err := functions.GetCookieUnauthIdValue(r)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			logger.Error(err.Error(), zap.String(cnst.RequestId, requestId))
			return
		}
		ctx := context.WithValue(r.Context(), cnst.UnauthIdCookieName, unauthId)
		ctx = context.WithValue(ctx, cnst.RequestId, requestId)
		r = r.WithContext(ctx)

		rec, ok := w.(*recorder.ResponseWriter)
		if !ok {
			functions.ErrorResponse(w, myerrors.InternalServerRu, http.StatusInternalServerError)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func LogInitRequest(r *http.Request, logger *zap.SugaredLogger, timeNow time.Time, requestId string) {
	msg := fmt.Sprintf("init request %s", requestId)
	startLog := &AccessLogStart{
		UserAgent:     r.UserAgent(),
		RealIp:        r.Header.Get("X-Real-IP"),
		ContentLength: r.ContentLength,
		URI:           r.RequestURI,
		Method:        r.Method,
		StartTime:     timeNow.Format(cnst.Timestamptz),
	}

	logger.Info(msg,
		zap.String("user_agent", startLog.UserAgent),
		zap.String("real_ip", startLog.RealIp),
		zap.Int64("content_length", startLog.ContentLength),
		zap.String("uri", startLog.URI),
		zap.String("method", startLog.Method),
		zap.String("start_time", startLog.StartTime),
		zap.String("request_id", requestId),
	)
}

func LogEndRequest(logger *zap.SugaredLogger, timeNow time.Time, requestId string, responseStatus int) {
	msg := fmt.Sprintf("request done %s", requestId)
	endLog := &AccessLogEnd{
		EndTime:        timeNow.Format(cnst.Timestamptz),
		LatencyHuman:   time.Since(timeNow).String(),
		LatencyMs:      time.Since(timeNow).String(),
		ResponseStatus: responseStatus,
	}
	logger.Info(msg,
		zap.String("end_time", endLog.EndTime),
		zap.String("latency_human", endLog.LatencyHuman),
		zap.String("latency_human_ms", endLog.LatencyMs),
		zap.String("request_id", requestId),
		zap.Int("response_status", responseStatus),
	)
}
