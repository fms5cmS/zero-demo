package user

import (
	"context"

	"zero-demo/user-api/internal/svc"
	"zero-demo/user-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoUpateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoUpateLogic(ctx context.Context, svcCtx *svc.ServiceContext) UserInfoUpateLogic {
	return UserInfoUpateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoUpateLogic) UserInfoUpate(req types.UserInfoUpdateReq) (resp *types.UserInfoUpdateResp, err error) {
	// todo: add your logic here and delete this line

	return
}
