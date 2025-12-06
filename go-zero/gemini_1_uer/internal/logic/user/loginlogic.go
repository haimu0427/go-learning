// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2
// todo: 你的工作主要集中在 internal/logic
package user

import (
	"context"
	"errors"
	"gemini_1_uer/internal/model"
	"gemini_1_uer/internal/svc"
	"gemini_1_uer/internal/types"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 获取 JWT Token 的辅助函数
func (l *LoginLogic) getJwtToken(secretKey string, iat, seconds int64, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId // 把 userId 塞进 token 里
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 1. 使用 ServiceContext 中的 UserModel 查询数据库
	// goctl 自动生成了 FindOneByUsername，因为我们在 SQL 里加了唯一索引
	userInfo, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)

	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.New("用户不存在")
		}
		return nil, err // 查询出错（如数据库挂了）
	}

	// 2. 校验密码 (这里简单演示明文比对，生产环境请使用 bcrypt)
	if userInfo.Password != req.Password {
		return nil, errors.New("密码错误")
	}

	// 3. 登录成功，返回数据
	return &types.LoginResp{
		Id:       userInfo.Id,
		Name:     userInfo.Username,
		Token:    "mock_token_" + req.Username, // 下一阶段我们再讲 JWT
		ExpireAt: "2025-12-31",
	}, nil
}
