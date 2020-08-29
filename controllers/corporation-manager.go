package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"

	"github.com/zengchen1024/cla-server/models"
)

type CorporationManagerController struct {
	beego.Controller
}

// @Title add corporation manager
// @Description add corporation manager
// @Param	body		body 	models.CorporationManagerCreateOption	true		"body for corporation manager"
// @Success 201 {int} map
// @Failure 403 body is empty
// @router / [post]
func (this *CorporationManagerController) Post() {
	var statusCode = 201
	var reason error
	var body interface{}

	defer func() {
		sendResponse(&this.Controller, statusCode, reason, body)
	}()

	var info models.CorporationManagerCreateOption
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &info); err != nil {
		reason = err
		statusCode = 400
		return
	}

	if err := (&info).Create(); err != nil {
		reason = err
		statusCode = 500
		return
	}

	body = "add manager successfully"
}

// @Title authenticate corporation manager
// @Description authenticate corporation manager
// @Param	body		body 	models.CorporationManagerAuthentication	true		"body for corporation manager info"
// @Success 201 {int} map
// @Failure 403 body is empty
// @router /auth [post]
func (this *CorporationManagerController) Authenticate() {
	var statusCode = 201
	var reason error
	var body interface{}

	defer func() {
		sendResponse(&this.Controller, statusCode, reason, body)
	}()

	var info models.CorporationManagerAuthentication
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &info); err != nil {
		reason = err
		statusCode = 400
		return
	}

	if v, err := (&info).Authenticate(); err != nil {
		reason = err
		statusCode = 500
		return
	} else {
		// body = "authenticate successfully"
		body = v
	}
}

// @Title Reset password
// @Description reset password
// @Param	body		body 	models.CorporationManagerResetPassword	true		"body for corporation manager info"
// @Success 201 {int} map
// @Failure 403 body is empty
// @router / [put]
func (this *CorporationManagerController) Update() {
	var statusCode = 202
	var reason error
	var body interface{}

	defer func() {
		sendResponse(&this.Controller, statusCode, reason, body)
	}()

	var info models.CorporationManagerResetPassword
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &info); err != nil {
		reason = err
		statusCode = 400
		return
	}

	if err := (&info).Reset(); err != nil {
		reason = err
		statusCode = 500
		return
	}

	body = "reset password successfully"
}