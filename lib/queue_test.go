// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package lib

import (
	"testing"
	"strings"
	"strconv"
)

func TestMuiltWriterLogger(t *testing.T) {
	q := NewQueue()
	if cap(q.Channel) != QUEUE_DEFAULT_CAPACITY {
		t.Error("len(q.Channel)!=QUEUE_DEFAULT_CAPACITY")
	}
	err := q.AppendDataSource([]string{"fdas", "fdgswfds3425243214321"})
	if err == nil {
		t.Error("Error in q.AppendDataSource")
	}else {
		t.Log(err)
	}

	err = q.AppendDataSource([]string{"038a4c750809d70be26c7b8d9aaa5da32567147ddfb465ef6cf186c82e1a3461", "038cdc3a81335cb2ebc670a115e026122f13803bc2fa59475d6b0ebd67d8f125"})
	if err != nil {
		t.Error("Error in q.AppendDataSource: " + err.Error())
	}
	t.Log("q.data: len " + strconv.Itoa(len(q.data)) + ":" + strings.Join(q.data, ", "))

	err = q.AppendFileDataSource("/Users/bruce/project/godev/src/gopush/apns/test_data/test.txt")
	if err != nil {
		t.Error("Error in q.AppendDataSource: " + err.Error())
	}
	t.Log("q.data: len " + strconv.Itoa(len(q.data)) + ":" + strings.Join(q.data, ", "))

}