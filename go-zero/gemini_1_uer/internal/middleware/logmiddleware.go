package middleware

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogMiddleware struct {
}

func NewLogMiddleware() *LogMiddleware {
	return &LogMiddleware{}
}

func (m *LogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Info("=== 进入中间件: 收到请求 ", r.RequestURI)
		// 放行，执行后续逻辑
		next(w, r)
		logx.Info("=== 离开中间件: 请求处理完毕")
	}
}
