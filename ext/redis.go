package ext

import (
	"errors"
	"fmt"
	"github.com/FZambia/sentinel"
	"github.com/garyburd/redigo/redis"
	sentinelRedis "github.com/gomodule/redigo/redis"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RedisCache struct {
	Abstract
	sync.RWMutex
	pool         map[string]*redis.Pool
	sentinelPool map[string]*sentinelRedis.Pool
}

var Redis *RedisCache

func init() {
	Redis = &RedisCache{
		RWMutex:      sync.RWMutex{},
		pool:         map[string]*redis.Pool{},
		sentinelPool: map[string]*sentinelRedis.Pool{},
	}
}

// 获取redis Pool
func (rc *RedisCache) Pool(identify string) *redis.Pool {
	// 读取缓存
	if redisPool, ok := rc.pool[identify]; ok {
		return redisPool
	}
	// 加锁
	rc.Lock()
	defer rc.Unlock()
	if redisPool, ok := rc.pool[identify]; ok {
		return redisPool
	}
	// 读取配置
	config := rc.config.Get(identify).ToStringMap()
	if len(config) == 0 {
		panic(fmt.Sprintf("get config by identify '%s' failed.", identify))
	}
	// 参数校验
	stringMap := map[string]string{
		"Host":     "",
		"Password": "",
	}
	for key := range stringMap {
		if value, ok := config[key]; ok {
			stringMap[key] = value
		}
	}
	for key, value := range stringMap {
		if value == "" && key != "Password" {
			panic(fmt.Sprintf("the config's '%s' can't be empty.", key))
		}
	}

	// 整型参数校验
	intMap := map[string]int{"Db": 0, "MaxIdle": 0, "MaxActive": 0, "IdleTimeout": 0}
	for key := range intMap {
		if value, ok := config[key]; ok {
			if v, err := strconv.Atoi(value); err == nil {
				intMap[key] = v
			}
		}
	}

	// 创建连接池
	rc.pool[identify] = &redis.Pool{
		MaxIdle:     intMap["MaxIdle"],
		MaxActive:   intMap["MaxActive"],
		IdleTimeout: time.Duration(intMap["IdleTimeout"]),
		Dial: func() (redis.Conn, error) {
			dialOption := []redis.DialOption{
				redis.DialReadTimeout(time.Duration(intMap["IdleTimeout"]) * time.Millisecond),
				redis.DialWriteTimeout(time.Duration(intMap["IdleTimeout"]) * time.Millisecond),
				redis.DialConnectTimeout(time.Duration(intMap["IdleTimeout"]) * time.Millisecond),
				redis.DialDatabase(intMap["Db"]),
			}
			if stringMap["Password"] != "" {
				dialOption = append(dialOption, redis.DialPassword(stringMap["Password"]))
			}
			return redis.Dial("tcp", stringMap["Host"], dialOption...)
		},
	}
	return rc.pool[identify]
}

// 获取redis SentinelPool
func (rc *RedisCache) SentinelPool(identify string) *sentinelRedis.Pool {
	// 读取缓存
	if redisPool, ok := rc.sentinelPool[identify]; ok {
		return redisPool
	}
	// 加锁
	rc.Lock()
	defer rc.Unlock()
	if redisPool, ok := rc.sentinelPool[identify]; ok {
		return redisPool
	}
	// 读取配置
	config := rc.config.Get(identify).ToStringMap()
	if len(config) == 0 {
		panic(fmt.Sprintf("get config by identify '%s' failed.", identify))
	}

	// Sentinel类
	stringMap := map[string]string{
		"HostArray":  "",
		"MasterName": "",
		"Password":   "",
	}
	for key := range stringMap {
		if value, ok := config[key]; ok {
			stringMap[key] = value
		}
	}
	for key, value := range stringMap {
		if value == "" && key != "Password" {
			panic(fmt.Sprintf("the config's '%s' can't be empty.", key))
		}
	}
	st := &sentinel.Sentinel{
		Addrs:      strings.Split(stringMap["HostArray"], ","),
		MasterName: stringMap["MasterName"],
		Dial: func(addr string) (sentinelRedis.Conn, error) {
			timeout := 500 * time.Millisecond
			c, err := sentinelRedis.DialTimeout("tcp", addr, timeout, timeout, timeout)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

	// Sentinel 连接池
	intMap := map[string]int{"Db": 0, "MaxIdle": 0, "MaxActive": 0, "IdleTimeout": 0}
	for key := range intMap {
		if value, ok := config[key]; ok {
			if v, err := strconv.Atoi(value); err == nil {
				intMap[key] = v
			}
		}
	}
	rc.sentinelPool[identify] = &sentinelRedis.Pool{
		MaxIdle:     intMap["MaxIdle"],
		MaxActive:   intMap["MaxActive"],
		Wait:        true,
		IdleTimeout: time.Duration(intMap["IdleTimeout"]) * time.Second,
		Dial: func() (sentinelRedis.Conn, error) {
			// 主地址
			masterAddr, err := st.MasterAddr()
			if err != nil {
				return nil, err
			}
			// dial设置
			dialOption := []sentinelRedis.DialOption{
				sentinelRedis.DialReadTimeout(time.Duration(intMap["IdleTimeout"]) * time.Millisecond),
				sentinelRedis.DialWriteTimeout(time.Duration(intMap["IdleTimeout"]) * time.Millisecond),
				sentinelRedis.DialConnectTimeout(time.Duration(intMap["IdleTimeout"]) * time.Millisecond),
				sentinelRedis.DialDatabase(intMap["Db"]),
			}
			if stringMap["Password"] != "" {
				dialOption = append(dialOption, sentinelRedis.DialPassword(stringMap["Password"]))
			}
			c, err := sentinelRedis.Dial("tcp", masterAddr, dialOption...)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c sentinelRedis.Conn, t time.Time) error {
			if !sentinel.TestRole(c, "master") {
				return errors.New("check role failed")
			} else {
				return nil
			}
		},
	}
	return rc.sentinelPool[identify]
}
