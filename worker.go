package etcenter

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const (
	errConn   = "network is unreachable"
	heartTime = 10
	expireTime = 15
)

type Worker struct {
	mu         sync.Mutex
	Key        string
	Val        string
	*time.Ticker
	context.Context
	ready bool
}

// 0 < heartTimeSec < expireSec
func NewWorker(ctx context.Context, key string) (*Worker, error) {
		m := &Worker{
		sync.Mutex{},
		key,
		fmt.Sprintf("%v%v%v", time.Now().Unix(), os.Getgid(), os.Getpid()),
		time.NewTicker(time.Duration(heartTime) * time.Second),
		ctx,
		false,
	}
	return m, nil
}

func (m *Worker) register() (bool, error) {
	return lock(context.TODO(),m.Key,m.Val,expireTime)
}

// unlock yourself lock, val is unique
func (m *Worker) unRegister() (bool, error) {
	return unLock(context.TODO(),m.Key,m.Val)
}

func (m *Worker) expire() (bool, error) {
	return expire(context.TODO(),m.Key,m.Val,expireTime)
}

// go run
func (m *Worker) SyncWatch(taskRun, taskStop func()) {
	defer m.Ticker.Stop()
	for {
		<-m.Ticker.C
		if m.Ready() == true {
			//If add expire fail ,maybe network timeout , maybe you need close worker and wait to try register again .
			//If not ok ,maybe key is del and others server are register succ. so we must run stop task .
			//The probability of error and not ok is very small
			ok, err := m.expire()
			if err != nil || !ok {
				m.mu.Lock()
				if taskStop != nil {
					taskStop()
				}
				m.unRegister()
				m.ready = false
				m.mu.Unlock()
				continue
			}

		} else {
			// if fail ,maybe others server are succ register,you need try registering one by one
			ok, err := m.register()
			if err != nil {
				log.Println(err.Error())
				continue
			}
			if ok {
				m.mu.Lock()
				taskRun()
				m.ready = true
				m.mu.Unlock()
				log.Println(fmt.Sprintf("register succ! key = %v ,val = %v", m.Key, m.Val))
			} else {
				log.Println("registering...")
			}
		}
	}

}

func (m *Worker) Ready() bool {
	return m.ready
}
