// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"user-api/internal/svc"
	"user-api/internal/types"
	"user/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	userResp, err := l.svcCtx.UserRpc.GetUser(l.ctx, &userclient.GetUserRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &types.UserInfoResp{
		Id:       userResp.Id,
		Username: userResp.Name,
		Email:    userResp.Name + "@example.com",
	}, nil
}
