package sentry

import (
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"google.golang.org/grpc/codes"
)

// The identifier of the Go-Zero SDK.
const (
	sdkIdentifier = "sentry.go-zero"
)

type SentryOption struct {
	sentry.ClientOptions
	Dsn              string
	EnableTracing    bool
	TracesSampleRate float64
	TracesSampler    sentry.TracesSampler
	Repanic          bool
	Timeout          time.Duration
	IgnoreCodes      []codes.Code
}

func Setup(opt SentryOption) error {
	if opt.ServerName == "" {
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "gozero-unknown"
		}

		opt.ServerName = hostname
	}

	opt.ClientOptions.Dsn = opt.Dsn
	opt.ClientOptions.EnableTracing = opt.EnableTracing
	opt.ClientOptions.TracesSampleRate = opt.TracesSampleRate
	opt.ClientOptions.TracesSampler = opt.TracesSampler
	opt.ClientOptions.ServerName = opt.ServerName

	if err := sentry.Init(opt.ClientOptions); err != nil {
		return err
	}

	hub := sentry.CurrentHub()
	client := hub.Client()
	if client != nil {
		client.SetSDKIdentifier(sdkIdentifier)
	}

	if opt.Timeout == 0 {
		opt.Timeout = time.Second * 2
	}

	return nil
}
