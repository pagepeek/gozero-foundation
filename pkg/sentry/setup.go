package sentry

import (
	"os"
	"time"

	"github.com/getsentry/sentry-go"
)

// The identifier of the Go-Zero SDK.
const (
	sdkIdentifier = "sentry.go-zero"
)

type SentryOption struct {
	sentry.ClientOptions
	Repanic bool
	Timeout time.Duration
}

func Setup(opt SentryOption) error {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "center-unknown"
	}

	opt.ClientOptions.ServerName = hostname

	err = sentry.Init(opt.ClientOptions)
	if err != nil {
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
