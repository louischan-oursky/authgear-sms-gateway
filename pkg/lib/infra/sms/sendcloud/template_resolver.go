package sendcloud

import (
	"errors"
	"fmt"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
)

type ISendCloudTemplateResolver interface {
	Resolve(templateName string, languageTag string) (*SendCloudTemplate, error)
}

type TemplateID string

type TemplateMessageType string

type AuthgearTemplateName string

type AuthgearLanguage string

type SendCloudTemplate struct {
	TemplateID      TemplateID
	TemplateMsgType TemplateMessageType
}

func NewSendCloudTemplate(template *config.SendCloudTemplate) *SendCloudTemplate {
	return &SendCloudTemplate{
		TemplateID:      TemplateID(template.TemplateID),
		TemplateMsgType: TemplateMessageType(template.TemplateMsgType),
	}
}

type ByLanguage struct {
	AuthgearLanguage AuthgearLanguage
	Template         *SendCloudTemplate
}

func NewByLanguage(authgearLanguage AuthgearLanguage, template *SendCloudTemplate) *ByLanguage {
	return &ByLanguage{
		AuthgearLanguage: authgearLanguage,
		Template:         template,
	}
}

type SendCloudTemplateAssignment struct {
	AuthgearTemplateName AuthgearTemplateName
	DefaultTemplate      *SendCloudTemplate
	ByLanguages          []*ByLanguage
	ByLanguageMap        map[AuthgearLanguage]*ByLanguage
}

func NewSendCloudTemplateAssignment(templateAssignment *config.SendCloudTemplateAssignment, templateIDMap map[TemplateID]*SendCloudTemplate) (*SendCloudTemplateAssignment, error) {
	byLanguages := make([]*ByLanguage, len(templateAssignment.ByLanguages))
	byLanguageMap := make(map[AuthgearLanguage]*ByLanguage)
	for i, byLanguage := range templateAssignment.ByLanguages {
		template := templateIDMap[TemplateID(byLanguage.TemplateID)]
		if template == nil {
			return nil, errors.New(fmt.Sprintf("Cannot find template with id %v", byLanguage.TemplateID))
		}
		b := NewByLanguage(AuthgearLanguage(byLanguage.AuthgearLanguage), template)
		byLanguages[i] = b
		byLanguageMap[b.AuthgearLanguage] = b
	}

	defaultTemplate := templateIDMap[TemplateID(templateAssignment.DefaultTemplateID)]
	if defaultTemplate == nil {
		return nil, errors.New(fmt.Sprintf("Cannot find template with id %v", templateAssignment.DefaultTemplateID))
	}

	return &SendCloudTemplateAssignment{
		AuthgearTemplateName: AuthgearTemplateName(templateAssignment.AuthgearTemplateName),
		DefaultTemplate:      defaultTemplate,
		ByLanguages:          byLanguages,
		ByLanguageMap:        byLanguageMap,
	}, nil
}

type SendCloudTemplateResolver struct {
	templates                           []*SendCloudTemplate
	templateIDMap                       map[TemplateID]*SendCloudTemplate
	templateAssignments                 []*SendCloudTemplateAssignment
	templateAssignmentMapByTemplateName map[AuthgearTemplateName]*SendCloudTemplateAssignment
}

var _ ISendCloudTemplateResolver = &SendCloudTemplateResolver{}

func NewSendCloudTemplateResolver(
	templates []*config.SendCloudTemplate,
	templateAssignments []*config.SendCloudTemplateAssignment,
) (*SendCloudTemplateResolver, error) {
	ts := make([]*SendCloudTemplate, len(templates))
	templateIDMap := make(map[TemplateID]*SendCloudTemplate)
	for i, template := range templates {
		t := NewSendCloudTemplate(template)
		ts[i] = t
		templateIDMap[t.TemplateID] = t
	}

	tas := make([]*SendCloudTemplateAssignment, len(templateAssignments))
	templateAssignmentMapByTemplateName := make(map[AuthgearTemplateName]*SendCloudTemplateAssignment)
	for i, templateAssignment := range templateAssignments {
		ta, err := NewSendCloudTemplateAssignment(templateAssignment, templateIDMap)
		if err != nil {
			return nil, err
		}
		tas[i] = ta
		templateAssignmentMapByTemplateName[ta.AuthgearTemplateName] = ta
	}

	return &SendCloudTemplateResolver{
		templates:                           ts,
		templateIDMap:                       templateIDMap,
		templateAssignments:                 tas,
		templateAssignmentMapByTemplateName: templateAssignmentMapByTemplateName,
	}, nil
}

func (s *SendCloudTemplateResolver) Resolve(templateName string, languageTag string) (*SendCloudTemplate, error) {
	templateAssignment := s.templateAssignmentMapByTemplateName[AuthgearTemplateName(templateName)]
	if templateAssignment == nil {
		return nil, errors.New(fmt.Sprintf("Could not found template assignment from template name %v", templateName))
	}
	byLanguage := templateAssignment.ByLanguageMap[AuthgearLanguage(languageTag)]

	if byLanguage == nil {
		return templateAssignment.DefaultTemplate, nil
	}
	return byLanguage.Template, nil
}
