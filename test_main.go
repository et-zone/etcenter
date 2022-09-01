package etcenter

import (
	"context"
	redis "github.com/go-redis/redis/v8"
	"testing"
)

func testMain(t testing.T){
	cli:= NewClient(redis.NewClient(&redis.Options{
		Addr:     "49.232.190.11:63790",
		Password: "", // no password set
		DB:       1,  // use default DB
	}))
	w,_:= NewWorker(context.TODO(),cli,"vvvv5vvvv",3,10)
	w.SyncWatch(func() {}, func() {})
}
