package kibana

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	cfg := Config{
		URL:      "http://127.0.0.1:5601",
		Username: "elastic",
		Password: "changeme",
		Insecure: true,
	}
	client, err := NewClient(cfg)
	if err != nil {
		t.Error(err)
	}
	if client == nil {
		t.Error("client was nil")
	}
}
