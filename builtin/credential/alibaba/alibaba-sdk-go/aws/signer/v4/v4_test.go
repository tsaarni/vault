package v4

import (
	"net/http"
	"testing"

	"time"

	"net/url"

	"github.com/hashicorp/vault/builtin/credential/alibaba/alibaba-sdk-go/aws/request"
)

const testURL = "http://ecs.aliyuncs.com/?Timestamp=2016-02-23T12:46:24Z&Format=XML&AccessKeyId=testid&Action=DescribeRegions&SignatureMethod=HMAC-SHA1&SignatureNonce=3ee8c1b8-83d3-44af-a94f-4e0ad82fd6cf&Version=2014-05-26&SignatureVersion=1.0"

var (
	curTimeFunc = func() time.Time {
		return time.Date(2016, time.February, 23, 12, 46, 24, 0, time.UTC)
	}
)

func TestStringToSign(t *testing.T) {
	expected := "GET&%2F&AccessKeyId%3Dtestid&Action%3DDescribeRegions&Format%3DXML&SignatureMethod%3DHMAC-SHA1&SignatureNonce%3D3ee8c1b8-83d3-44af-a94f-4e0ad82fd6cf&SignatureVersion%3D1.0&Timestamp%3D2016-02-23T12%253A46%253A24Z&Version%3D2014-05-26"
	result, err := stringToSign("GET", testURL)
	if err != nil {
		t.Fatal(err)
	}
	if result != expected {
		t.Fatalf("expected %s but received %s", expected, result)
	}
}

func TestSignature(t *testing.T) {
	expected := "VyBL52idtt+oImX0NZC+2ngk15Q="
	result, _ := signature("GET", testURL, "testsecret")
	if result != expected {
		t.Fatalf("expected %s but received %s", expected, result)
	}
}

func TestSignedURL(t *testing.T) {
	u, err := url.Parse("http://ecs.aliyuncs.com/?Format=XML&AccessKeyId=testid&Action=DescribeRegions&SignatureNonce=3ee8c1b8-83d3-44af-a94f-4e0ad82fd6cf&Version=2014-05-26")
	if err != nil {
		t.Fatal(err)
	}

	httpReq := &http.Request{
		URL:    u,
		Method: "GET",
	}

	s := &signRequestHandler{accessKeySecret: "testsecret"}

	r := &request.Request{
		HTTPRequest: httpReq,
	}
	if err := s.signSDKRequestWithCurrTime(r, curTimeFunc); err != nil {
		t.Fatal(err)
	}

	expected := "http://ecs.aliyuncs.com/?SignatureVersion=1.0&Action=DescribeRegions&Format=XML&SignatureNonce=3ee8c1b8-83d3-44af-a94f-4e0ad82fd6cf&Version=2014-05-26&AccessKeyId=testid&Signature=VyBL52idtt+oImX0NZC+2ngk15Q%3D&SignatureMethod=HMAC-SHA1&Timestamp=2016-02-23T12%3A46%3A24Z"
	expectedAsUrl, err := url.Parse(expected)
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range expectedAsUrl.Query() {
		actual := r.HTTPRequest.URL.Query().Get(k)
		if actual != v[0] {
			t.Fatalf("expected %s for %s, but received %s", v[0], k, actual)
		}
	}
}
