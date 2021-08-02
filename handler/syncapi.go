package handler

import (
	"apitesting/entity"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type Entity struct {
	Api     entity.Api
	Testing entity.Testing
}

type SyncApiReq struct {
	Project  string   `valid:"Required"`
	Entities []Entity `valid:"Required"`
}

var SyncApi = JsonHandlerFunc(func(r *http.Request) Jsonresp {
	req := SyncApiReq{}
	err := Validate(r, &req)
	if err != nil {
		return JsonrespInterParamsErr(err)
	}

	err = entity.GormDB().Transaction(func(tx *gorm.DB) error {
		err := tx.Where("project = ?", req.Project).Delete(entity.Api{}).Error
		if err != nil {
			return fmt.Errorf("SyncApi Api Delete %s", err.Error())
		}
		err = tx.Where("project = ? and user_id = 0", req.Project).Delete(entity.Testing{}).Error
		if err != nil {
			return fmt.Errorf("SyncApi Testing Delete %s", err.Error())
		}
		for _, e := range req.Entities {
			err = tx.Create(&e.Api).Error
			if err != nil {
				return fmt.Errorf("SyncApi Apis Create %s", err.Error())
			}
			e.Testing.ApiId = e.Api.ID
			err = tx.Create(&e.Testing).Error
			if err != nil {
				return fmt.Errorf("SyncApi Testings Create %s", err.Error())
			}
		}
		return nil
	})

	if err != nil {
		return JsonrespInterServerErr(err)
	}
	return JsonrespSuccess(nil)
})
