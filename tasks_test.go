package mdmstorage

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestExtractedParamsTruckTypeRoundTrip(t *testing.T) {
	want := "b_double|semi_trailer"
	original := &ExtractedParams{TruckType: &want}

	payload, err := original.ToJson()
	if err != nil {
		t.Fatalf("marshal ExtractedParams: %v", err)
	}
	if !strings.Contains(payload, `"truck_type":"b_double|semi_trailer"`) {
		t.Fatalf("payload %q does not contain truck_type", payload)
	}

	decoded, err := ToExtractedParams(payload)
	if err != nil {
		t.Fatalf("unmarshal ExtractedParams: %v", err)
	}
	if decoded.TruckType == nil || *decoded.TruckType != want {
		t.Fatalf("TruckType = %v, want %q", decoded.TruckType, want)
	}
}

func TestExtractedParamsTruckTypeCompatibility(t *testing.T) {
	legacy, err := ToExtractedParams(`{"key":"legacy","mode":"truck"}`)
	if err != nil {
		t.Fatalf("unmarshal legacy ExtractedParams: %v", err)
	}
	if legacy.TruckType != nil {
		t.Fatalf("legacy TruckType = %v, want nil", legacy.TruckType)
	}

	payload, err := json.Marshal(ExtractedParams{})
	if err != nil {
		t.Fatalf("marshal empty ExtractedParams: %v", err)
	}
	if strings.Contains(string(payload), "truck_type") {
		t.Fatalf("nil TruckType must be omitted: %s", payload)
	}
}

func TestExtractedParamsTruckTypeTags(t *testing.T) {
	field, ok := reflect.TypeOf(ExtractedParams{}).FieldByName("TruckType")
	if !ok {
		t.Fatal("TruckType field not found")
	}
	if got := field.Tag.Get("json"); got != "truck_type,omitempty" {
		t.Fatalf("json tag = %q, want truck_type,omitempty", got)
	}
	if got := field.Tag.Get("form"); got != "truck_type,omitempty" {
		t.Fatalf("form tag = %q, want truck_type,omitempty", got)
	}
}
