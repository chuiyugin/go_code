package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf

	ShortUrlDB struct {
		DSN string
	}

	SequenceDB struct {
		DSN string
	}

	ShortUrlBlackList []string

	ShortDoamin string
}
