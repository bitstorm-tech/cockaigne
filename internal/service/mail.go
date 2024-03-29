package service

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	brevo "github.com/getbrevo/brevo-go/lib"
)

var br *brevo.APIClient

func init() {
	var cfg = brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", os.Getenv("BREVO_API_KEY"))
	br = brevo.NewAPIClient(cfg)
}

func SendAccountActivationMail(email string, baseUrl string) error {
	activationCode := rand.Intn(900000) + 100000

	err := SetActivationCode(email, activationCode)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"ActivationCode": activationCode,
		"ActivationUrl":  fmt.Sprintf("%s/signup-complete?email=%s&code=%d", baseUrl, email, activationCode),
	}

	templateId := os.Getenv("BREVO_ACTIVATE_ACCOUNT_TEMPLATE_ID")

	return sendMail(templateId, email, params)
}

func SendPasswordChangeEmail(email string, code string, baseUrl string) error {
	templateId := os.Getenv("BREVO_CHANGE_PASSWORD_TEMPLATE_ID")

	params := map[string]interface{}{
		"ChangePasswordUrl": fmt.Sprintf("%s/password-change/%s", baseUrl, code),
	}

	return sendMail(templateId, email, params)
}

func SendEmailChangeEmail(email string, code string, baseUrl string) error {
	templateId := os.Getenv("BREVO_CHANGE_EMAIL_TEMPLATE_ID")

	params := map[string]interface{}{
		"ChangeEmailUrl": fmt.Sprintf("%s/email-change/%s", baseUrl, code),
	}

	return sendMail(templateId, email, params)
}

func sendMail(templateId string, email string, params map[string]interface{}) error {
	templId, err := strconv.Atoi(templateId)
	if err != nil {
		return err
	}

	mail := brevo.SendSmtpEmail{
		To: []brevo.SendSmtpEmailTo{{
			Email: email,
		}},
		TemplateId: int64(templId),
		Params:     params,
	}

	_, _, err = br.TransactionalEmailsApi.SendTransacEmail(context.TODO(), mail)

	return err
}
