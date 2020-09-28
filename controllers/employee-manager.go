package controllers

import (
	"fmt"

	"github.com/astaxie/beego"

	"github.com/opensourceways/app-cla-server/dbmodels"
	"github.com/opensourceways/app-cla-server/models"
	"github.com/opensourceways/app-cla-server/util"
)

type EmployeeManagerController struct {
	beego.Controller
}

func (this *EmployeeManagerController) Prepare() {
	apiPrepare(&this.Controller, []string{PermissionCorporAdmin}, nil)
}

// @Title Post
// @Description add employee managers
// @Param	:cla_org_id	path 	string					true		"cla org id"
// @Param	body		body 	models.EmployeeManagerCreateOption	true		"body for employee manager"
// @Success 201 {int} map
// @router / [post]
func (this *EmployeeManagerController) Post() {
	this.addOrDeleteManagers(true)
}

// @Title Delete
// @Description delete employee manager
// @Param	:cla_org_id	path 	string					true		"cla org id"
// @Param	body		body 	models.EmployeeManagerCreateOption	true		"body for employee manager"
// @Success 204 {string} delete success!
// @router / [delete]
func (this *EmployeeManagerController) Delete() {
	this.addOrDeleteManagers(false)
}

// @Title GetAll
// @Description get all employee managers
// @Param	:cla_org_id	path 	string					true		"cla org id"
// @Success 200 {object} dbmodels.CorporationManagerListResult
// @router /:cla_org_id [get]
func (this *EmployeeManagerController) GetAll() {
	var statusCode = 0
	var errCode = 0
	var reason error
	var body interface{}

	defer func() {
		sendResponse(&this.Controller, statusCode, errCode, reason, body, "list employee managers")
	}()

	claOrgID, err := fetchStringParameter(&this.Controller, ":cla_org_id")
	if err != nil {
		reason = err
		errCode = util.ErrInvalidParameter
		statusCode = 400
		return
	}

	corpEmail, err := getApiAccessUser(&this.Controller)
	if err != nil {
		reason = err
		statusCode = 500
		return
	}

	r, err := models.ListCorporationManagers(
		claOrgID, corpEmail, dbmodels.RoleManager,
	)
	if err != nil {
		reason = err
		statusCode, errCode = convertDBError(err)
		return
	}

	body = r
}

func (this *EmployeeManagerController) addOrDeleteManagers(toAdd bool) {
	var statusCode = 0
	var errCode = 0
	var reason error
	var body interface{}

	defer func() {
		op := "add"
		if !toAdd {
			op = "delete"
		}
		body = fmt.Sprintf("%s employee manager successfully", op)

		sendResponse(&this.Controller, statusCode, errCode, reason, body, fmt.Sprintf("%s employee managers", op))
	}()

	claOrgID, corpEmail, err := parseCorpManagerUser(&this.Controller)
	if err != nil {
		reason = err
		errCode = util.ErrUnknownToken
		statusCode = 401
		return
	}

	//TODO: check if the cla org has the adminstrator, otherwise he/she can't do this operation.

	var info models.EmployeeManagerCreateOption
	if err := fetchInputPayload(&this.Controller, &info); err != nil {
		reason = err
		errCode = util.ErrInvalidParameter
		statusCode = 400
		return
	}

	if err := (&info).Validate(); err != nil {
		reason = err
		errCode = util.ErrInvalidParameter
		statusCode = 400
		return
	}

	if !isSameCorp(corpEmail, info.Emails[0]) {
		reason = fmt.Errorf("not same corp")
		errCode = util.ErrNotSameCorp
		statusCode = 400
		return
	}

	if toAdd {
		err = (&info).Create(claOrgID)
	} else {
		err = (&info).Delete(claOrgID)
	}
	if err != nil {
		reason = err
		statusCode, errCode = convertDBError(err)
	}
}
