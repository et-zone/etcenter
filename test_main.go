package etcenter

import (
	"context"
	"testing"
)

func testMain(t testing.T){
	cli:= NewClient(Options{
		Addr:     "49.232.190.11:63790",
		Password: "", // no password set
		DB:       1,  // use default DB
	})
	w,_:= NewWorker(context.TODO(),cli,"vvvv5vvvv",3,10)
	w.SyncWatch(func() {}, func() {})
}
