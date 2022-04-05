package svc

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"zero-demo/user-api/internal/config"
	"zero-demo/user-api/internal/middleware"
	"zero-demo/user-api/model"
	"zero-demo/user-rpc/usercenter"
)

type ServiceContext struct {
	Config            config.Config
	TestMiddleware    rest.Middleware
	TestMiddleware2   rest.Middleware
	PrincessModel     model.PrincessModel
	PrincessDataModel model.PrincessDataModel
	// RPC client
	UserRpcClient usercenter.Usercenter
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:            c,
		PrincessModel:     model.NewPrincessModel(sqlx.NewMysql(c.DB.DataSource)),
		PrincessDataModel: model.NewPrincessDataModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		TestMiddleware:    middleware.NewTestMiddleware().Handle,
		TestMiddleware2:   middleware.NewTestMiddleware2Middleware().Handle,
		// 添加客户端拦截器
		UserRpcClient: usercenter.NewUsercenter(zrpc.MustNewClient(c.UserRpcConf, zrpc.WithUnaryClientInterceptor(TestClientInterceptor))),
	}
}

func TestClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// 向另一个服务传入 metadata
	md := metadata.New(map[string]string{"username": "fms5cmS"})
	// 注意这里使用的函数
	ctx = metadata.NewOutgoingContext(ctx, md)

	fmt.Println("发送前 ==========>")
	err := invoker(ctx, method, req, reply, cc, opts...)
	fmt.Println("发送后 ==========>")
	return err
}
