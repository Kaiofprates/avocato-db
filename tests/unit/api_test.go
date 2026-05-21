package unit

import (
	"avocato-db/src/api/handlers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBlockHandler_Security(t *testing.T) {
	// We don't need a real service because the handler should reject the ID 
	// before calling the service.
	handler := handlers.NewGetBlockHandler(nil)

	tests := []struct {
		name       string
		id         string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Reject Index",
			id:         "1",
			wantStatus: http.StatusBadRequest,
			wantBody:   "Invalid UUID format",
		},
		{
			name:       "Reject Hash",
			id:         "a3e75136c37a20f48c05ea89db8047889a5520bbda7f5c461661141a99a60e1f",
			wantStatus: http.StatusBadRequest,
			wantBody:   "Invalid UUID format",
		},
		{
			name:       "Reject Junk",
			id:         "not-even-close",
			wantStatus: http.StatusBadRequest,
			wantBody:   "Invalid UUID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/v1/block/"+tt.id, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.wantBody)
		})
	}
}

func TestAppendResponse_Fields(t *testing.T) {
	// Test if AppendResponse struct has all required fields for JSON 
	// to prevent regressions in field naming or missing fields.
	resp := handlers.AppendResponse{}
	resp.Data.UUID = "test-uuid"
	resp.Data.Index = 1

	data, err := json.Marshal(resp)
	assert.NoError(t, err)

	var decoded map[string]interface{}
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	innerData := decoded["data"].(map[string]interface{})
	assert.NotNil(t, innerData["uuid"], "uuid field is missing from JSON response")
	assert.Equal(t, "test-uuid", innerData["uuid"])
}
