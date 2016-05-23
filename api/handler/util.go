// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package handler

import (
	"net/http"
	"errors"
	"strconv"
)

const (
	API_CODE_OK = iota

//post method
	API_CODE_POST_NEEDED
//param required
	API_CODE_PARAM_REQUIRED
	API_CODE_PARAM_ERROR
	API_CODE_QUEUE_BUILD
	API_CODE_TASK_ERROR

	DEVICEID_SEP = ","
)


func formatNormalResponceHeader(w http.ResponseWriter) {
	w.Header().Add("server", "gopush")
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
}

func GetParamString(r *http.Request, name string) (string, error) {
	if param, ok := r.Form[name]; ok {
		// fetch first one
		for _, value := range param {
			return value, nil
		}
	}else {
		return "", errors.New("Param " + name + " not found")
	}

	return "", nil
}

func GetParamInt(r *http.Request, name string) (int, error) {
	param, err := GetParamString(r, name)
	if err != nil {
		return 0, err
	} else {
		int64Param, err := strconv.ParseInt(param, 10, 64)
		//may truncate
		intParam := int(int64Param)
		if err != nil {
			return 0, err
		}else {
			return intParam, nil
		}
	}
}

func GetParamArrayString(r *http.Request, name string) ([]string, error) {
	if param, ok := r.Form[name]; ok {
		var result []string
		for _, value := range param {
			result = append(result, value)
		}

		return result, nil
	}else {
		return nil, errors.New("Param " + name + " not found")
	}
}