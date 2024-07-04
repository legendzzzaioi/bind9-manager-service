package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	DataSource string
	BindPath   string
	JwtAuth    struct {
		AccessSecret string
		AccessExpire int
	}
}
