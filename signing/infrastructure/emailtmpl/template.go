package emailtmpl

import (
	"fmt"
	"text/template"

	"github.com/opensourceways/app-cla-server/signing/domain/emailservice"
	"github.com/opensourceways/app-cla-server/util"
)

const (
	TmplCorporationSigning    = "corporation signing"
	TmplIndividualSigning     = "individual signing"
	TmplEmployeeSigning       = "employee signing"
	TmplNotifyingManager      = "notifying manager"
	TmplVerificationCode      = "verification code"
	TmplAddingCorpEmailDomain = "adding corp email domain"
	TmplAddingCorpAdmin       = "adding corp admin"
	TmplAddingCorpManager     = "adding corp manager"
	TmplRemovingCorpManager   = "removing corp manager"
	TmplActivatingEmployee    = "activating employee"
	TmplInactivaingEmployee   = "inactivating employee"
	TmplRemovingingEmployee   = "removing employee"
	TmplPasswordRetrieval     = "password retrieval"
	TmplEmailVerification     = "email verification"
)

var msgTmpl = map[string]*template.Template{}

type EmailMessage = emailservice.EmailMessage

func Init() error {
	items := map[string]string{
		TmplCorporationSigning:    "./conf/email-template/corporation-signing.tmpl",
		TmplIndividualSigning:     "./conf/email-template/individual-signing.tmpl",
		TmplEmployeeSigning:       "./conf/email-template/employee-signing.tmpl",
		TmplNotifyingManager:      "./conf/email-template/notifying-corp-manager.tmpl",
		TmplVerificationCode:      "./conf/email-template/verification-code.tmpl",
		TmplAddingCorpEmailDomain: "./conf/email-template/adding-corp-email-domain.tmpl",
		TmplAddingCorpAdmin:       "./conf/email-template/adding-corp-admin.tmpl",
		TmplAddingCorpManager:     "./conf/email-template/adding-corp-manager.tmpl",
		TmplRemovingCorpManager:   "./conf/email-template/removing-corp-manager.tmpl",
		TmplActivatingEmployee:    "./conf/email-template/activating-employee.tmpl",
		TmplInactivaingEmployee:   "./conf/email-template/inactivating-employee.tmpl",
		TmplRemovingingEmployee:   "./conf/email-template/removing-employee.tmpl",
		TmplPasswordRetrieval:     "./conf/email-template/password-retrieval.tmpl",
		TmplEmailVerification:     "./conf/email-template/email_verification.tmpl",
	}

	for name, path := range items {
		tmpl, err := util.NewTemplate(name, path)
		if err != nil {
			return err
		}
		msgTmpl[name] = tmpl
	}

	return nil
}

func findTmpl(name string) *template.Template {
	v, ok := msgTmpl[name]
	if ok {
		return v
	}
	return nil
}

func genEmailMsg(tmplName string, data interface{}) (msg EmailMessage, err error) {
	tmpl := findTmpl(tmplName)
	if tmpl == nil {
		err = fmt.Errorf("failed to generate email msg: didn't find msg template: %s", tmplName)

		return
	}

	msg.Content, err = util.RenderTemplate(tmpl, data)

	return
}

type CorporationSigning struct {
	Org         string
	Date        string
	AdminName   string
	ProjectURL  string
	SigningInfo string
}

func (data *CorporationSigning) GenEmailMsg() (EmailMessage, error) {
	return genEmailMsg(TmplCorporationSigning, data)
}

type IndividualSigning struct {
	Name string
}

func (data *IndividualSigning) GenEmailMsg() (EmailMessage, error) {
	return genEmailMsg(TmplIndividualSigning, data)
}

type VerificationCode struct {
	Email      string
	Org        string
	Code       string
	ProjectURL string
}

func (data *VerificationCode) GenEmailMsg() (EmailMessage, error) {
	return genEmailMsg(TmplVerificationCode, data)
}

type AddingCorpEmailDomain struct {
	Corp       string
	Org        string
	Code       string
	ProjectURL string
}

func (cse AddingCorpEmailDomain) GenEmailMsg() (EmailMessage, error) {
	return genEmailMsg(TmplAddingCorpEmailDomain, cse)
}

type AddingCorpManager struct {
	Admin            bool
	ID               string
	User             string
	Email            string
	Password         []byte
	Org              string
	ProjectURL       string
	URLOfCLAPlatform string
}

func (data *AddingCorpManager) GenEmailMsg() (msg EmailMessage, err error) {
	if data.Admin {
		msg, err = genEmailMsg(TmplAddingCorpAdmin, data)
	} else {
		msg, err = genEmailMsg(TmplAddingCorpManager, data)
	}

	msg.HasSecret = true

	return
}

type RemovingCorpManager struct {
	User       string
	Org        string
	ProjectURL string
}

func (data *RemovingCorpManager) GenEmailMsg() (EmailMessage, error) {
	return genEmailMsg(TmplRemovingCorpManager, data)
}

type EmployeeSigning struct {
	Name       string
	Org        string
	ProjectURL string
	Managers   string
}

func (data *EmployeeSigning) GenEmailMsg() (EmailMessage, error) {
	return genEmailMsg(TmplEmployeeSigning, data)
}

type NotifyingManager struct {
	EmployeeEmail    string
	ProjectURL       string
	URLOfCLAPlatform string
	Org              string
}

func (data *NotifyingManager) GenEmailMsg() (EmailMessage, error) {
	return genEmailMsg(TmplNotifyingManager, data)
}

type EmployeeNotification struct {
	Removing bool
	Active   bool
	Inactive bool

	Name       string
	ProjectURL string
	Manager    string
	Org        string
}

func (data *EmployeeNotification) GenEmailMsg() (EmailMessage, error) {
	if data.Active {
		return genEmailMsg(TmplActivatingEmployee, data)
	}

	if data.Inactive {
		return genEmailMsg(TmplInactivaingEmployee, data)
	}

	if data.Removing {
		return genEmailMsg(TmplRemovingingEmployee, data)
	}

	return EmailMessage{}, fmt.Errorf("do nothing")
}

type PasswordRetrieval struct {
	Timeout      int64
	Org          string
	ResetURL     string
	RetrievalURL string
}

func (p PasswordRetrieval) GenEmailMsg() (EmailMessage, error) {
	msg, err := genEmailMsg(TmplPasswordRetrieval, p)
	if err != nil {
		return msg, err
	}

	//adapter send html tmpl content
	msg.MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	return msg, nil
}

type EmailVerification struct {
	Code string
}

func (data *EmailVerification) GenEmailMsg() (EmailMessage, error) {
	return genEmailMsg(TmplEmailVerification, data)
}
