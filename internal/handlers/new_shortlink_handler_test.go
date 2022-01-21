package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"range/internal/storage"
	"testing"
	"time"
)

type mockConfigger struct{}

func (m mockConfigger) GetHost() string {
	return "host.com"
}

type mockShortLinkWriter struct {
	mock.Mock
}

func (m *mockShortLinkWriter) Save(request storage.ShortLinkRequest) (storage.ShortLinkResponse, error) {
	args := m.Called(request)
	return args.Get(0).(storage.ShortLinkResponse), args.Error(1)
}

func TestNewShortLinkHandler(t *testing.T) {
	tests := []struct {
		name     string
		req      ShortLinkRequest
		wantCode int
		sto      func() shortLinkWriter
		cfg      configger
	}{
		{
			name:     "nothingburger_error_is_bad_request",
			wantCode: http.StatusBadRequest,
			cfg:      mockConfigger{},
			sto: func() shortLinkWriter {
				return nil
			},
		},
		{
			name:     "conflict_from_storage_is_http_conflict",
			req:      ShortLinkRequest{Link: "http://example.com"},
			wantCode: http.StatusConflict,
			cfg:      mockConfigger{},
			sto: func() shortLinkWriter {
				w := &mockShortLinkWriter{}
				w.On("Save", mock.AnythingOfType("ShortLinkRequest")).Return(storage.ShortLinkResponse{}, storage.ErrConflict{})
				return w
			},
		},
		{
			name:     "other_error_from_storage_is_server_error",
			req:      ShortLinkRequest{Link: "http://example.com"},
			wantCode: http.StatusInternalServerError,
			cfg:      mockConfigger{},
			sto: func() shortLinkWriter {
				w := &mockShortLinkWriter{}
				w.On("Save", mock.AnythingOfType("ShortLinkRequest")).Return(storage.ShortLinkResponse{}, errors.New("whoopsie daisy"))
				return w
			},
		},
		{
			name:     "happy_path",
			req:      ShortLinkRequest{Link: "http://example.com"},
			wantCode: http.StatusOK,
			cfg:      mockConfigger{},
			sto: func() shortLinkWriter {
				w := &mockShortLinkWriter{}
				hash := "clever-hash"
				w.On("Save", mock.AnythingOfType("ShortLinkRequest")).Return(storage.ShortLinkResponse{
					Link:      "",
					Hash:      &hash,
					Suffix:    nil,
					CreatedAt: time.Now(),
				}, nil)
				return w
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewShortLinkHandler(tt.sto(), tt.cfg)
			w := httptest.NewRecorder()
			j, _ := json.Marshal(tt.req)
			r := httptest.NewRequest(http.MethodPost, "/target", bytes.NewBuffer(j))
			h(w, r)
			if w.Code != tt.wantCode {
				t.Errorf("got %v, want %v", w.Code, tt.wantCode)
			}
		})
	}
}
