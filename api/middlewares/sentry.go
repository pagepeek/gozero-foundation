package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/getsentry/sentry-go"
	sentryx "github.com/pagepeek/gozero-foundation/pkg/sentry"
	"github.com/zeromicro/go-zero/rest"
)

func NewSentryMiddleware(opt sentryx.SentryOption) rest.Middleware {
	err := sentryx.Setup(opt)
	if err != nil {
		panic(err)
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			handle(w, r, next, opt)
		}
	}
}

func handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc, opt sentryx.SentryOption) {
	ctx := r.Context()
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub().Clone()
	}

	transaction := sentry.StartTransaction(
		sentry.SetHubOnContext(ctx, hub),
		fmt.Sprintf("%s %s", r.Method, r.URL.Path),
		[]sentry.SpanOption{
			sentry.ContinueTrace(hub, r.Header.Get(sentry.SentryTraceHeader), r.Header.Get(sentry.SentryBaggageHeader)),
			sentry.WithOpName("http.server"),
			sentry.WithTransactionName(fmt.Sprintf("%s %s", r.Method, r.URL.Path)),
			sentry.WithTransactionSource(sentry.SourceRoute),
			sentry.WithSpanOrigin(sentry.SpanOriginManual),
			func(s *sentry.Span) {
				s.SetData("http.request.method", r.Method)
				s.SetData("http.request.url", r.URL.Path)
			},
		}...,
	)

	hub.Scope().SetRequest(r)

	defer transaction.Finish()
	defer func() {
		if p := recover(); p != nil {
			eventID := hub.Recover(
				fmt.Errorf("%+v\n\n%s", p, strings.TrimSpace(string(debug.Stack()))),
			)

			if eventID != nil {
				hub.Flush(opt.Timeout)
			}

			if opt.Repanic {
				panic(p)
			}

			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	next(w, r)
}
