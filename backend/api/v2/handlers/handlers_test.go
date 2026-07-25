package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateEngineHoursRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantErr bool
	}{
		{
			name:    "valid payload with null engine_hours_after",
			payload: `{"user_id":1,"engine_hours_before":400,"engine_hours_after":null,"type":1}`,
			wantErr: false,
		},
		{
			name:    "valid payload with null engine_hours_after as string",
			payload: `{"user_id":1,"engine_hours_before":"400","engine_hours_after":null,"type":1}`,
			wantErr: false,
		},
		{
			name:    "valid payload with non-null engine_hours_after",
			payload: `{"user_id":1,"engine_hours_before":400,"engine_hours_after":405.5,"type":1}`,
			wantErr: false,
		},
		{
			name:    "missing required user_id",
			payload: `{"engine_hours_before":400,"type":1}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/api/v2/boat/engine-hour/update", bytes.NewBufferString(tt.payload))
			req := &UpdateEngineHoursRequest{}
			err := ReadBodyAndValidate(r, req, UpdateEngineHoursValidationErrors)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadBodyAndValidate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
