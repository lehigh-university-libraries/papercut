package doi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestJoinDate(t *testing.T) {
	tests := []struct {
		input DateParts
		want  string
	}{
		// Test cases for valid input dates
		{
			DateParts{
				Dates: [][]int{
					{
						2022,
						1,
						15,
					},
				},
			},
			"2022-01-15",
		},
		{
			DateParts{
				Dates: [][]int{
					{
						2022,
						9,
					},
				},
			},
			"2022-09",
		},
		{
			DateParts{
				Dates: [][]int{
					{
						2022,
					},
				},
			},
			"2022",
		},
		// Test cases for invalid input dates
		{
			DateParts{
				Dates: [][]int{{}},
			},
			"invalid date",
		},
		{
			DateParts{
				Dates: [][]int{
					{
						2022,
						1,
						32,
					},
				},
			},
			"invalid date",
		},
	}

	for _, test := range tests {
		got := JoinDate(test.input)
		if got != test.want {
			t.Errorf("JoinDate(%v) = %v; want %v", test.input, got, test.want)
		}
	}
}

func TestGetObject(t *testing.T) {
	// Create an example article JSON
	article := Article{
		ReferenceCount: 1,
		Publisher:      "Test Publisher",
		Issue:          "Test Issue",
		Title:          "Test Title",
	}

	// Marshal the article into JSON
	articleJSON, err := json.Marshal(article)
	if err != nil {
		t.Errorf("Error marshaling article JSON: %v", err)
		return
	}

	// Mock server to simulate HTTP responses
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/test1" {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Test response 1"))
		} else if r.URL.Path == "/test2" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			// Write the article JSON to the response
			w.Write(articleJSON)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	tests := []struct {
		url               string
		acceptContentType string
		want              []byte
		wantErr           bool
	}{
		// Test case for plain text response
		{url: ts.URL + "/test1", acceptContentType: "text/plain", want: []byte("Test response 1"), wantErr: false},
		// Test case for JSON response
		{url: ts.URL + "/test2", acceptContentType: "application/json", want: articleJSON, wantErr: false},
		// Test case for non-existent URL
		{url: ts.URL + "/notfound", acceptContentType: "text/plain", want: nil, wantErr: true},
	}

	for _, test := range tests {
		got, err := GetObject(test.url, test.acceptContentType)
		if (err != nil) != test.wantErr {
			t.Errorf("GetObject(%s, %s) error = %v, wantErr %v", test.url, test.acceptContentType, err, test.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetObject(%s, %s) = %v, want %v", test.url, test.acceptContentType, got, test.want)
		}
	}
}
