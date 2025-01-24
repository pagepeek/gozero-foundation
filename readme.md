# 下载

```bash
go get github.com/pagepeek/gozero-foundation@latest
```

## rpc服务使用sentry

```go
package main

import (
	"flag"
	"fmt"
	"time"

	"yourproject/rpc/internal/config"
	_ "yourproject/rpc/internal/ent/runtime"
	"yourproject/rpc/internal/svc"

	"github.com/getsentry/sentry-go"
	_ "github.com/lib/pq"
	sentryx "github.com/pagepeek/gozero-foundation/pkg/sentry"
	"github.com/pagepeek/gozero-foundation/rpc/interceptors"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func main() {
	var c config.Config
	conf.MustLoad("config.yaml", &c)

	ctx := svc.NewServiceContext(c)
	group := service.NewServiceGroup()
	defer group.Stop()

	sentinel := interceptors.NewSentryInterceptor(sentryx.SentryOption{
		Dsn: c.SentryDNS,
		// 启动追踪
		EnableTracing: true,
		// 追踪采样率
		TracesSampleRate: 1.0,
		// 自定义追踪采样
		TracesSampler: func(ctx sentry.SamplingContext) float64 {
			if ctx.Span.Name == "project.Server/HelloWorld" {
				return 0.0
			}

			return 1.0
		},
		// 捕获panic时是否重放,true会继续向上层Panic(err)
		Repanic: false,
		// flush sentry 时间
		Timeout: 2 * time.Second,
		// 忽略的错误码
		IgnoreCodes: []codes.Code{codes.NotFound, codes.PermissionDenied},
	})

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		// register your service
	})

	// register interceptor
	s.AddStreamInterceptors(sentinel.StreamInterceptor())
	s.AddUnaryInterceptors(sentinel.UnaryInterceptor())

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	group.Start()
}
```

## api服务使用sentry

```go
package main

import (
	"flag"
	"fmt"
	"time"

	"pagepeek/center/api/internal/config"
	"pagepeek/center/api/internal/errors"
	"pagepeek/center/api/internal/handler"
	"pagepeek/center/api/internal/svc"

	"github.com/getsentry/sentry-go"
	"github.com/pagepeek/gozero-foundation/api/middlewares"
	sentryx "github.com/pagepeek/gozero-foundation/pkg/sentry"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/center-api.yaml", "the config file")

func main() {
	var c config.Config
	conf.MustLoad("config.yaml", &c)

	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	server.Use(middlewares.NewSentryMiddleware(sentryx.SentryOption{
		Dsn: c.SentryDNS,
		// 启动追踪
		EnableTracing: true,
		// 追踪采样率
		TracesSampleRate: 1.0,
		// 自定义追踪采样
		TracesSampler: func(ctx sentry.SamplingContext) float64 {
			if ctx.Span.Name == "GET /hello_world" {
				return 0.0
			}

			return 1.0
		},
		// 捕获panic时是否重放,true会继续向上层Panic(err)
		Repanic: false,
		// flush sentry 时间
		Timeout: 2 * time.Second,
	}))

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	httpx.SetErrorHandler(errors.DefaultErrorHandler)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
```
