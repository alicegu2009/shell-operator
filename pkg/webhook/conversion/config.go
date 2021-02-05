package conversion

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
