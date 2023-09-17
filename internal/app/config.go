package app

import (
	"github.com/core-go/client"
	"github.com/core-go/core"
	mid "github.com/core-go/log/middleware"
	"github.com/core-go/log/zap"
	"github.com/core-go/sql"
)

type Config struct {
	Server     core.ServerConf     `mapstructure:"server"`
	Sql        sql.Config          `mapstructure:"sql"`
	Client     client.ClientConfig `mapstructure:"client"`
	Log        log.Config          `mapstructure:"log"`
	MiddleWare mid.LogConfig       `mapstructure:"middleware"`
	Action     *core.ActionConfig  `mapstructure:"action"`
}
