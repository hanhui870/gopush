package lib

import (
	"testing"
	"strconv"
	"fmt"
)

func TestTaskQueueCycleOperation(t *testing.T) {
	tq := &TaskQueue{tasks:make([]*Task, 5)}

	pos, err := tq.Add(&DeviceQueue{Position:1}, nil)
	if err != nil {
		t.Fatal("Add queue task faild: " + err.Error())
	}
	if pos != 0 {
		t.Fatal("Add queue task pos error: " + strconv.Itoa(pos))
	}
	tq.Add(&DeviceQueue{Position:2}, nil)
	tq.Add(&DeviceQueue{Position:3}, nil)
	pos, err = tq.Add(&DeviceQueue{Position:4}, nil)
	if pos != 3 {
		t.Fatal("Add queue task pos error: " + strconv.Itoa(pos))
	}
	tq.Add(&DeviceQueue{Position:5}, nil)
	pos, err = tq.Add(&DeviceQueue{Position:6}, nil)
	fmt.Println(tq.tasks)
	if err == nil {
		t.Fatal("Add queue task should faild, queue full.")
	}else {
		fmt.Println(err)
	}

	task, err := tq.Read()
	if err != nil {
		t.Fatal("Read queue task faild: " + err.Error())
	}else {
		t.Log("Task position:" + strconv.Itoa(task.list.Position))
	}

	err = tq.Pop()
	if err != nil {
		t.Fatal("Pop queue task faild: " + err.Error())
	}
	fmt.Println(tq.tasks)

	pos, err = tq.Add(&DeviceQueue{Position:6}, nil)
	fmt.Println(tq.tasks)
	if err != nil {
		t.Fatal("Add queue task faild: " + err.Error())
	}
	if pos != 4 {
		t.Fatal("Add queue task pos error: " + strconv.Itoa(pos))
	}

	task, err = tq.Read()
	if err != nil {
		t.Fatal("Read queue task faild: " + err.Error())
	}else {
		t.Log("Task position:" + strconv.Itoa(task.list.Position))
	}
}