package vesapi

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type request struct {
	Method  string
	Headers map[string]string
	Path    string
	Body    string
}

type response struct {
	code int
	body string
}

func TestGetVehicleStatus(t *testing.T) {

	testCases := []struct {
		name     string
		request  request
		response *response
		vehicle  *VehicleStatus
		err      error
	}{
		{
			name: "invalid url returns error",
			request: request{
				Path: "/vehicle-enquiry/v1/vehicles",
			},
			response: &response{
				code: 404,
				body: `not_found`,
			},
			err:     errors.New("HTTP 404: not_found"),
			vehicle: &VehicleStatus{},
		},
		{
			name: "returns vehicle status",
			request: request{
				Path: "/vehicle-enquiry/v1/vehicles",
			},
			response: &response{
				code: 200,
				body: `{"colour":"red","taxStatus":"expired","markedForExport":true,"dateOfLastV5CIssued":"2020-01-01"}`,
			},
			vehicle: &VehicleStatus{
				Colour:          "red",
				TaxStatus:       "expired",
				MarkedForExport: true,
				DateOfLastV5CIssued: Date{
					Time: time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, tc := range testCases {
		reqCount := 0
		expAPIKey := "12435"

		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				reqCount++

				apiKey := r.Header.Get("X-Api-Key")
				if apiKey != expAPIKey {
					t.Errorf("Expected api key header '%s' but got '%s'", expAPIKey, apiKey)
				}

				ct := r.Header.Get("Content-Type")
				if ct != "application/json" {
					t.Errorf("Expected content type 'application/json' but got '%s'", ct)
				}

				if r.Method != http.MethodPost {
					t.Errorf("Expected method '%s' but got '%s'", http.MethodPost, r.Method)
				}

				if tc.response != nil {
					w.WriteHeader(tc.response.code)
					fmt.Fprintf(w, tc.response.body)
					return
				}
			}))

			c := NewClient("12435", server.URL)
			v, err := c.GetVehicleStatus(tc.name)

			if fmt.Sprint(err) != fmt.Sprint(tc.err) {
				t.Errorf("Expected errors to match: got '%s' want: '%s'", err, tc.err)
			}

			if !reflect.DeepEqual(v, tc.vehicle) {
				t.Error("Expected vehicle status to match but didn't", v, tc.vehicle)
			}

			server.Close()
		})
	}
}
