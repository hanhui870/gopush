// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package handler

import (
	"net/http"
)

func formatNormalResponceHeader(w http.ResponseWriter) {
	w.Header().Add("server", "gopush")
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
}

func GetParamString(name string) {

}

func GetParamInt(name string) {

}

func GetParamArray(name string) {

}