package db

import "github.com/go-redis/redis/v8"

var Client *redis.Client

var Cluster *redis.ClusterClient
