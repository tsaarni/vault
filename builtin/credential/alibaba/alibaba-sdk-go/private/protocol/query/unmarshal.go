package query

//go:generate go run -tags codegen ../../../models/protocol_tests/generate.go ../../../models/protocol_tests/output/query.json unmarshal_test.go

import (
	"encoding/json"

	"github.com/hashicorp/vault/builtin/credential/alibaba/alibaba-sdk-go/aws/awserr"
	"github.com/hashicorp/vault/builtin/credential/alibaba/alibaba-sdk-go/aws/request"
)

// UnmarshalHandler is a named request handler for unmarshaling query protocol requests
var UnmarshalHandler = request.NamedHandler{Name: "awssdk.query.Unmarshal", Fn: Unmarshal}

// UnmarshalMetaHandler is a named request handler for unmarshaling query protocol request metadata
var UnmarshalMetaHandler = request.NamedHandler{Name: "awssdk.query.UnmarshalMeta", Fn: UnmarshalMeta}

// Unmarshal unmarshals a response for an AWS Query service.
func Unmarshal(r *request.Request) {
	defer r.HTTPResponse.Body.Close()
	if r.DataFilled() {
		if err := json.NewDecoder(r.HTTPResponse.Body).Decode(r.Data); err != nil {
			r.Error = awserr.New("SerializationError", "failed decoding Query response", err)
			return
		}
	}
}

// UnmarshalMeta unmarshals header response values for an AWS Query service.
func UnmarshalMeta(r *request.Request) {
	r.RequestID = r.HTTPResponse.Header.Get("X-Amzn-Requestid")
}
