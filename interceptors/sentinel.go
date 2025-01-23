package interceptors

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// The identifier of the Go-Zero SDK.
const (
	sdkIdentifier = "sentry.go-zero"
)

type SentinelOption struct {
	Dsn              string
	Repanic          bool
	Timeout          time.Duration
	IgnoreCodes      []codes.Code
	TracesSampleRate float64
	TracesSampler    sentry.TracesSampler
}

type SentinelInterceptor struct {
	reporter    *sentry.Client
	timeout     time.Duration
	repanic     bool
	ignoreCodes []codes.Code
}

func NewSentinelInterceptor(opt SentinelOption) *SentinelInterceptor {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "center-unknown"
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn:              opt.Dsn,
		EnableTracing:    true,
		ServerName:       hostname,
		TracesSampleRate: opt.TracesSampleRate,
		TracesSampler:    opt.TracesSampler,
	})
	if err != nil {
		panic(err)
	}

	hub := sentry.CurrentHub()
	client := hub.Client()
	if client != nil {
		client.SetSDKIdentifier(sdkIdentifier)
	}

	if opt.Timeout == 0 {
		opt.Timeout = time.Second * 2
	}

	return &SentinelInterceptor{reporter: client, repanic: opt.Repanic, timeout: opt.Timeout, ignoreCodes: opt.IgnoreCodes}
}

func (i *SentinelInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		var resp any
		err := i.handle(ctx, info.FullMethod, func() error {
			var err error
			resp, err = handler(ctx, req)

			return err
		})

		return resp, err
	}
}

func (i *SentinelInterceptor) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, s grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := i.handle(s.Context(), info.FullMethod, func() error {
			return handler(srv, s)
		})

		return err
	}
}

func (i *SentinelInterceptor) handle(ctx context.Context, name string, next func() error) (err error) {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub().Clone()
	}

	transaction := sentry.StartTransaction(sentry.SetHubOnContext(ctx, hub), name,
		[]sentry.SpanOption{
			sentry.WithOpName("rpc.server"),
			sentry.WithTransactionSource(sentry.SourceCustom),
			sentry.WithSpanOrigin(sentry.SpanOriginManual),
			func(s *sentry.Span) {
				s.SetData("rpc.full_method", name)
			},
		}...,
	)

	defer transaction.Finish()
	defer func() {
		if p := recover(); p != nil {
			err = i.report(hub, fmt.Errorf("%+v\n\n%s", p, strings.TrimSpace(string(debug.Stack()))))
		}
	}()

	err = next()
	if err != nil {
		if stErr, ok := err.(interface{ GRPCStatus() *status.Status }); !ok {
			hub.CaptureException(err)
		} else {
			status := stErr.GRPCStatus()
			// 业务错误不做记录
			if lo.Contains(i.ignoreCodes, status.Code()) {
				return
			}

			switch status.Code() {
			case codes.OK:
				transaction.Status = sentry.SpanStatusOK
			case codes.Canceled:
				transaction.Status = sentry.SpanStatusCanceled
			case codes.Unknown:
				transaction.Status = sentry.SpanStatusUnknown
			case codes.InvalidArgument:
				transaction.Status = sentry.SpanStatusInvalidArgument
			case codes.DeadlineExceeded:
				transaction.Status = sentry.SpanStatusDeadlineExceeded
			case codes.NotFound:
				transaction.Status = sentry.SpanStatusNotFound
			case codes.AlreadyExists:
				transaction.Status = sentry.SpanStatusAlreadyExists
			case codes.PermissionDenied:
				transaction.Status = sentry.SpanStatusPermissionDenied
			case codes.ResourceExhausted:
				transaction.Status = sentry.SpanStatusResourceExhausted
			case codes.FailedPrecondition:
				transaction.Status = sentry.SpanStatusFailedPrecondition
			case codes.Aborted:
				transaction.Status = sentry.SpanStatusAborted
			case codes.OutOfRange:
				transaction.Status = sentry.SpanStatusOutOfRange
			case codes.Unimplemented:
				transaction.Status = sentry.SpanStatusUnimplemented
			case codes.Internal:
				transaction.Status = sentry.SpanStatusInternalError
			case codes.Unavailable:
				transaction.Status = sentry.SpanStatusUnavailable
			case codes.DataLoss:
				transaction.Status = sentry.SpanStatusDataLoss
			case codes.Unauthenticated:
				transaction.Status = sentry.SpanStatusUnauthenticated
			default:
				transaction.Status = sentry.SpanStatusUndefined
			}

			hub.CaptureException(stErr.GRPCStatus().Err())
		}

		hub.Flush(i.timeout)
	}

	return
}

func (i *SentinelInterceptor) report(hub *sentry.Hub, err error) error {
	eventID := hub.Recover(err)

	if eventID != nil {
		hub.Flush(i.timeout)
	}

	if i.repanic {
		panic(err)
	}

	return status.Errorf(codes.Internal, "panic: %v", err)
}
