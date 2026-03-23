package alfabank

import (
	"encoding/json"
	"testing"
)

func TestFlexibleStringUnmarshal(t *testing.T) {
	t.Run("from string", func(t *testing.T) {
		var value FlexibleString
		if err := json.Unmarshal([]byte(`"810"`), &value); err != nil {
			t.Fatalf("unmarshal failed: %v", err)
		}
		if value != "810" {
			t.Fatalf("unexpected value: %q", value)
		}
	})

	t.Run("from number", func(t *testing.T) {
		var value FlexibleString
		if err := json.Unmarshal([]byte(`810`), &value); err != nil {
			t.Fatalf("unmarshal failed: %v", err)
		}
		if value != "810" {
			t.Fatalf("unexpected value: %q", value)
		}
	})
}

func TestFlexibleInt64Unmarshal(t *testing.T) {
	t.Run("from string", func(t *testing.T) {
		var value FlexibleInt64
		if err := json.Unmarshal([]byte(`"123"`), &value); err != nil {
			t.Fatalf("unmarshal failed: %v", err)
		}
		if value != 123 {
			t.Fatalf("unexpected value: %d", value)
		}
	})

	t.Run("from number", func(t *testing.T) {
		var value FlexibleInt64
		if err := json.Unmarshal([]byte(`123`), &value); err != nil {
			t.Fatalf("unmarshal failed: %v", err)
		}
		if value != 123 {
			t.Fatalf("unexpected value: %d", value)
		}
	})

	t.Run("from null", func(t *testing.T) {
		var value FlexibleInt64 = 99
		if err := json.Unmarshal([]byte(`null`), &value); err != nil {
			t.Fatalf("unmarshal failed: %v", err)
		}
		if value != 0 {
			t.Fatalf("unexpected value: %d", value)
		}
	})
}

func TestRefundCollectionUnmarshal(t *testing.T) {
	t.Run("single object", func(t *testing.T) {
		var value RefundCollection
		if err := json.Unmarshal([]byte(`{"externalRefundId":"rf-1","actionCode":0}`), &value); err != nil {
			t.Fatalf("unmarshal failed: %v", err)
		}
		if len(value) != 1 || value[0].ExternalRefundID != "rf-1" || value[0].ActionCode != "0" {
			t.Fatalf("unexpected value: %#v", value)
		}
	})

	t.Run("array", func(t *testing.T) {
		var value RefundCollection
		if err := json.Unmarshal([]byte(`[ {"externalRefundId":"rf-1"}, {"externalRefundId":"rf-2"} ]`), &value); err != nil {
			t.Fatalf("unmarshal failed: %v", err)
		}
		if len(value) != 2 || value[1].ExternalRefundID != "rf-2" {
			t.Fatalf("unexpected value: %#v", value)
		}
	})
}

func TestEncodeForm(t *testing.T) {
	type payload struct {
		Amount      int64             `form:"amount"`
		Features    []string          `form:"features,omitempty"`
		JSONParams  map[string]string `form:"jsonParams,omitempty,json"`
		Description string            `form:"description,omitempty"`
	}

	values, err := encodeForm(payload{
		Amount:     1200,
		Features:   []string{"VERIFY", "FORCE_TDS"},
		JSONParams: map[string]string{"k": "v"},
	})
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	if values.Get("amount") != "1200" {
		t.Fatalf("unexpected amount: %q", values.Get("amount"))
	}
	if got := values["features"]; len(got) != 2 || got[0] != "VERIFY" || got[1] != "FORCE_TDS" {
		t.Fatalf("unexpected features: %#v", got)
	}
	if values.Get("jsonParams") != `{"k":"v"}` {
		t.Fatalf("unexpected jsonParams: %q", values.Get("jsonParams"))
	}
	if _, ok := values["description"]; ok {
		t.Fatalf("description should be omitted: %#v", values)
	}
}