package etcenter

import "sync"

//key 之所以需要续期，是这个key影响到程序的执行，基本上是针对脚本执行才会关心续期的问题。
//如：我们确定每5s执行一次，每次执行时长一定不会超过3s，那么我们也就没有必要考虑续期的问题、
// 也就是说，必须一次性执行完的任务，我们需要续期解决问题，如果可以拆分每5s执行，每次去2个操作，那么就不存在续期的问题
//针对较长时间的脚本任务，防止多服务并发执行，此时我们针对这些小部分的key进行续期

type SyncExpire struct {
	m sync.Map
}

func NewSyncExpire()*SyncExpire{
	return &SyncExpire{
		sync.Map{},
	}
}

type KV struct {
	Key string
	Val string
	ExCount int
	MaxCount int
}

func (s *SyncExpire) Store(kv *KV){
	s.m.Store(kv.Key,kv)
}

func (s *SyncExpire) Load(key string)(*KV,bool){
	val,ok:=s.m.Load(key)
	if !ok{
		return &KV{},false
	}
	return val.(*KV),ok
}

// 性能 空跑 100w 个耗时130ms
func (s *SyncExpire) Range(f func(kv *KV)bool){
	s.m.Range(func(key interface{},val interface{})bool{
		v:=val.(*KV)
		if v.MaxCount>0{
			if v.ExCount>=v.MaxCount{
				s.m.Delete(key)
			}else {
				if f(v){
					v.ExCount++
				}else {
					s.m.Delete(key)
				}
			}
		}
		return true
	})
}
