package handler

import (
	"apitesting/entity"
	"net/http"
)

type SaveReq struct {
	Entity map[string]interface{}
}

func NewSaveFunc(t entity.FromIdSaver) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		req := SaveReq{}
		err := Validate(r, &req)
		if err != nil {
			return JsonrespInterParamsErr(err), nil
		}

		err, _ = t.FromIdSave(req.Entity)
		if err != nil {
			return JsonrespInterServerErr(err), nil
		}
		return JsonrespSuccess(nil), nil
	}
}

var (
	SaveEnv     = JsonHandlerFunc(NewSaveFunc(entity.Env{}))
	SaveUser    = JsonHandlerFunc(NewSaveFunc(entity.User{}))
	SaveTesting = JsonHandlerFunc(NewSaveFunc(entity.Testing{}))
)
