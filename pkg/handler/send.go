package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/authgear/authgear-server/pkg/util/httputil"
	"github.com/authgear/authgear-server/pkg/util/validation"
	sms_infra "github.com/authgear/authgear-sms-gateway/pkg/lib/infra/sms"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms"
)

type SendHandler struct {
	Logger     *slog.Logger
	SMSService *sms.SMSService
}

var RequestSchema = validation.NewMultipartSchema("SendRequestSchema")

var _ = RequestSchema.Add("SendRequestSchema", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"app_id": { "type": "string" },
		"to": { "type": "string" },
		"body": { "type": "string" },
		"app_id": { "type": "string" },
		"message_type": { "type": "string" },
		"template_name": { "type": "string" },
		"language_tag": { "type": "string" },
		"template_variables": { "$refs": "#/$defs/TemplateVariables" }
	},
	"required": ["to", "body", "template_name", "language_tag", "template_variables"]
}
`)
var _ = RequestSchema.Add("TemplateVariables", sms_infra.TemplateVariablesSchema)

func init() {
	RequestSchema.Instantiate()
}

func (h *SendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body RequestBody
	err := httputil.BindJSONBody(r, w, RequestSchema.Validator(), &body)
	if err != nil {
		h.write(w, &ResponseBody{
			Code:             CodeInvalidRequest,
			ErrorDescription: err.Error(),
		})
		return
	}

	h.Logger.Info(fmt.Sprintf("Attempt to send sms to %v. Body: %v. AppID: %v", body.To, body.Body, body.AppID))

	sendResult, err := h.SMSService.Send(
		body.AppID,
		&sms_infra.SendOptions{
			To:                body.To,
			Body:              body.Body,
			TemplateName:      body.TemplateName,
			LanguageTag:       body.LanguageTag,
			TemplateVariables: body.TemplateVariables,
		},
	)
	// TODO: handle err.
	if err != nil {
		h.write(w, &ResponseBody{
			Code:             CodeUnknownError,
			ErrorDescription: err.Error(),
		})
		return
	}

	h.write(w, &ResponseBody{
		Code:                       CodeOK,
		UnderlyingHTTPResponseBody: string(sendResult.ClientResponse),
		SegmentCount:               sendResult.SegmentCount,
	})
}

func (h *SendHandler) write(w http.ResponseWriter, body *ResponseBody) {
	statusCode := body.Code.HTTPStatusCode()
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := encoder.Encode(body)
	if err != nil {
		panic(err)
	}
}
