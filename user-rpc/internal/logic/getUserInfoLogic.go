package logic

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"

	"zero-demo/user-rpc/internal/svc"
	"zero-demo/user-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *pb.GetUserInfoReq) (*pb.GetUserInfoResp, error) {
	// 获取客户端传入的 metadata，注意这里使用的函数
	if md, exists := metadata.FromIncomingContext(l.ctx); exists {
		names := md.Get("username")
		if len(names) > 0 {
			fmt.Printf("username ===========> %s\n", names[0])
		}
	}

	fmt.Println("server get user info ---------> ")
	userMap := map[int64]string{
		1: "佩克莉姆 from  rpc",
		2: "可可萝 from rpc",
		3: "凯露 from rpc",
	}
	nickname := "unknown"
	if name, exists := userMap[in.Id]; exists {
		nickname = name
	}
	return &pb.GetUserInfoResp{UserModel: &pb.UserModel{Id: in.Id, Nickname: nickname}}, nil
}
