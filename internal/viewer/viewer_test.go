package viewer

import (
	"encoding/json"
	"testing"
)

func TestOpenPayloadJSON(t *testing.T) {
	payload := OpenPayload{
		Images: []string{"a.png", "b.png"},
		Start:  1,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded OpenPayload
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if len(decoded.Images) != 2 || decoded.Images[1] != "b.png" || decoded.Start != 1 {
		t.Fatalf("decoded payload mismatch: %+v", decoded)
	}
}

func TestNewWailsViewerReal(t *testing.T) {
	v := NewWailsViewerReal()
	if v == nil {
		t.Fatal("expected non-nil viewer client")
	}
}
