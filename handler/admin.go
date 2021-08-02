package handler

import (
	"apitesting/entity"
	"net/http"
)

var AdminId = JsonHandlerFunc(func(r *http.Request) Jsonresp {
	id, err := entity.User{}.AdminId()
	if err != nil {
		return JsonrespInterServerErr(err)
	}
	return JsonrespSuccess(id)
})
