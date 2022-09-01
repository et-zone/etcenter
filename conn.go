package etcenter
import (
	redis "github.com/go-redis/redis/v8"
)

const (
	client = "client"
	clusterClient = "clusterClient"
)
type Client struct {
	*redis.Client //normal client or sentinelClietn
	*redis.ClusterClient // culster
	Type string
}

func NewClient(cli *redis.Client)*Client {
	return &Client{
		cli,nil, client,
	}
}

func NewClusterClient(cli *redis.ClusterClient)*Client {
	return &Client{
		nil,cli, clusterClient,
	}
}