package config

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	GlobalConfig *AppConfig
	DB           *gorm.DB
	Rdb          *redis.Client
	Ctx          = context.Background()
)
