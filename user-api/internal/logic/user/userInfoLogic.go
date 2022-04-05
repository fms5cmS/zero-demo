package user

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"zero-demo/user-api/internal/svc"
	"zero-demo/user-api/internal/types"
	"zero-demo/user-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) UserInfoLogic {
	return UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	fmt.Println("client get user info ---------> ")
	// 走 rpc 服务的方式
	princess, err := l.svcCtx.UserRpcClient.GetUserInfo(l.ctx, &pb.GetUserInfoReq{Id: req.UserId})
	if err != nil {
		return nil, err
	}

	// 走数据库的方式
	// 1. 测试 log 的 Encoding 配置
	//if err = l.one(); err != nil {
	//	logx.Errorf("err: %+v", err)
	//}
	// 2. 测试 log 的 Level 配置
	//logx.Error("test error level log")
	// 3. 查询数据库
	//princess, err := l.svcCtx.PrincessModel.FindOne(req.UserId)
	//if err != nil && !errors.Is(err, model.ErrNotFound) {
	//	return nil, errors.New("查询失败")
	//}
	//if princess == nil {
	//	return nil, errors.New("不存在")
	//}
	resp = &types.UserInfoResp{UserId: req.UserId, Nickname: princess.UserModel.Nickname}
	return
}

func (l *UserInfoLogic) one() error {
	return l.Two()
}

func (l *UserInfoLogic) Two() error {
	return errors.Wrap(errors.New("test error"), "测试堆栈类型的日志")
}
