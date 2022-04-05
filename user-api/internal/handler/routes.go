// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	user "zero-demo/user-api/internal/handler/user"
	"zero-demo/user-api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.TestMiddleware},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/user/create",
					Handler: user.UserCreateHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/user/info/:userId",
					Handler: user.UserInfoHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/userapi/v1"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.TestMiddleware2},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/user/infoupdate",
					Handler: user.UserInfoUpateHandler(serverCtx),
				},
			}...,
		),
	)
}
