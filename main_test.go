package main

import (
	"bytes"
	"encoding/json"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// func TestCreatePostAccountResponse(t *testing.T) {
// 	type args struct {
// 		account Account
// 	}
// 	tests := []struct {
// 		name  string
// 		args  args
// 		wantR PostAccountResponse
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if gotR := CreatePostAccountResponse(tt.args.account); !reflect.DeepEqual(gotR, tt.wantR) {
// 				t.Errorf("CreatePostAccountResponse() = %v, want %v", gotR, tt.wantR)
// 			}
// 		})
// 	}
// }

func TestGetDayOfTheWeek(t *testing.T) {
	tests := []struct {
		name        string
		DateOfBirth string
		WeekDay     string
	}{
		{"October 4 2022", "2022-10-04", "Tuesday"},
		{"No Date", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dayOfTheWeek := GetDayOfTheWeek(tt.DateOfBirth)
			if dayOfTheWeek != tt.WeekDay {
				t.Errorf("expected %v got %v", tt.WeekDay, dayOfTheWeek)
			}
		})
	}
}

func Test_getRFC3339EST(t *testing.T) {
	tests := []struct {
		name       string
		CreatedOn  int
		RFC3339EST string
	}{
		{"Oct 5 2022 12:50:33 EST", 1664988633, "2022-10-05T12:50:33-05:00"},
		{"Oct 5 2022 7:04:30 EST", 1665011070, "2022-10-05T19:04:30-05:00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getRFC3339EST(tt.CreatedOn)
		})
	}
}

func TestPostToAccountsRoute(t *testing.T) {
	tests := []struct {
		name             string
		data             []byte
		expectedResponse []PostAccountResponse
	}{
		{"default", []byte(`[{ "user_id":1,  "name":"Joe Smith",  "date_of_birth":"1983-05-12",  "created_on":1642612034 },{ "user_id":2,  "name":"Jane Doe",  "date_of_birth":"1990-08-06",  "created_on":1642612034 }]`), []PostAccountResponse{
			{UserID: 1, Name: "Joe Smith", DayofTheWeek: "Thursday", RFC3339EST: "2022-01-19T12:07:14-05:00"}, {UserID: 2, Name: "Jane Doe", DayofTheWeek: "Monday", RFC3339EST: "2022-01-19T12:07:14-05:00"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := SetUpRouter()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/accounts", bytes.NewReader(tt.data))
			router.ServeHTTP(w, req)
			var actualResponse []PostAccountResponse
			err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
			failOnError(err, "failed to unmarshal body")
			if !reflect.DeepEqual(actualResponse, tt.expectedResponse) {
				t.Errorf("FAILED: expected: %v, got %v\n", tt.expectedResponse, actualResponse)
			}
		})
	}
}

func TestConvertJPEGToPNG(t *testing.T) {
	tests := []struct {
		name     string
		jpegPath string
		pngPath  string
	}{
		{"default", "data.jpeg", "response.png"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := SetUpRouter()

			// get data for jpeg
			jpegBytes, err := ioutil.ReadFile(tt.jpegPath)
			failOnError(err, "failed to open jpegPath")

			req:= httptest.NewRequest("POST", "/convert", bytes.NewReader(jpegBytes))
			req.Header.Add("Content-Type", "multipart/form-data")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			jpegMeta, err := jpeg.DecodeConfig(bytes.NewReader(jpegBytes))
			failOnError(err, "failed on decode config for jpeg")
			// encode as png and then ensure it has the correct size and aspect ratio
			b := w.Body.Bytes()
			// pngImage, err:=png.Decode(bytes.NewReader(b))
			pngMeta, err := png.DecodeConfig(bytes.NewReader(b))
			failOnError(err, "failed to decode png")
			if pngMeta.Width > 256 || pngMeta.Height > 256 {
				t.Errorf("FAILED: Expected height and width less than 256 Actual H: %v W: %v", pngMeta.Height, pngMeta.Width)
			}
			jpegAR := CalculateAspectRatio(jpegMeta.Height, jpegMeta.Width)
			pngAr := CalculateAspectRatio(pngMeta.Height, pngMeta.Width)
			if pngAr != jpegAR {
				t.Errorf("FAILED: Aspect Ratios are not the same. Expected: %v Actual: %v", jpegAR, pngAr)
			}

		})
	}
}
