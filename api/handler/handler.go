// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package handler

import (
	"net/http"
	"fmt"
)

func Send(w http.ResponseWriter, r *http.Request) {
	formatNormalResponceHeader(w)

	fmt.Fprintln(w, `{"success":true,"message":"hello world."}`)
}

func AddDevice(w http.ResponseWriter, r *http.Request) {

}
