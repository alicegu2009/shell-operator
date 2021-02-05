package types

import (
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ConversionEvent struct {
	CrdName string
	Review  *v1.ConversionReview
	Objects []unstructured.Unstructured
}

// Mimic a v1.ConversionReview structure but with the array of unstructured Objects
// instead of the array of runtime.RawExtension
func (c ConversionEvent) GetReview() map[string]interface{} {
	return map[string]interface{}{
		"kind":       c.Review.Kind,
		"apiVersion": c.Review.APIVersion,
		"request": map[string]interface{}{
			"uid":               c.Review.Request.UID,
			"desiredAPIVersion": c.Review.Request.DesiredAPIVersion,
			"objects":           c.Objects,
		},
	}
}
