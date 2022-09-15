package etcenter

import "context"

var cli  *Client
func Set(c *Client){
	cli= c
}

func lock(ctx context.Context,key ,val string,expireTime int) (bool, error) {
	lua := "if redis.call('setnx',KEYS[1],ARGV[1]) == 1 then redis.call('expire',KEYS[1],ARGV[2]) return 1 else return 0 end"
	var num int64
	switch cli.Type {
	case client:
		v, err := cli.Client.Eval(ctx, lua, []string{key}, val, int64(expireTime)).Result()
		if err != nil {
			return false, err
		}
		num = v.(int64)
	case clusterClient:
		v, err := cli.ClusterClient.Eval(ctx, lua, []string{key}, val, int64(expireTime)).Result()
		if err != nil {
			return false, err
		}
		num = v.(int64)
	}

	if num == 1 {
		return true, nil
	}
	return false, nil
}

func unLock(ctx context.Context,key ,val string) (bool, error) {
	lua := "if redis.call('get',KEYS[1]) == ARGV[1] then return redis.call('del',KEYS[1]) else return 0 end"
	var num int64
	switch cli.Type {
	case client:
		v, err := cli.Client.Eval(ctx, lua, []string{key}, val).Result()
		if err != nil {
			return false, err
		}
		num = v.(int64)
	case clusterClient:
		v, err := cli.ClusterClient.Eval(ctx, lua, []string{key}, val).Result()
		if err != nil {
			return false, err
		}
		num = v.(int64)
	}

	if num == 1 {
		return true, nil
	}
	return false, nil
}

func expire(ctx context.Context,key,val string,expireTime int) (bool, error) {
	lua := "if redis.call('get',KEYS[1]) == ARGV[1] then return redis.call('expire',KEYS[1],ARGV[2]) else return 0 end"
	var num int64
	switch cli.Type {
	case client:
		v, err := cli.Client.Eval(ctx, lua, []string{key}, val, expireTime).Result()
		if err != nil {
			return false, err
		}
		num = v.(int64)
	case clusterClient:
		v, err := cli.ClusterClient.Eval(ctx, lua, []string{key}, val, expireTime).Result()
		if err != nil {
			return false, err
		}
		num = v.(int64)
	}
	if num == 1 {
		return true, nil
	}
	return false, nil
}

