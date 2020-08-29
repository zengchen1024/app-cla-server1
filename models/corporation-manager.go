package models

import "github.com/zengchen1024/cla-server/dbmodels"

const (
	RoleAdmin   = "admin"
	RoleManager = "manager"
)

type CorporationManagerCreateOption struct {
	CLAOrgID string `json:"cla_org_id"`
	Email    string `json:"email"`
}

func (this *CorporationManagerCreateOption) Create() error {
	pw := "123456"
	opt := []dbmodels.CorporationManagerCreateOption{
		{
			Role:          RoleAdmin,
			Email:         this.Email,
			Password:      pw,
			CorporationID: emailSuffixToKey(this.Email),
		},
	}
	return dbmodels.GetDB().AddCorporationManager(this.CLAOrgID, opt, 1)
}

type CorporationManagerAuthentication struct {
	Platform string `json:"platform"`
	OrgID    string `json:"org_id"`
	RepoID   string `json:"repo_id"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func (this CorporationManagerAuthentication) Authenticate() (dbmodels.CorporationManagerCheckResult, error) {
	opt := dbmodels.CorporationManagerCheckInfo{
		Platform: this.Platform,
		OrgID:    this.OrgID,
		RepoID:   this.RepoID,
		User:     this.User,
		Password: this.Password,
	}

	return dbmodels.GetDB().CheckCorporationManagerExist(opt)
}

type CorporationManagerResetPassword struct {
	CLAOrgID    string `json:"cla_org_id"`
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (this CorporationManagerResetPassword) Reset() error {
	opt := dbmodels.CorporationManagerResetPassword{
		Email:       this.Email,
		OldPassword: this.OldPassword,
		NewPassword: this.NewPassword,
	}

	return dbmodels.GetDB().ResetCorporationManagerPassword(this.CLAOrgID, opt)
}

type CorporationManagerListOption struct {
	CLAOrgID string `json:"cla_org_id"`
	Role     string `json:"role"`
	Email    string `json:"email"`
}

func (this CorporationManagerListOption) List() ([]dbmodels.CorporationManagerListResult, error) {
	opt := dbmodels.CorporationManagerListOption{
		Role:          this.Role,
		CorporationID: emailSuffixToKey(this.Email),
	}
	return dbmodels.GetDB().ListCorporationManager(this.CLAOrgID, opt)
}