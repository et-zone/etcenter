package etcenter

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

type Worker struct {
	mu sync.Mutex
	cli *Client
	heartTime int  //second
	expiration int //second
	Key string
	Val string
	*time.Ticker
	context.Context
	runing bool
}

//  heartTime ,expiration ( second) & heartTime < expiration
func NewWorker(ctx context.Context,cli *Client,key string,heartTime,expiration int )(*Worker,error){
	if heartTime<=0||expiration<=0 ||heartTime>=expiration{
		return nil,errors.New("args: 0 < heartTime < expiration ")
	}
	m:= &Worker{
		sync.Mutex{},
		cli,
		heartTime,
		expiration,
		key,
		fmt.Sprintf("%v%v%v",time.Now().Unix(),os.Getgid(),os.Getpid()),
		time.NewTicker(time.Duration(heartTime)*time.Second),
		ctx,
		false,
	}
	return m,nil
}

func (m *Worker)register()(bool,error){
	lua:="if redis.call('setnx',KEYS[1],ARGV[1]) == 1 then redis.call('expire',KEYS[1],ARGV[2]) return 1 else return 0 end"
	var num int64
	switch m.cli.Type {
	case client:
		v,err:=m.cli.Client.Eval(m.Context,lua,[]string{m.Key},m.Val,int64(m.expiration)).Result()
		if err!=nil{
			return false,err
		}
		num=v.(int64)
	case clusterClient:
		v,err:=m.cli.ClusterClient.Eval(m.Context,lua,[]string{m.Key},m.Val,int64(m.expiration)).Result()
		if err!=nil{
			return false,err
		}
		num=v.(int64)
	}

	if num==1{
		m.mu.Lock()
		m.runing=true
		m.mu.Unlock()
		return true,nil
	}
	return false,nil
}

// unlock yourself lock, val is unique
func (m *Worker)unRegister(ctx context.Context, key, val string)(bool,error){
	lua:="if redis.call('get',KEYS[1]) == ARGV[1] then return redis.call('del',KEYS[1]) else return 0 end"
	var num int64
	switch m.cli.Type {
	case client:
		v,err:=m.cli.Client.Eval(ctx,lua,[]string{key},val).Result()
		if err!=nil{
			return false,err
		}
		num=v.(int64)
	case clusterClient:
		v,err:=m.cli.ClusterClient.Eval(ctx,lua,[]string{key},val).Result()
		if err!=nil{
			return false,err
		}
		num=v.(int64)
	}

	if num==1{
		m.mu.Lock()
		m.runing=false
		m.mu.Unlock()
		return true,nil
	}
	return false,nil
}

func (m *Worker)watch()(bool,error){
	lua:="if redis.call('get',KEYS[1]) == ARGV[1] then return redis.call('expire',KEYS[1],ARGV[2]) else return 0 end"
	var num int64
	switch m.cli.Type {
	case client:
		v,err:=m.cli.Client.Eval(m.Context,lua,[]string{m.Key},m.Val,m.expiration).Result()
		if err!=nil{
			return false,err
		}
		num=v.(int64)
	case clusterClient:
		v,err:=m.cli.ClusterClient.Eval(m.Context,lua,[]string{m.Key},m.Val,m.expiration).Result()
		if err!=nil{
			return false,err
		}
		num=v.(int64)
	}

	if num==1{
		return true,nil
	}
	return false,nil
}

func (m *Worker)SyncWatch(){
	defer m.Ticker.Stop()
	for{
		<-m.Ticker.C
		if m.runing==true{
			ok,err:=m.watch()
			fmt.Println(ok,err)
		}else {
			ok,err:=m.register()
			if err!=nil{
				fmt.Println(err.Error())
			}
			if ok{
				fmt.Println("register succ!")
			}else {
				fmt.Println("registering...")
			}
		}
	}
}

func  (m *Worker)Runing()bool{
	return m.runing
}