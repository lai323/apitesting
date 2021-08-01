package handler

import (
	"fmt"
	"net/http"

	"apitesting/entity"
)

type QueryReq struct {
	Filter entity.Filter
}

func NewSelectFunc(t entity.FilterSelecter) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		req := QueryReq{}
		err := Validate(r, &req)
		if err != nil {
			return JsonrespInterParamsErr(err), nil
		}

		rows, err := t.SelectFromFilter(req.Filter)
		if err != nil {
			return JsonrespInterServerErr(fmt.Errorf("SelectFromFilter %s SelectFromFilter %s", t.TableName(), err.Error())), nil
		}
		total, err := t.CountFromFilter(req.Filter)
		if err != nil {
			return JsonrespInterServerErr(fmt.Errorf("SelectFromFilter %s TotalCount %s", t.TableName(), err.Error())), nil

		}
		return JsonrespSuccess(map[string]interface{}{
			"rows":  rows,
			"total": total,
		}), err
	}
}

var (
	SelectEnv     = JsonHandlerFunc(NewSelectFunc(entity.Env{}))
	SelectUser    = JsonHandlerFunc(NewSelectFunc(entity.User{}))
	SelectApi     = JsonHandlerFunc(NewSelectFunc(entity.Api{}))
	SelectTesting = JsonHandlerFunc(NewSelectFunc(entity.Testing{}))
)
