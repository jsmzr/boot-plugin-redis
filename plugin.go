package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	plugin "github.com/jsmzr/boot-plugin"
	"github.com/jsmzr/boot-plugin-redis/db"
	"github.com/spf13/viper"
)

type RedisPlugin struct{}

const configPrefix = "boot.redis."
const clusterConfigPrefix = "boot.redis.cluster."
const sentinelConfigPrefix = "boot.redis.cluster."

var defaultConfig map[string]interface{} = map[string]interface{}{
	"enabled": true,
	"order":   10,
	// single, cluster, sentinel
	"type":     "single",
	"db":       0,
	"poolSize": 10,
}

func (r *RedisPlugin) Enabled() bool {
	return viper.GetBool(configPrefix + "enabled")
}

func (r *RedisPlugin) Order() int {
	return viper.GetInt(configPrefix + "order")
}

func (r *RedisPlugin) Load() error {
	clusterType := viper.GetString(configPrefix + "type")
	switch clusterType {
	case "single":
		return initSingle()
	case "cluster":
		return initCluster()
	case "sentinel":
		return initSentinel()
	}
	return fmt.Errorf("not found redis type:[%s], please use [single,cluster,sentinel]", clusterType)
}

func initSingle() error {
	config := redis.Options{
		Addr:     viper.GetString(configPrefix + "address"),
		Password: viper.GetString(configPrefix + "password"),
		DB:       viper.GetInt(configPrefix + "db"),
		Username: viper.GetString(configPrefix + "username"),
		PoolSize: viper.GetInt(configPrefix + "poolSize"),
	}
	rdb := redis.NewClient(&config)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return err
	} else {
		db.Client = rdb
		return nil
	}
}

func initCluster() error {
	config := redis.ClusterOptions{
		Addrs:    viper.GetStringSlice(clusterConfigPrefix + "cluster"),
		Username: viper.GetString(clusterConfigPrefix + "username"),
		Password: viper.GetString(clusterConfigPrefix + "password"),
	}

	rdb := redis.NewClusterClient(&config)
	c := context.Background()
	if err := rdb.ForEachShard(c, func(ctx context.Context, client *redis.Client) error {
		return client.Ping(ctx).Err()
	}); err != nil {
		return err
	} else {
		db.Cluster = rdb
		return nil
	}
}

func initSentinel() error {
	config := redis.FailoverOptions{
		// sentinel use boot.redis.sentinel.password, username, address, masterName
		MasterName:       viper.GetString(sentinelConfigPrefix + "masterName"),
		SentinelAddrs:    viper.GetStringSlice(sentinelConfigPrefix + "address"),
		SentinelUsername: viper.GetString(sentinelConfigPrefix + "username"),
		SentinelPassword: viper.GetString(sentinelConfigPrefix + "password"),
		// node use boot.redis.password,username,db
		Username: viper.GetString(configPrefix + "username"),
		Password: viper.GetString(configPrefix + "password"),
		DB:       viper.GetInt(configPrefix + "db"),
		// pool
		PoolSize: viper.GetInt(configPrefix + "poolSize"),
	}
	rdb := redis.NewFailoverClient(&config)

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return err
	} else {
		db.Client = rdb
		return nil
	}
}

func init() {
	for key := range defaultConfig {
		viper.SetDefault(configPrefix+key, defaultConfig[key])
	}
	plugin.Register("redis", &RedisPlugin{})
}
