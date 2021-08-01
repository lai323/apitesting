package handler

import (
	"apitesting/entity"
	"net/http"
)

var AdminId = JsonHandlerFunc(func(r *http.Request) (interface{}, error) {
	id, err := entity.User{}.AdminId()
	if err != nil {
		return JsonrespInterServerErr(err), nil
	}
	return JsonrespSuccess(id), nil
})
