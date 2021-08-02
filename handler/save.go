package handler

import (
	"apitesting/entity"
	"net/http"
)

type SaveReq struct {
	Entity map[string]interface{}
}

func NewSaveFunc(t entity.FromIdSaver) func(*http.Request) Jsonresp {
	return func(r *http.Request) Jsonresp {
		req := SaveReq{}
		err := Validate(r, &req)
		if err != nil {
			return JsonrespInterParamsErr(err)
		}

		err, _ = t.FromIdSave(req.Entity)
		if err != nil {
			return JsonrespInterServerErr(err)
		}
		return JsonrespSuccess(nil)
	}
}

var (
	SaveEnv     = JsonHandlerFunc(NewSaveFunc(entity.Env{}))
	SaveUser    = JsonHandlerFunc(NewSaveFunc(entity.User{}))
	SaveTesting = JsonHandlerFunc(NewSaveFunc(entity.Testing{}))
)
