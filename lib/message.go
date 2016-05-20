// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package lib

import (
	"encoding/json"
)

type MessageInterface interface {
	// fetch message title
	GetTitle() string

	// fetch message body
	GetBody() string

	// fetch sound info
	GetSound() string

	// fetch custom info
	GetCustom() map[string]string

	GetUuid() string

	MarshalJSON() (string, error)
}

type Message struct {
	Title  string `json:"title"`
	Body   string `json:"body"`

	//custom field
	Custom map[string]string `json:"custom"`

	Sound  string `json:"sound"`

	Uuid   string `json:"uuid"`
}

func (m *Message)MarshalJSON() (string, error) {
	jb, err := json.Marshal(m)
	if err != nil {
		return "", nil
	}else {
		return string(jb), nil
	}
}

// fetch message title
func (m *Message)GetTitle() string {
	return m.Title
}

// fetch message body
func (m *Message)GetBody() string {
	return m.Body
}

// fetch sound info
func (m *Message)GetSound() string {
	return m.Sound
}

// fetch custom info
func (m *Message)GetCustom() map[string]string {
	return m.Custom
}

// fetch sound info
func (m *Message)GetUuid() string {
	return m.Uuid
}