package interceptors

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/getsentry/sentry-go"
	sentryx "github.com/pagepeek/gozero-foundation/pkg/sentry"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SentryInterceptor struct {
	opt sentryx.SentryOption
}

func NewSentryInterceptor(opt sentryx.SentryOption) *SentryInterceptor {
	err := sentryx.Setup(opt)
	if err != nil {
		panic(err)
	}

	return &SentryInterceptor{opt: opt}
}

func (i *SentryInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
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

func (i *SentryInterceptor) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, s grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := i.handle(s.Context(), info.FullMethod, func() error {
			return handler(srv, s)
		})

		return err
	}
}

func (i *SentryInterceptor) handle(ctx context.Context, name string, next func() error) (err error) {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub().Clone()
	}

	transaction := sentry.StartTransaction(sentry.SetHubOnContext(ctx, hub), name,
		[]sentry.SpanOption{
			sentry.WithOpName("rpc.server"),
			sentry.WithTransactionName(name),
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
			if lo.Contains(i.opt.IgnoreCodes, status.Code()) {
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

		hub.Flush(i.opt.Timeout)
	}

	return
}

func (i *SentryInterceptor) report(hub *sentry.Hub, err error) error {
	eventID := hub.Recover(err)

	if eventID != nil {
		hub.Flush(i.opt.Timeout)
	}

	if i.opt.Repanic {
		panic(err)
	}

	return status.Errorf(codes.Internal, "panic: %v", err)
}
