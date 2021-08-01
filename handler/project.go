package handler

import (
	"apitesting/entity"
	"fmt"
	"net/http"
)

var ApiProjects = JsonHandlerFunc(func(r *http.Request) (interface{}, error) {
	ret := []string{}
	err := entity.GormDB().Model(&entity.Api{}).Select("project").Group("project").Find(&ret).Error
	if err != nil {
		return JsonrespInterServerErr(fmt.Errorf("ApiProjects Group Find %s", err.Error())), nil
	}
	return JsonrespSuccess(ret), nil
})

var TestingProjects = JsonHandlerFunc(func(r *http.Request) (interface{}, error) {
	ret := []string{}
	err := entity.GormDB().Model(&entity.Testing{}).Select("project").Group("project").Find(&ret).Error
	if err != nil {
		return JsonrespInterServerErr(fmt.Errorf("TestingProjects Group Find %s", err.Error())), nil
	}
	return JsonrespSuccess(ret), nil
})
