package handler

import (
	"apitesting/entity"
	"net/http"
)

type DelReq struct {
	ID uint `validate:"required"`
}

func NewDelFunc(t entity.FromIdDeleter) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		req := DelReq{}
		err := Validate(r, &req)
		if err != nil {
			return JsonrespInterParamsErr(err), nil
		}
		err = t.FromIdDel(req.ID)
		if err != nil {
			return JsonrespInterServerErr(err), nil
		}
		return JsonrespSuccess(nil), nil
	}
}

var (
	DelEnv     = JsonHandlerFunc(NewDelFunc(entity.Env{}))
	DelUser    = JsonHandlerFunc(NewDelFunc(entity.User{}))
	DelApi     = JsonHandlerFunc(NewDelFunc(entity.Api{}))
	DelTesting = JsonHandlerFunc(NewDelFunc(entity.Testing{}))
)
