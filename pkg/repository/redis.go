package repository

import (
	"context"
	"fmt"
	"learn-golang-bain/configs"
	"net"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()
)

func initSingleRedis() *redis.Client {
	redisConfig := configs.GetConfig().Redis
	client := redis.NewClient(&redis.Options{
		Addr:		redisConfig.Host + ":" + redisConfig.Port,
		Username:	redisConfig.User,
		Password:	redisConfig.Password,
	})

	return client
}

func initRedisSentinel() (failOverClient *redis.Client, masterClient *redis.Client, err error) {
	redisConfig := configs.GetConfig().Redis
	
	failOverClient = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName: redisConfig.MasterName,
		Password:   redisConfig.Password,
		SentinelAddrs: []string{
			redisConfig.SentinelIP1 + ":" + redisConfig.SentinelPort,
			redisConfig.SentinelIP2 + ":" + redisConfig.SentinelPort,
			redisConfig.SentinelIP3 + ":" + redisConfig.SentinelPort},
		MaxRetries: -1,
	})
	
	sentinel := redis.NewSentinelClient(&redis.Options{
		Addr: ":" + redisConfig.SentinelPort,
		MaxRetries: -1,
	})
	defer sentinel.Close()

	masterAddr, err := sentinel.GetMasterAddrByName(ctx, redisConfig.MasterName).Result()
	if err != nil {
		return failOverClient, masterClient, err
	}

	masterAddrString := net.JoinHostPort(masterAddr[0], masterAddr[1])
	masterClient = redis.NewClient(&redis.Options{
		Addr:       masterAddrString,
		MaxRetries: -1,
	})

	return failOverClient, masterClient, err
}

func RedisSet(key string, value interface{}, exp /*in second*/ int) error {
	redisHA := configs.GetConfig().Redis.EnableHA

	// datalog.Info.Printf("...Adding key %s and its value to redis.\n", key)

	var usedClient *redis.Client
	if redisHA == "true" {
		failOverClient, usedClient, err := initRedisSentinel()
		if err != nil {
			return err
		}
		defer failOverClient.Close()
		defer usedClient.Close()
	} else {
		usedClient := initSingleRedis()
		defer usedClient.Close()
	}

	err := usedClient.Set(ctx, key, value, time.Duration(exp)).Err()
	if err != nil {
		return fmt.Errorf("error occurred when setting key : %s and value : %v to redis. detail :  %s", key, value, err.Error())
	}

	return nil
	// datalog.Info.Printf("Key : %s and its value added to redis.\n", key)
}

func RedisGet(key string) (string, error) {
	redisHA := configs.GetConfig().Redis.EnableHA

	// datalog.Info.Printf("...Getting value by key : %s from redis\n", key)
	var usedClient *redis.Client
	if redisHA == "true" {
		failOverClient, usedClient, err := initRedisSentinel()
		if err != nil {
			return "", fmt.Errorf("err : %s", err.Error())
		}
		defer failOverClient.Close()
		defer usedClient.Close()
	} else {
		usedClient := initSingleRedis()
		defer usedClient.Close()
	}

	val, err := usedClient.Get(ctx, key).Result()
	if err != nil {
		if err.Error() == redis.Nil.Error() {
			return "", nil
		}

		return "", fmt.Errorf("error occurred when getting value by key : %s from redis. detail : %s", key, err.Error())
	}

	return val, nil
}

func GetRedisKey(key string) (keys []string, err error) {
	redisHA := configs.GetConfig().Redis.EnableHA

	// datalog.Info.Printf("...Getting key which contains '%s' from redis\n", key)

	var usedClient *redis.Client
	if redisHA == "true" {
		failOverClient, usedClient, err := initRedisSentinel()
		if err != nil {
			return keys, err
		}
		defer failOverClient.Close()
		defer usedClient.Close()
	} else {
		usedClient := initSingleRedis()
		defer usedClient.Close()
	}

	keys, err = usedClient.Keys(ctx, key).Result()
	if err != nil {
		return keys, fmt.Errorf("error occurred when Getting key by which contains '%s' from redis. detail : %s", key, err.Error())
	}

	if len(keys) == 0 {
		// datalog.Info.Printf("There is no key which contains %s.\n", key)
	}

	return keys, err
}