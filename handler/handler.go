package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type Jsonresp struct {
	Code    string      `json:"code"`
	Ret     int         `json:"ret"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var validate = validator.New()

func Validate(r *http.Request, dst interface{}) error {
	data := map[string]interface{}{}
	for k, v := range r.URL.Query() {
		data[k] = v
	}
	err := r.ParseForm()
	if err != nil {
		return err
	}
	for k, v := range r.Form {
		data[k] = v
	}
	if r.Method == "POST" && r.Header.Get("content-type") == "application/json" {
		err = json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			return err
		}
	}

	err = mapstructure.WeakDecode(data, dst)
	if err != nil {
		return err
	}
	return validate.Struct(dst)
}

type JsonServer interface {
	Serve(r *http.Request) (interface{}, error)
}

func JsonHandler(s JsonServer) http.Handler {
	return JsonHandlerFunc(s.Serve)
}

type JsonHandlerFunc func(r *http.Request) (interface{}, error)

func (f JsonHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("panic", r).Error("ServeHTTP panic")
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%#v", r)
			}
			WriteJsonresp(JsonrespInterServerErr(err), w)
		}
	}()
	resp, err := f(r)

	var jsonresp Jsonresp
	if err != nil {
		jsonresp = JsonrespInterServerErr(err)
	} else {
		jsonresp = JsonrespSuccess(resp)
	}
	WriteJsonresp(jsonresp, w)
}

func JsonrespSuccess(data interface{}) Jsonresp {
	return Jsonresp{"0000", 200, "请求成功", data}
}

func JsonrespInterServerErr(err error) Jsonresp {
	return Jsonresp{"9999", 200, err.Error(), nil}
}

func JsonrespInterParamsErr(err error) Jsonresp {
	msg := fmt.Sprintf("Invalid parameters: %s", err.Error())
	return Jsonresp{"4001", 200, msg, nil}
}

func JsonrespInterOperationNotAllowed(err error) Jsonresp {
	msg := fmt.Sprintf("Operation not allowed: %s", err.Error())
	return Jsonresp{"4002", 200, msg, nil}
}

func JsonrespInterForbidden(err error) Jsonresp {
	msg := fmt.Sprintf("Forbidden: %s", err.Error())
	return Jsonresp{"4003", 200, msg, nil}
}

func WriteJsonresp(resp Jsonresp, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,X-Requested-With,authToken")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Expose-Headers", "*")
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		resp = JsonrespInterServerErr(fmt.Errorf("%#v json Encode error %s", resp, err.Error()))
		json.NewEncoder(w).Encode(resp)
	}
}
