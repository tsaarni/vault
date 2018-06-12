package query

import (
	"github.com/hashicorp/vault/builtin/credential/alibaba/alibaba-sdk-go/aws/request"
)

const format = "JSON"

func NewBuildHandler(accessKeyID string) request.NamedHandler {
	b := &buildHandler{accessKeyID: accessKeyID}
	return request.NamedHandler{Name: "awssdk.query.Build", Fn: b.Build}
}

type buildHandler struct {
	accessKeyID string
}

// Build builds a request for an AWS Query service.
func (b *buildHandler) Build(r *request.Request) {
	if r.HTTPRequest.URL.RawQuery != "" {
		r.HTTPRequest.URL.RawQuery += "&"
	}
	r.HTTPRequest.URL.RawQuery += "Format=" + format
	r.HTTPRequest.URL.RawQuery += "&"
	r.HTTPRequest.URL.RawQuery += "AccessKeyId=" + b.accessKeyID
	// TODO do I need to add content-type or accept headers?
}
