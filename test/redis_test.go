package test

import (
	"github.com/go-redis/redis"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{"172.17.210.84:7001", "172.17.210.85:7002", "172.17.210.85:7001", "172.17.210.86:7002", "172.17.210.86:7001", "172.17.210.84:7002"},
		Password: "8Mbh8Ykz",
	})
	rKey := "xq:boat:race:20200619:401123670"
	boatRace, err := rdb.HGetAll(rKey).Result()
	if err != nil {
		t.Log(err)
	}
	t.Log(boatRace)
	rdb.HIncrBy(rKey, "speed_up", 10)
	rdb.Expire(rKey, 86400*2*time.Second)
	res, err := rdb.HSetNX(rKey, "speed", 200).Result()
	t.Log(res, err)
}
