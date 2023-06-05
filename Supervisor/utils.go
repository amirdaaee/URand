package supervisor

import (
	config "URand/Config"

	"github.com/redis/go-redis/v9"
)

var rds *redis.Client
var cfg *config.SettingsType

func getRedis() *redis.Client {
	if rds == nil {
		opts, _ := redis.ParseURL(string(getCfg().RedisURL))
		rds = redis.NewClient(opts)
	}
	return rds
}
func getCfg() *config.SettingsType {
	if cfg == nil {
		cfg = config.Config()

	}
	return cfg
}
