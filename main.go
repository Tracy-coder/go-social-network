// Code generated by hertz generator.

package main

import (
	"context"
	"go-social-network/biz/logic"
	"go-social-network/configs"
	"go-social-network/data"
	"go-social-network/minio"
	"go-social-network/pkg/sentinel"
	"go-social-network/pkg/snowflake"
	"go-social-network/pkg/tracer"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/logger/accesslog"
	prometh "github.com/hertz-contrib/monitor-prometheus"
	hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"
	hertzSentinel "github.com/hertz-contrib/opensergo/sentinel/adapter"
)

func main() {
	data.InitData()
	configs.InitConfig()
	sentinel.InitSentinel()
	minio.InitMinio()
	snowflake.InitSnowFlake()
	// logger.InitLogger()

	go logic.Consumer4SyndicateStatus()

	p := tracer.InitTracer("go-social-network")
	defer p.Shutdown(context.Background())
	tracer, cfg := hertztracing.NewServerTracer()
	h := server.Default(tracer, server.WithTracer(
		prometh.NewServerTracer("127.0.0.1:9091", "/hertz",
			prometh.WithEnableGoCollector(true), // enable go runtime metric collector
		),
	))

	h.Use(hertztracing.ServerMiddleware(cfg))
	h.Use(hertzSentinel.SentinelServerMiddleware(
		hertzSentinel.WithServerBlockFallback(func(c context.Context, ctx *app.RequestContext) {
			ctx.JSON(http.StatusTooManyRequests, nil)
			ctx.Abort()
		}),
	))

	h.Use(accesslog.New(accesslog.WithFormat("[${time}] ${status} - ${latency} ${method} ${path} ${queryParams}")))
	h.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	register(h)

	h.Spin()
}
