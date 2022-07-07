package repository

import (
	"context"
	"fmt"
	"net"
	"time"

	rs "github.com/go-redis/redis/v8"
	"github.com/gomodule/redigo/redis"
)

var (
	ctx = context.Background()
)

func initRedisPool(ctx context.Context) *redis.Pool {
	redisHost, _ := helpers.GetEnvVar("REDIS_HOST")
	redisPort, _ := helpers.GetEnvVar("REDIS_PORT")
	redisPassword, _ := helpers.GetEnvVar("REDIS_PASSWORD")
	redisUser, _ := helpers.GetEnvVar("REDIS_USER")

	pool := &redis.Pool{
		MaxIdle:   0,
		MaxActive: 12000,
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			conn, err := redis.DialContext(ctx, "tcp", redisHost+":"+redisPort)
			if err != nil {
				return nil, fmt.Errorf("error occurred when init redis connection. detail : %s", err.Error())
			}

			if redisUser != "" && redisPassword != "" {
				if _, err := conn.Do("AUTH", redisUser, redisPassword); err != nil {
					conn.Close()
					return nil, fmt.Errorf("error occurred when init redis connection. detail : %s", err.Error())
				}
			}

			if redisPassword != "" {
				if _, err := conn.Do("AUTH", redisPassword); err != nil {
					conn.Close()
					return nil, fmt.Errorf("error occurred when init redis connection. detail : %s", err.Error())
				}
			}

			return conn, err
		},
	}

	return pool
}

func InitRedisSentinel() (string, *rs.Client, *rs.Client, error) {
	redisMasterName, _ := helpers.GetEnvVar("REDIS_MASTERNAME")
	redisAuthHA, _ := helpers.GetEnvVar("REDIS_AUTH_HA")
	redisSentinelIP1, _ := helpers.GetEnvVar("REDIS_SENTINEL_IP1")
	redisSentinelIP2, _ := helpers.GetEnvVar("REDIS_SENTINEL_IP2")
	redisSentinelIP3, _ := helpers.GetEnvVar("REDIS_SENTINEL_IP3")
	redisSentinelPort, _ := helpers.GetEnvVar("REDIS_SENTINEL_PORT")

	client := rs.NewFailoverClient(&rs.FailoverOptions{
		MasterName: redisMasterName,
		Password:   redisAuthHA,
		SentinelAddrs: []string{
			redisSentinelIP1 + ":" + redisSentinelPort,
			redisSentinelIP2 + ":" + redisSentinelPort,
			redisSentinelIP3 + ":" + redisSentinelPort},
		MaxRetries: -1,
	})
	
	sentinel := rs.NewSentinelClient(&rs.Options{
		Addr: ":" + redisSentinelPort,
		MaxRetries: -1,
	})
	defer sentinel.Close()

	masterAddr, err := sentinel.GetMasterAddrByName(ctx, redisMasterName).Result()
	if err != nil {
		return "", nil, nil, err
	}

	masterAddrString := net.JoinHostPort(masterAddr[0], masterAddr[1])
	master := rs.NewClient(&rs.Options{
		Addr:       masterAddrString,
		MaxRetries: -1,
	})

	return masterAddrString, client, master, nil
}

func SetRedis(key string, value interface{}, expiredtimeinsecond int) error {
	redisHA, _ := helpers.GetEnvVar("REDIS_HA")

	datalog.Info.Printf("...Adding key %s and its value to redis.\n", key)

	if redisHA == "true" {
		masterAddr, client, master, err := InitRedisSentinel()
		if err != nil {
			return fmt.Errorf("redis addr : %s err : %s", masterAddr, err.Error())
		}
		defer client.Close()
		defer master.Close()

		if expiredtimeinsecond == 0 {
			expiredtimeinsecond = helpers.DefaultExpiredTimeInSecond
		}

		err = master.Set(ctx, key, value, time.Duration(expiredtimeinsecond)).Err()
		if err != nil {
			return fmt.Errorf("error occurred when setting key : %s and value : %v to redis. redis addr : %s detail :  %s", key, value, masterAddr, err.Error())
		}

		datalog.Info.Printf("Key : %s and its value added to redis.\n", key)
		return nil
	} else {
		singleClient, err := initRedisPool(ctx).GetContext(ctx)
		if err != nil {
			return err
		}
		defer singleClient.Close()
	
		if expiredtimeinsecond == 0 {
			expiredtimeinsecond = helpers.DefaultExpiredTimeInSecond
		}
	
		setStatus, err := redis.String(singleClient.Do("SET", key, value, "EX", expiredtimeinsecond))
		if err != nil {
			return fmt.Errorf("error occurred when setting key : %s and value : %v to redis. detail : %s", key, value, err.Error())
		}
	
		datalog.Info.Printf("Set status : %s\n", setStatus)
		datalog.Info.Printf("Key : %s and its value added to redis.\n", key)
		return nil
	}
}

func GetRedis(key string) (string, error) {
	redisHA, _ := helpers.GetEnvVar("REDIS_HA")

	datalog.Info.Printf("...Getting value by key : %s from redis\n", key)

	if redisHA == "true" {
		masterAddr, client, master, err := InitRedisSentinel()
		if err != nil {
			return "", fmt.Errorf("redis addr : %s err : %s", masterAddr, err.Error())
		}
		defer client.Close()
		defer master.Close()

		val, err := client.Get(ctx, key).Result()
		if err != nil {
			if err.Error() == rs.Nil.Error() {
				return "", nil
			}

			return "", fmt.Errorf("error occurred when getting value by key : %s from redis. detail : %s", key, err.Error())
		}

		return val, nil
	} else {
		singleClient, err := initRedisPool(ctx).GetContext(ctx)
		if err != nil {
			return "", err
		}
		defer singleClient.Close()
	
		val, err := redis.String(singleClient.Do("GET", key))
		if err != nil {
			if err.Error() == redis.ErrNil.Error() {
				return "", nil
			}
	
			return "", fmt.Errorf("error occurred when getting value by key : %s from redis. detail : %s", key, err.Error())
		}
	
		return val, nil
	}
}

func FindKey(key string) (res string, err error) {
	redisHA, _ := helpers.GetEnvVar("REDIS_HA")

	datalog.Info.Printf("...Finding key(s) by which contains '%s' from redis\n", key)

	var resp []string
	if redisHA == "true" {
		var (
			masterAddr string
			client, master *rs.Client
		)
		masterAddr, client, master, err = InitRedisSentinel()
		if err != nil {
			return "", fmt.Errorf("redis addr : %s err : %s", masterAddr, err.Error())
		}
		defer client.Close()
		defer master.Close()

		resp, err = client.Keys(ctx, key).Result()
		if err != nil {
			return "", fmt.Errorf("error occurred when Getting key which contains '%s' from redis. detail : %s", key, err.Error())
		}

		if len(resp) > 0 {
			res = "keys found"
			datalog.Info.Println("keys found: ", resp)
			return
		}

		return res, err
	} else {
		var singleClient redis.Conn
		singleClient, err = initRedisPool(ctx).GetContext(ctx)
		if err != nil {
			return
		}
		defer singleClient.Close()
	
		resp, err = redis.Strings(singleClient.Do("KEYS", key))
		if err != nil {
			return
		}
	
		if len(resp) > 0 {
			res = "keys found"
			datalog.Info.Println("keys found: ", resp)
			return
		}

		return res, err
	}

	
}

func GetRedisKey(key string) ([]string, error) {
	redisHA, _ := helpers.GetEnvVar("REDIS_HA")

	datalog.Info.Printf("...Getting key which contains '%s' from redis\n", key)

	var keyResult []string
	if redisHA == "true" {
		masterAddr, client, master, err := InitRedisSentinel()
		if err != nil {
			return nil, fmt.Errorf("redis addr : %s err : %s", masterAddr, err.Error())
		}
		defer client.Close()
		defer master.Close()

		keyResult, err = client.Keys(ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("error occurred when Getting key which contains '%s' from redis. detail : %s", key, err.Error())
		}
	} else {
		singleClient, err := initRedisPool(ctx).GetContext(ctx)
		if err != nil {
			return nil, err
		}
		defer singleClient.Close()
	
		keyResult, err = redis.Strings(singleClient.Do("KEYS", key))
		if err != nil {
			return nil, fmt.Errorf("error occurred when Getting key by which contains '%s' from redis. detail : %s", key, err.Error())
		}
	}

	if len(keyResult) == 0 {
		datalog.Info.Printf("There is no key which contains %s.\n", key)
	}

	return keyResult, nil
}