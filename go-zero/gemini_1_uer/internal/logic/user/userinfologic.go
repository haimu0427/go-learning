// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"encoding/json"

	"gemini_1_uer/internal/svc"
	"gemini_1_uer/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo() (resp *types.UserInfoResp, err error) {
	// 1. 从 context 中获取 userId
	// "userId" 这个 key 必须和你生成 token 时 claims["userId"] 的 key 一致
	// go-zero 解析出来的数字默认是 json.Number 类型
	userIdNumber := l.ctx.Value("userId").(json.Number)
	userId, _ := userIdNumber.Int64()

	// 2. 拿着 userId 去调用 RPC 获取详细信息
	// (这里你需要去 mall-rpc 加一个 GetUserInfo 的接口，为了演示简单，我们先模拟返回)

	return &types.UserInfoResp{
		Id:     userId,
		Name:   "Authenticated User",
		Mobile: "13800000000",
	}, nil
}
