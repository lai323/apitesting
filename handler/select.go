package handler

import (
	"fmt"
	"net/http"

	"apitesting/entity"
)

type QueryReq struct {
	Filter entity.Filter
}

func NewSelectFunc(t entity.FilterSelecter) func(*http.Request) Jsonresp {
	return func(r *http.Request) Jsonresp {
		req := QueryReq{}
		err := Validate(r, &req)
		if err != nil {
			return JsonrespInterParamsErr(err)
		}

		rows, err := t.SelectFromFilter(req.Filter)
		if err != nil {
			return JsonrespInterServerErr(fmt.Errorf("SelectFromFilter %s SelectFromFilter %s", t.TableName(), err.Error()))
		}
		total, err := t.CountFromFilter(req.Filter)
		if err != nil {
			return JsonrespInterServerErr(fmt.Errorf("SelectFromFilter %s TotalCount %s", t.TableName(), err.Error()))

		}
		return JsonrespSuccess(map[string]interface{}{
			"rows":  rows,
			"total": total,
		})
	}
}

var (
	SelectEnv     = JsonHandlerFunc(NewSelectFunc(entity.Env{}))
	SelectUser    = JsonHandlerFunc(NewSelectFunc(entity.User{}))
	SelectApi     = JsonHandlerFunc(NewSelectFunc(entity.Api{}))
	SelectTesting = JsonHandlerFunc(NewSelectFunc(entity.Testing{}))
)
