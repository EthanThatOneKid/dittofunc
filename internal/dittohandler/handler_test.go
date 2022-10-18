package dittoclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseConfigFromRequest(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *ProgramConfig
		want1   *GenerateConfig
		wantErr bool
	}{
		name: "invalid request method",
		args: args{
			r: httptest.NewRequest(http.MethodGet, "/", nil),
		},
		wantErr: true,
	},
	{
		name: "invalid content type",
		args: args{
			r: httptest.NewRequest(http.MethodPost, "/", nil),
		},
		wantErr: true,
	},
	{ 
		name: "invalid request body",
		args: args{
			r: httptest.NewRequest(http.MethodPost, "/", nil),
		},
		wantErr: true,
	},
	{
		name: "invalid program config",
		args: args{
			r: httptest.NewRequest(http.MethodPost, "/", nil),
		},
		wantErr: true,
	},
	{
		name: "invalid generate config",
		args: args{
			r: httptest.NewRequest(http.MethodPost, "/", nil),
		},
		wantErr: true,
	},
	{
		name: "valid request",
		args: args{
			r: httptest.NewRequest(http.MethodPost, "/", nil),
		},
		want:    &ProgramConfig{},
		want1:   &GenerateConfig{},
		wantErr: false,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseConfigFromRequest(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConfigFromRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseConfigFromRequest() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ParseConfigFromRequest() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}