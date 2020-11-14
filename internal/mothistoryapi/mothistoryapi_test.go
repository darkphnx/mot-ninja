package mothistoryapi

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

func TestGetVehicleHistory(t *testing.T) {

	testCases := []struct {
		name     string
		request  request
		response *response
		vehicle  *Vehicle
		err      error
	}{
		{
			name: "invalid url returns error",
			request: request{
				Path: "/vehicle-enquiry/v1/vehicles?registration=P239FWP",
			},
			response: &response{
				code: 404,
				body: `not_found`,
			},
			err:     errors.New("HTTP 404: not_found"),
			vehicle: nil,
		},
		{
			name: "returns vehicle mot history",
			request: request{
				Path: "/vehicle-enquiry/v1/vehicles?registration=P239FWP",
			},
			response: &response{
				code: 200,
				body: `[{
					"registration":"P239FWP",
					"make":"MAZDA",
					"model":"MPV",
					"firstUsedDate":"1996.12.31",
					"fuelType":"Diesel",
					"primaryColour": "White",
					"vehicleId":"n_wLOetTguVjsHCoUEhspw==",
					"registrationDate":"1996.08.01",
					"manufactureDate":"1996.12.31",
					"engineSize":"1998",
					"motTests":[{
						"completedDate":"2020.10.21 08:17:47",
						"testResult":"PASSED",
						"expiryDate":"2021.10.20",
						"odometerValue":"200413",
						"odometerUnit":"mi",
						"motTestNumber":"901662956826",
						"odometerResultType":"READ",
						"rfrAndComments":[{
							"text":"Nearside Front Track rod end ball joint dust cover damaged or  deteriorated, but preventing the ingress of dirt (2.1.3 (g) (i))",
							"type":"MINOR",
							"dangerous":false
						}]
					}]
				}]`,
			},
			vehicle: &Vehicle{
				Registration: "P239FWP",
				Make:         "MAZDA",
				Model:        "MPV",
				FirstUsedDate: DottedDate{
					Time: time.Date(1996, 12, 31, 0, 0, 0, 0, time.UTC),
				},
				FuelType:      "Diesel",
				PrimaryColour: "White",
				VehicleID:     "n_wLOetTguVjsHCoUEhspw==",
				RegistrationDate: DottedDate{
					Time: time.Date(1996, 8, 1, 0, 0, 0, 0, time.UTC),
				},
				ManufactureDate: DottedDate{
					Time: time.Date(1996, 12, 31, 0, 0, 0, 0, time.UTC),
				},
				EngineSize: 1998,
				MotTests: []MotTest{
					MotTest{
						CompletedDate: DottedTime{
							Time: time.Date(2020, 10, 21, 8, 17, 47, 0, time.UTC),
						},
						TestResult: "PASSED",
						ExpiryDate: DottedDate{
							Time: time.Date(2021, 10, 20, 0, 0, 0, 0, time.UTC),
						},
						OdometerValue:      200413,
						OdometerUnit:       "mi",
						MotTestNumber:      "901662956826",
						OdometerResultType: "READ",
						RfrAndComments: []RfrAndComment{
							RfrAndComment{
								Text:      "Nearside Front Track rod end ball joint dust cover damaged or  deteriorated, but preventing the ingress of dirt (2.1.3 (g) (i))",
								Type:      "MINOR",
								Dangerous: false,
							},
						},
					},
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

				ct := r.Header.Get("Accept")
				if ct != "application/json+v6" {
					t.Errorf("Expected accept 'application/json' but got '%s'", ct)
				}

				if r.Method != http.MethodGet {
					t.Errorf("Expected method '%s' but got '%s'", http.MethodGet, r.Method)
				}

				if tc.response != nil {
					w.WriteHeader(tc.response.code)
					fmt.Fprintf(w, tc.response.body)
					return
				}
			}))

			c := NewClient("12435", server.URL)
			v, err := c.GetVehicleHistory(tc.name)

			if fmt.Sprint(err) != fmt.Sprint(tc.err) {
				t.Errorf("Expected errors to match: got '%s' want: '%s'", err, tc.err)
			}

			if !reflect.DeepEqual(v, tc.vehicle) {
				t.Error("Expected vehicle to match but didn't", v, tc.vehicle)
			}

			server.Close()
		})
	}
}
