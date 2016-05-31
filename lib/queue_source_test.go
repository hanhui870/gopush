// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package lib

import (
	"testing"
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