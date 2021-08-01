package handler

import (
	"apitesting/entity"
	"net/http"
)

type SaveApiReq struct {
	Api     map[string]interface{}
	Testing map[string]interface{}
}

var SaveApi = JsonHandlerFunc(func(r *http.Request) (interface{}, error) {
	req := SaveApiReq{}
	err := Validate(r, &req)
	if err != nil {
		return JsonrespInterParamsErr(err), nil
	}

	err, api := entity.Api{}.FromIdSave(req.Api)
	if err != nil {
		return JsonrespInterServerErr(err), nil
	}
	req.Testing["ApiId"] = api.(entity.Api).ID
	err, _ = entity.Testing{}.FromIdSave(req.Testing)
	if err != nil {
		return JsonrespInterServerErr(err), nil
	}
	return JsonrespSuccess(nil), nil
})
