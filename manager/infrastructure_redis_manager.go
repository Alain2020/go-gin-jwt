package manager

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

type InfrastructureRedisManagerEntity interface {
	GetRedisClient() *redis.Client
}

type InfrastructureRedisManager struct {
	redisClient *redis.Client
}

func NewInfrastructureRedisManager(appConfigManager AppConfigManagerManagerEntity) InfrastructureRedisManagerEntity {
	redisHost := appConfigManager.GetAppConfig().RedisHost
	redisPort := appConfigManager.GetAppConfig().RedisPort
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
		DB:   0,
	})
	return &InfrastructureRedisManager{redisClient: redisClient}
}

func (i *InfrastructureRedisManager) GetRedisClient() *redis.Client {
	return i.redisClient
}
