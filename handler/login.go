package handler

import (
	"fmt"
	"net/http"

	"apitesting/entity"
	"apitesting/token"
)

type LoginReq struct {
	ID       uint
	Name     string `validate:"required"`
	Password string `validate:"required"`
}

var TokenCodec = token.NewTokenCodec([]byte("AQdJz7ZGbv8GepLS"), []byte{'_'})

var Login = JsonHandlerFunc(func(r *http.Request) (interface{}, error) {
	req := LoginReq{}
	err := Validate(r, &req)
	if err != nil {
		return JsonrespInterParamsErr(err), nil
	}

	u := &entity.User{}
	exist, err := entity.GormFirst(u, &entity.User{Name: req.Name}, "Name")
	if err != nil {
		return JsonrespInterServerErr(fmt.Errorf("Login GormFirst %s", err.Error())), nil
	}
	if !exist {
		return JsonrespInterOperationNotAllowed(fmt.Errorf("User %s not exist", req.Name)), nil
	}
	if req.Password != u.Password {
		return JsonrespInterOperationNotAllowed(fmt.Errorf("Bad Password")), nil
	}
	t := TokenCodec.Bytes(u.ID, u.Name)
	return map[string]interface{}{
		"Token": string(t),
		"Id":    u.ID,
		"Name":  u.Name,
	}, nil
})
