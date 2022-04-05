package user

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"zero-demo/user-api/model"

	"zero-demo/user-api/internal/svc"
	"zero-demo/user-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) UserCreateLogic {
	return UserCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserCreateLogic) UserCreate(req types.UserCreateReq) (resp *types.UserCreateResp, err error) {
	err = l.svcCtx.PrincessModel.TransactCtxSelf(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		princess := &model.Princess{
			Name:     req.Name,
			Nickname: req.Nickname,
		}
		result, err := l.svcCtx.PrincessModel.TransactInsert(ctx, session, princess)
		if err != nil {
			return err
		}
		princessID, _ := result.LastInsertId()
		princessData := &model.PrincessData{PrincessId: princessID, Data: "xxxxxxxxxxx"}
		if _, err = l.svcCtx.PrincessDataModel.TransactInsert(ctx, session, princessData); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("创建用户失败")
	}
	resp = &types.UserCreateResp{Flag: true}
	return
}
