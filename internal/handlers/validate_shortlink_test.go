package handlers

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShortLinkRequest_Validate(t *testing.T) {
	tests := []struct {
		name string
		req  ShortLinkRequest
		cfg  configger
		want error
	}{
		{
			name: "nothingburger_error",
			want: errors.New("parse \"\": empty url"),
		},
		{
			name: "host_matching_is_error",
			cfg:  mockConfigger{},
			req:  ShortLinkRequest{Link: "http://" + mockConfigger{}.GetHost() + "/new"},
			want: errors.New("can't do that"),
		},
		{
			name: "errors_from_parse_get_forwarded",
			cfg:  mockConfigger{},
			req:  ShortLinkRequest{Link: "not a link"},
			want: errors.New("parse \"not a link\": invalid URI for request"),
		},
		{
			name: "happy_path",
			cfg:  mockConfigger{},
			req:  ShortLinkRequest{Link: "http://example.com"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.req.Validate(tt.cfg)
			if got == nil {
				assert.NoError(t, tt.want)
			} else {
				assert.EqualError(t, tt.want, got.Error())
			}
		})
	}
}
