package etcenter

//
//func testExpire(){
//	s:=etcenter.NewSyncExpire()
//	//for i:=0;i<1000000;i++{
//	//	s.Store(&etcenter.KV{Key: "aaa"+strconv.Itoa(i+1),Val: "aaa"+strconv.Itoa(i+1),MaxCount: 100})
//	//}
//	s.Store(&etcenter.KV{Key: "aaa",Val: "aaa",MaxCount: 4})
//	s.Store(&etcenter.KV{Key: "bbb",Val: "bbb",MaxCount: 7})
//	s.Store(&etcenter.KV{Key: "ccc",Val: "ccc",MaxCount: 10})
//	s.Store(&etcenter.KV{Key: "ddd",Val: "ddd",MaxCount: -1})
//	for{
//		t:=time.Now()
//		s.Range(func(kv *etcenter.KV) bool {
//			fmt.Println(fmt.Sprintf("key=%v,count=%v",kv.Key,kv.ExCount))
//			//fmt.Sprintf("key=%v,count=%v",kv.Key,kv.ExCount)
//			return true
//		})
//		fmt.Println(time.Since(t))
//		time.Sleep(time.Second)
//	}
//
//}
//
//func testWorker(){
//	cli:= etcenter.NewClient(etcenter.Options{
//		Addr:     "49.232.190.114:63790",
//		Password: "", // no password set
//		DB:       1,  // use default DB
//	})
//	etcenter.Set(cli)
//	w,_:= etcenter.NewWorker(context.TODO(),"vvvv5vvvv")
//	w.SyncWatch(func() {fmt.Println("start run...")}, func() {fmt.Println("end run...")})
//}