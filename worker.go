package etcenter

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const retry = 5
const errConn = "network is unreachable"

type Worker struct {
	mu         sync.Mutex
	cli        *Client
	heartTime  int //second
	expiration int //second
	Key        string
	Val        string
	*time.Ticker
	context.Context
	ready bool
}

//  heartTime ,expiration ( second) & heartTime < expiration
func NewWorker(ctx context.Context, cli *Client, key string, heartTime, expiration int) (*Worker, error) {
	if heartTime <= 0 || expiration <= 0 || heartTime >= expiration {
		return nil, errors.New("args: 0 < heartTime < expiration ")
	}
	m := &Worker{
		sync.Mutex{},
		cli,
		heartTime,
		expiration,
		key,
		fmt.Sprintf("%v%v%v", time.Now().Unix(), os.Getgid(), os.Getpid()),
		time.NewTicker(time.Duration(heartTime) * time.Second),
		ctx,
		false,
	}
	return m, nil
}

func (m *Worker) register() (bool, error) {
	lua := "if redis.call('setnx',KEYS[1],ARGV[1]) == 1 then redis.call('expire',KEYS[1],ARGV[2]) return 1 else return 0 end"
	var num int64
	switch m.cli.Type {
	case client:
		v, err := m.cli.Client.Eval(m.Context, lua, []string{m.Key}, m.Val, int64(m.expiration)).Result()
		if err != nil {
			return false, err
		}
		num = v.(int64)
	case clusterClient:
		v, err := m.cli.ClusterClient.Eval(m.Context, lua, []string{m.Key}, m.Val, int64(m.expiration)).Result()
		if err != nil {
			return false, err
		}
		num = v.(int64)
	}

	if num == 1 {
		m.mu.Lock()
		m.ready = true
		m.mu.Unlock()
		return true, nil
	}
	return false, nil
}

// unlock yourself lock, val is unique
func (m *Worker) unRegister() (bool, error) {
	lua := "if redis.call('get',KEYS[1]) == ARGV[1] then return redis.call('del',KEYS[1]) else return 0 end"
	var num int64
	switch m.cli.Type {
	case client:
		v, err := m.cli.Client.Eval(m.Context, lua, []string{m.Key}, m.Val).Result()
		if err != nil {
			return false, err
		}
		num = v.(int64)
	case clusterClient:
		v, err := m.cli.ClusterClient.Eval(m.Context, lua, []string{m.Key}, m.Val).Result()
		if err != nil {
			return false, err
		}
		num = v.(int64)
	}

	if num == 1 {
		m.mu.Lock()
		m.ready = false
		m.mu.Unlock()
		return true, nil
	}
	return false, nil
}

func (m *Worker) expire() (bool, error) {
	lua := "if redis.call('get',KEYS[1]) == ARGV[1] then return redis.call('expire',KEYS[1],ARGV[2]) else return 0 end"
	var num int64
	switch m.cli.Type {
	case client:
		v, err := m.cli.Client.Eval(m.Context, lua, []string{m.Key}, m.Val, m.expiration).Result()
		if err != nil {
			return false, err
		}
		num = v.(int64)
	case clusterClient:
		v, err := m.cli.ClusterClient.Eval(m.Context, lua, []string{m.Key}, m.Val, m.expiration).Result()
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

// go run
func (m *Worker) SyncWatch(taskRun,taskStop func()) {
		defer m.Ticker.Stop()
		for {
			<-m.Ticker.C
			if m.Ready() == true {
				//If add expire fail ,maybe network timeout , maybe you need close worker and wait to try register again .
				//If not ok ,maybe key is del and others server are register succ. so we must run stop task .
				//The probability of error and not ok is very small
				ok, err := m.expire()
				if err!=nil||!ok{
					if taskStop!=nil{
						m.mu.Lock()
						taskStop()
						m.mu.Unlock()
						ok,err=m.unRegister()
						if err!=nil||!ok{
							time.Sleep(time.Duration(m.expiration)*time.Second)
							m.mu.Lock()
							m.ready=false
							m.mu.Unlock()
						}
					}
				}
			} else {
				// if fail ,maybe others server are succ register,you need try registering one by one
				ok, err := m.register()
				if err != nil {
					log.Println(err.Error())
					continue
				}
				if ok {
					taskRun()
					log.Println("register succ!")
				} else {
					log.Println("registering...")
				}
			}
		}

}

func (m *Worker) Ready() bool {
	return m.ready
}