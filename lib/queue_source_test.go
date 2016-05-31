// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package lib

import (
	"testing"
	"net/http"
	"fmt"
	"log"
	"time"
)

func TestQueueSourceMysqlTesting(t *testing.T) {
	qscfg:=&QueueSourceConfig{Method:QUEUE_SOURCE_METHOD_MYSQL, MysqlDsn:"root:@tcp(localhost:3306)/test?autocommit=true", Value:"select PushID from `device_tokens` group by PushID"}

	qs, err:=NewQueueSourceByConfig(qscfg)
	if err!=nil {
		t.Errorf("Found Error: %v", err)
	}else{
		list, err:=qs.GetData()
		if err!=nil {
			t.Errorf("Found Error qs.GetData: %v", err)
		}else{
			if len(list)<=0 {
				t.Errorf("Found Error qs.GetData: len(list)<=0")
			}else{
				t.Logf("Fetch data: %v", list)
			}

		}
	}
}

func TestQueueSourceApiTesting(t *testing.T) {
	go func() {
		//if /bar/指定, /bar自动301到/bar/,也可以两个都指定
		http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "3523544012e5491b3fe8cf6627eddd123d6aa4191fbebf371191a3ce7d4c02ac\nefdd029e3e62ab46bf089bfe7084d3261471b6f9e0e4225f9851b4e5b8e7f57e")
		})

		log.Println("Server started...")
		log.Fatal(http.ListenAndServe(":9998", nil))
	}()

	//wait server up
	time.Sleep(time.Second)
	//api queue test
	qscfg:=&QueueSourceConfig{Method:QUEUE_SOURCE_METHOD_API, ApiPrefix:"http://127.0.0.1:9998/test?queue=", Value:"test"}

	qs, err:=NewQueueSourceByConfig(qscfg)
	if err!=nil {
		t.Errorf("Found Error: %v", err)
	}else{
		list, err:=qs.GetData()
		if err!=nil {
			t.Errorf("Found Error qs.GetData: %v", err)
		}else{
			if len(list)<=0 {
				t.Errorf("Found Error qs.GetData: len(list)<=0")
			}else{
				t.Logf("Fetch data: %v", list)
			}
		}
	}

}

func TestQueueSourceFileTesting(t *testing.T) {
	qscfg:=&QueueSourceConfig{Method:QUEUE_SOURCE_METHOD_FILE, FilePath:"/Users/bruce/project/godev/src/gopush/runtime/data/%s.txt", Value:"queue"}

	qs, err:=NewQueueSourceByConfig(qscfg)
	if err!=nil {
		t.Errorf("Found Error: %v", err)
	}else{
		list, err:=qs.GetData()
		if err!=nil {
			t.Errorf("Found Error qs.GetData: %v", err)
		}else{
			if len(list)<=0 {
				t.Errorf("Found Error qs.GetData: len(list)<=0")
			}else{
				t.Logf("Fetch data: %v", list)
			}
		}
	}
}