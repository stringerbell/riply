package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewError(t *testing.T) {
	tests := []struct {
		name  string
		error error
		want  err
	}{
		{
			name:  "nil error shows internal detail",
			error: nil,
			want:  err{Error: true, Detail: "internal"},
		},
		{
			name:  "errors get forwarded",
			error: errors.New("some error"),
			want:  err{Error: true, Detail: "some error"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			got := newError(w, tt.error, http.StatusInternalServerError)
			want := err{}
			json.Unmarshal(got, &want)
			if !reflect.DeepEqual(want, tt.want) {
				t.Errorf("mismatched err got: %v, want: %v", want, tt.want)
			}
		})
	}
}
