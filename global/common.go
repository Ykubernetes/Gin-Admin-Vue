package global

import (
	"gitee.com/go-server/conf"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	Config  *conf.Config
	DB      *gorm.DB
	Log     *logrus.Logger
	RedisDB *redis.Client
)
