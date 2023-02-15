package access

import (
	"reflect"
	"testing"
)

func TestRecord_Normalize(t *testing.T) {
	tests := []struct {
		name   string
		record Record
		want   *Record
	}{{
		name: "headers with certain keys should be redacted",
		record: Record{
			RequestHeaders: map[string][]string{
				"authorization":             {"AValue"},
				"grpcgateway-authorization": {"AValue"},
				"cookie":                    {"AValue"},
				"grpcgateway-cookie":        {"AValue"},
			}, ResponseHeaders: map[string][]string{
				"set-cookie": {"AValue"},
			},
		},
		want: &Record{
			RequestHeaders: map[string][]string{
				"authorization":             {"[REDACTED]"},
				"grpcgateway-authorization": {"[REDACTED]"},
				"cookie":                    {"[REDACTED]"},
				"grpcgateway-cookie":        {"[REDACTED]"},
			}, ResponseHeaders: map[string][]string{
				"set-cookie": {"[REDACTED]"},
			},
		},
	}, {
		name: "header keys should be lower cased",
		record: Record{
			RequestHeaders:  map[string][]string{"AKey": {"AValue"}},
			ResponseHeaders: map[string][]string{"AKey": {"AValue"}}},
		want: &Record{
			RequestHeaders:  map[string][]string{"akey": {"AValue"}},
			ResponseHeaders: map[string][]string{"akey": {"AValue"}}},
	}, {
		name: "an already prune record should stay unchanged",
		record: Record{
			RequestURL: "https://my.zitadel.cloud/",
			RequestHeaders: map[string][]string{
				"authorization": {"[REDACTED]"},
			},
			ResponseHeaders: map[string][]string{},
		},
		want: &Record{
			RequestURL: "https://my.zitadel.cloud/",
			RequestHeaders: map[string][]string{
				"authorization": {"[REDACTED]"},
			},
			ResponseHeaders: map[string][]string{},
		},
	}, {
		name: "empty record should stay empty",
		record: Record{
			RequestHeaders:  map[string][]string{},
			ResponseHeaders: map[string][]string{},
		},
		want: &Record{
			RequestHeaders:  map[string][]string{},
			ResponseHeaders: map[string][]string{},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.record.Normalize(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}
