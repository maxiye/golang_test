package test

import (
	"github.com/go-redis/redis"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{"172.17.12.152:6381", "172.17.12.152:6382", "172.17.12.152:6383", "172.17.12.152:6384", "172.17.12.152:6385", "172.17.12.152:6386"},
		Password: "im2NCnCwweA=",
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
