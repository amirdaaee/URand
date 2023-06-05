package config

import (
	"context"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func (f *NSType) UnmarshalText(text []byte) error {
	ns := strings.Split(string(text), ":")
	nsName := ns[0]
	nsLen := int64(3)
	if len(ns) > 1 {
		_nsLen, err := strconv.ParseInt(ns[1], 10, 64)
		if err != nil {
			return err
		}
		nsLen = _nsLen
	}
	*f = NSType{Name: nsName, Length: uint(nsLen)}
	return nil
}
func (f *RedisURLType) UnmarshalText(text []byte) error {
	opts, err := redis.ParseURL(string(text))
	if err != nil {
		return err
	}
	ping := redis.NewClient(opts).Ping(context.Background())
	if ping.Err() != nil {
		return ping.Err()
	}
	*f = RedisURLType(text)
	return nil
}

func (f *LogLevelType) UnmarshalText(text []byte) error {
	ll, err := logrus.ParseLevel(string(text))
	if err != nil {
		return err
	}
	logrus.SetLevel(ll)
	*f = LogLevelType(text)
	return nil
}
