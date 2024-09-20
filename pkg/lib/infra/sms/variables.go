package sms

type TemplateVariables struct {
	AppName     string `json:"app_name"`
	ClientID    string `json:"client_id"`
	Code        string `json:"code"`
	Email       string `json:"email"`
	HasPassword bool   `json:"has_password"`
	Host        string `json:"host"`
	Link        string `json:"link"`
	Password    string `json:"password"`
	Phone       string `json:"phone"`
	State       string `json:"state"`
	UILocales   string `json:"ui_locales"`
	URL         string `json:"url"`
	XState      string `json:"x_state"`
}

var TemplateVariablesSchema = `{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"app_name": { "type": "string" },
		"client_id": { "type": "string" },
		"code": { "type": "string" },
		"email": { "type": "string" },
		"has_password": { "type": "boolean" },
		"host": { "type": "string" },
		"link": { "type": "string" },
		"password": { "type": "string" },
		"phone": { "type": "string" },
		"state": { "type": "string" },
		"ui_locales": { "type": "string" },
		"url": { "type": "string" },
		"x_state": { "type": "string" }
	},
	"required": []
}`
