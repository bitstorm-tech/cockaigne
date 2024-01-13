package service

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	brevo "github.com/getbrevo/brevo-go/lib"
	"go.uber.org/zap"
)

var br *brevo.APIClient

func init() {
	var cfg = brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", os.Getenv("BREVO_API_KEY"))
	br = brevo.NewAPIClient(cfg)
}

type ActivationMailParams struct {
	ActivationCode int
	ActivationUrl  string
}

func SendAccountActivationMail(email string, baseUrl string) error {
	activationCode := rand.Intn(900000) + 100000

	err := SetActivationCode(email, activationCode)
	if err != nil {
		return err
	}

	params := interface{}(ActivationMailParams{
		ActivationCode: activationCode,
		ActivationUrl:  fmt.Sprintf("%s/signup-complete?email=%s&code=%d", baseUrl, email, activationCode),
	})

	templateIdString := os.Getenv("BREVO_SEND_ACTIVATION_CODE_TEMPLATE_ID")
	templateId, err := strconv.Atoi(templateIdString)
	if err != nil {
		zap.L().Sugar().Errorf("can't convert template ID string '%s': %+v", templateIdString, err)
		return err
	}

	mail := brevo.SendSmtpEmail{
		To: []brevo.SendSmtpEmailTo{{
			Email: email,
		}},
		TemplateId: int64(templateId),
		Params:     &params,
	}

	_, _, err = br.TransactionalEmailsApi.SendTransacEmail(context.TODO(), mail)

	return err
}
