package handler

import (
	"apitesting/entity"
	"fmt"
	"net/http"
)

var ApiProjects = JsonHandlerFunc(func(r *http.Request) Jsonresp {
	ret := []string{}
	err := entity.GormDB().Model(&entity.Api{}).Select("project").Group("project").Find(&ret).Error
	if err != nil {
		return JsonrespInterServerErr(fmt.Errorf("ApiProjects Group Find %s", err.Error()))
	}
	return JsonrespSuccess(ret)
})

var TestingProjects = JsonHandlerFunc(func(r *http.Request) Jsonresp {
	ret := []string{}
	err := entity.GormDB().Model(&entity.Testing{}).Select("project").Group("project").Find(&ret).Error
	if err != nil {
		return JsonrespInterServerErr(fmt.Errorf("TestingProjects Group Find %s", err.Error()))
	}
	return JsonrespSuccess(ret)
})
