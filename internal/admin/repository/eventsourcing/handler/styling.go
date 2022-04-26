package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/caos/logging"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"

	"github.com/zitadel/zitadel/internal/domain"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/eventstore/v1/query"
	"github.com/zitadel/zitadel/internal/eventstore/v1/spooler"
	iam_es_model "github.com/zitadel/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/zitadel/zitadel/internal/iam/repository/view/model"
	"github.com/zitadel/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/zitadel/zitadel/internal/static"
)

const (
	stylingTable = "adminapi.styling"
)

type Styling struct {
	handler
	static       static.Storage
	subscription *v1.Subscription
	resourceUrl  string
}

func newStyling(handler handler, static static.Storage, localDevMode bool) *Styling {
	h := &Styling{
		handler: handler,
		static:  static,
	}
	prefix := ""
	if localDevMode {
		prefix = "/login"
	}
	h.resourceUrl = prefix + "/resources/dynamic" //TODO: ?

	h.subscribe()

	return h
}

func (m *Styling) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

func (m *Styling) ViewModel() string {
	return stylingTable
}

func (m *Styling) Subscription() *v1.Subscription {
	return m.subscription
}

func (_ *Styling) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.OrgAggregate, iam_es_model.IAMAggregate}
}

func (m *Styling) CurrentSequence() (uint64, error) {
	sequence, err := m.view.GetLatestStylingSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (m *Styling) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := m.view.GetLatestStylingSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(m.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *Styling) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate, iam_es_model.IAMAggregate:
		err = m.processLabelPolicy(event)
	}
	return err
}

func (m *Styling) processLabelPolicy(event *es_models.Event) (err error) {
	policy := new(iam_model.LabelPolicyView)
	switch event.Type {
	case iam_es_model.LabelPolicyAdded, model.LabelPolicyAdded:
		err = policy.AppendEvent(event)
	case iam_es_model.LabelPolicyChanged, model.LabelPolicyChanged,
		iam_es_model.LabelPolicyLogoAdded, model.LabelPolicyLogoAdded,
		iam_es_model.LabelPolicyLogoRemoved, model.LabelPolicyLogoRemoved,
		iam_es_model.LabelPolicyIconAdded, model.LabelPolicyIconAdded,
		iam_es_model.LabelPolicyIconRemoved, model.LabelPolicyIconRemoved,
		iam_es_model.LabelPolicyLogoDarkAdded, model.LabelPolicyLogoDarkAdded,
		iam_es_model.LabelPolicyLogoDarkRemoved, model.LabelPolicyLogoDarkRemoved,
		iam_es_model.LabelPolicyIconDarkAdded, model.LabelPolicyIconDarkAdded,
		iam_es_model.LabelPolicyIconDarkRemoved, model.LabelPolicyIconDarkRemoved,
		iam_es_model.LabelPolicyFontAdded, model.LabelPolicyFontAdded,
		iam_es_model.LabelPolicyFontRemoved, model.LabelPolicyFontRemoved,
		iam_es_model.LabelPolicyAssetsRemoved, model.LabelPolicyAssetsRemoved:
		policy, err = m.view.StylingByAggregateIDAndState(event.AggregateID, int32(domain.LabelPolicyStatePreview))
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)

	case iam_es_model.LabelPolicyActivated, model.LabelPolicyActivated:
		policy, err = m.view.StylingByAggregateIDAndState(event.AggregateID, int32(domain.LabelPolicyStatePreview))
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
		if err != nil {
			return err
		}
		err = m.generateStylingFile(policy)
	default:
		return m.view.ProcessedStylingSequence(event)
	}
	if err != nil {
		return err
	}
	return m.view.PutStyling(policy, event)
}

func (m *Styling) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-2m9fs", "id", event.AggregateID).WithError(err).Warn("something went wrong in label policy handler")
	return spooler.HandleError(event, err, m.view.GetLatestStylingFailedEvent, m.view.ProcessedStylingFailedEvent, m.view.ProcessedStylingSequence, m.errorCountUntilSkip)
}

func (m *Styling) OnSuccess() error {
	return spooler.HandleSuccess(m.view.UpdateStylingSpoolerRunTimestamp)
}

func (m *Styling) generateStylingFile(policy *iam_model.LabelPolicyView) error {
	reader, size, err := m.writeFile(policy)
	if err != nil {
		return err
	}
	return m.uploadFilesToBucket(policy.AggregateID, "text/css", reader, size)
}

func (m *Styling) writeFile(policy *iam_model.LabelPolicyView) (io.Reader, int64, error) {
	cssContent := ""
	cssContent += ":root {"
	if policy.PrimaryColor != "" {
		palette := m.generateColorPaletteRGBA255(policy.PrimaryColor)
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-primary-%v: %s;", i, color)
		}
	}

	if policy.BackgroundColor != "" {
		palette := m.generateColorPaletteRGBA255(policy.BackgroundColor)
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-background-%v: %s;", i, color)
		}
	}
	if policy.WarnColor != "" {
		palette := m.generateColorPaletteRGBA255(policy.WarnColor)
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-warn-%v: %s;", i, color)
		}
	}
	if policy.FontColor != "" {
		palette := m.generateColorPaletteRGBA255(policy.FontColor)
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-text-%v: %s;", i, color)
		}
	}
	var fontname string
	if policy.FontURL != "" {
		split := strings.Split(policy.FontURL, "/")
		fontname = split[len(split)-1]
		cssContent += fmt.Sprintf("--zitadel-font-family: %s;", fontname)
	}
	cssContent += "}"
	if policy.FontURL != "" {
		cssContent += fmt.Sprintf(fontFaceTemplate, fontname, m.resourceUrl, policy.AggregateID, policy.FontURL)
	}
	cssContent += ".lgn-dark-theme {"
	if policy.PrimaryColorDark != "" {
		palette := m.generateColorPaletteRGBA255(policy.PrimaryColorDark)
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-primary-%v: %s;", i, color)
		}
	}
	if policy.BackgroundColorDark != "" {
		palette := m.generateColorPaletteRGBA255(policy.BackgroundColorDark)
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-background-%v: %s;", i, color)
		}
	}
	if policy.WarnColorDark != "" {
		palette := m.generateColorPaletteRGBA255(policy.WarnColorDark)
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-warn-%v: %s;", i, color)
		}
	}
	if policy.FontColorDark != "" {
		palette := m.generateColorPaletteRGBA255(policy.FontColorDark)
		for i, color := range palette {
			cssContent += fmt.Sprintf("--zitadel-color-text-%v: %s;", i, color)
		}
	}
	cssContent += "}"

	data := []byte(cssContent)
	buffer := bytes.NewBuffer(data)
	return buffer, int64(buffer.Len()), nil
}

const fontFaceTemplate = `
@font-face {
	font-family: '%s';
	font-style: normal;
	font-display: swap;
	src: url(%s?orgId=%s&filename=%s);
}
`

func (m *Styling) uploadFilesToBucket(aggregateID, contentType string, reader io.Reader, size int64) error {
	fileName := domain.CssPath + "/" + domain.CssVariablesFileName
	_, err := m.static.PutObject(context.Background(), aggregateID, fileName, contentType, reader, size, true)
	return err
}

func (m *Styling) generateColorPaletteRGBA255(hex string) map[string]string {
	palette := make(map[string]string)
	defaultColor := gamut.Hex(hex)

	color50, ok := colorful.MakeColor(gamut.Lighter(defaultColor, 1.0))
	if ok {
		palette["50"] = cssRGB(color50.RGB255())
	}

	color100, ok := colorful.MakeColor(gamut.Lighter(defaultColor, 0.8))
	if ok {
		palette["100"] = cssRGB(color100.RGB255())
	}

	color200, ok := colorful.MakeColor(gamut.Lighter(defaultColor, 0.6))
	if ok {
		palette["200"] = cssRGB(color200.RGB255())
	}

	color300, ok := colorful.MakeColor(gamut.Lighter(defaultColor, 0.4))
	if ok {
		palette["300"] = cssRGB(color300.RGB255())
	}

	color400, ok := colorful.MakeColor(gamut.Lighter(defaultColor, 0.1))
	if ok {
		palette["400"] = cssRGB(color400.RGB255())
	}

	color500, ok := colorful.MakeColor(defaultColor)
	if ok {
		palette["500"] = cssRGB(color500.RGB255())
	}

	color600, ok := colorful.MakeColor(gamut.Darker(defaultColor, 0.1))
	if ok {
		palette["600"] = cssRGB(color600.RGB255())
	}

	color700, ok := colorful.MakeColor(gamut.Darker(defaultColor, 0.2))
	if ok {
		palette["700"] = cssRGB(color700.RGB255())
	}

	color800, ok := colorful.MakeColor(gamut.Darker(defaultColor, 0.3))
	if ok {
		palette["800"] = cssRGB(color800.RGB255())
	}

	color900, ok := colorful.MakeColor(gamut.Darker(defaultColor, 0.4))
	if ok {
		palette["900"] = cssRGB(color900.RGB255())
	}

	colorContrast, ok := colorful.MakeColor(gamut.Contrast(defaultColor))
	if ok {
		palette["contrast"] = cssRGB(colorContrast.RGB255())
	}

	return palette
}

func cssRGB(r, g, b uint8) string {
	return fmt.Sprintf("rgb(%v, %v, %v)", r, g, b)
}
