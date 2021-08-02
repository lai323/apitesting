package handler

import (
	"apitesting/entity"
	"net/http"
)

type DelReq struct {
	ID uint `validate:"required"`
}

func NewDelFunc(t entity.FromIdDeleter) func(*http.Request) Jsonresp {
	return func(r *http.Request) Jsonresp {
		req := DelReq{}
		err := Validate(r, &req)
		if err != nil {
			return JsonrespInterParamsErr(err)
		}
		err = t.FromIdDel(req.ID)
		if err != nil {
			return JsonrespInterServerErr(err)
		}
		return JsonrespSuccess(nil)
	}
}

var (
	DelEnv     = JsonHandlerFunc(NewDelFunc(entity.Env{}))
	DelUser    = JsonHandlerFunc(NewDelFunc(entity.User{}))
	DelApi     = JsonHandlerFunc(NewDelFunc(entity.Api{}))
	DelTesting = JsonHandlerFunc(NewDelFunc(entity.Testing{}))
)
