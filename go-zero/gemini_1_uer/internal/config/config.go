// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	// 对应 YAML 中的 Mysql
	Mysql struct {
		DataSource string
	}
	// 对应 YAML 中的 Cache
	Cache cache.CacheConf
}
