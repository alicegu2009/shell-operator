package conversion

import (
	"strings"

	"github.com/flant/shell-operator/pkg/utils/string_helper"
)

// CrdName also used as a first element in Path field for spec.conversion.clientConfig in CRD.
// It should be already an url safe.

// ConversionWebhookConfig
type ConversionWebhookConfig struct {
	Conversions []ConversionRule
	CrdName     string // The name is used as a suffix to create different URLs for clientConfig in CRD.
	Metadata    struct {
		Name         string
		DebugName    string
		LogLabels    map[string]string
		MetricLabels map[string]string
	}
}

type ConversionRule struct {
	FromVersion string `json:"fromVersion"`
	ToVersion   string `json:"toVersion"`
}

func (r ConversionRule) String() string {
	return r.FromVersion + "->" + r.ToVersion
}

func (r ConversionRule) ShortFromVersion() string {
	return string_helper.TrimGroup(r.FromVersion)
}

func (r ConversionRule) ShortToVersion() string {
	return string_helper.TrimGroup(r.ToVersion)
}

func RuleFromString(in string) ConversionRule {
	idx := strings.Index(in, "->")
	return ConversionRule{
		FromVersion: in[0:idx],
		ToVersion:   in[idx+2:],
	}
}
