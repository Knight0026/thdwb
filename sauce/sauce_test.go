package sauce

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	structs "thdwb/structs"

	"github.com/stretchr/testify/assert"
)

func TestGetResource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`OK!`))
	}))

	defer server.Close()

	url, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	resource := GetResource(url, nil)
	var testResorce *structs.Resource

	assert.IsType(t, testResorce, resource, "Expecting: Resource Struct")
	assert.Equal(t, "OK!", resource.Body, "Expecting: OK!")
}
