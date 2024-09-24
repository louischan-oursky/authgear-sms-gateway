package sms

import (
	"fmt"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/infra/sms/accessyou/apis"
)

type AccessYouClient struct {
	BaseUrl   string
	Client    *http.Client
	AccountNo string
	User      string
	Pwd       string
	Sender    string
	Logger    *slog.Logger
}

func NewAccessYouClient(
	baseUrl string,
	accountNo string,
	user string,
	pwd string,
	sender string,
	logger *slog.Logger,
) *AccessYouClient {
	if baseUrl == "" {
		baseUrl = "http://sms.accessyou-anyip.com"
	}
	return &AccessYouClient{
		BaseUrl:   baseUrl,
		Client:    &http.Client{},
		AccountNo: accountNo,
		User:      user,
		Pwd:       pwd,
		Sender:    sender,
		Logger:    logger,
	}
}

var plusHyphensRegexp = regexp.MustCompile(`[\+\-]+`)

func fixPhoneNumber(phoneNumber string) string {
	// Access you phone number should have no + and -
	return plusHyphensRegexp.ReplaceAllString(phoneNumber, "")
}

func (n *AccessYouClient) Send(options *SendOptions) (*SendResult, error) {
	to := fixPhoneNumber(string(options.To))

	respData, sendSMSResponse, err := apis.SendSMS(
		n.Client,
		n.BaseUrl,
		n.AccountNo,
		n.User,
		n.Pwd,
		n.Sender,
		to,
		options.Body,
	)
	n.Logger.Info(fmt.Sprintf("%v", sendSMSResponse))

	if err != nil {
		n.Logger.Error(fmt.Sprintf("%v", err))
		return nil, err
	}

	return &SendResult{
		ClientResponse: respData,
		Success:        sendSMSResponse.Status == "100",
	}, nil
}

var _ RawClient = &AccessYouClient{}
