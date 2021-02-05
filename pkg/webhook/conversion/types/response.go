package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

/*

FailedMessage:
Response is Success if empty string:
"result": {
  "status": "Success"
},
Response is Failed:
"result": {
  "status": "Failed",
  "message": FailedMessage
}

ConvertedObjects:
# Objects must match the order of request.objects, and have apiVersion set to <request.desiredAPIVersion>.
# kind, metadata.uid, metadata.name, and metadata.namespace fields must not be changed by the webhook.
# metadata.labels and metadata.annotations fields may be changed by the webhook.
# All other changes to metadata fields by the webhook are ignored.
*/
type ConversionResponse struct {
	FailedMessage    string                      `json:"failedMessage"`
	ConvertedObjects []unstructured.Unstructured `json:"convertedObjects,omitempty"`
}

func ConversionResponseFromFile(filePath string) (*ConversionResponse, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %s", filePath, err)
	}

	if len(data) == 0 {
		return nil, nil
	}
	return ConversionResponseFromBytes(data)
}

func ConversionResponseFromBytes(data []byte) (*ConversionResponse, error) {
	return ConversionResponseFromReader(bytes.NewReader(data))
}

func ConversionResponseFromReader(r io.Reader) (*ConversionResponse, error) {
	response := new(ConversionResponse)

	dec := json.NewDecoder(r)

	err := dec.Decode(response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *ConversionResponse) Dump() string {
	b := new(strings.Builder)
	b.WriteString("ConversionResponse(failedMessage=")
	b.WriteString(r.FailedMessage)
	if len(r.ConvertedObjects) > 0 {
		b.WriteString(",convertedObjects.len=")
		b.WriteString(strconv.FormatInt(int64(len(r.ConvertedObjects)), 10))
	}
	b.WriteString(")")
	return b.String()
}
