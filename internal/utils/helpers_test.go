package utils

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestFetchEmails(t *testing.T) {
	testCases := []struct {
		name             string
		mockURL          string
		mockResponseBody string
		expectedEmails   []string
		expectedErr      error
	}{
		{
			name:             "Valid Response",
			mockURL:          "/mock/emails",
			mockResponseBody: "Contact us at info@example.com or support@example.com",
			expectedEmails:   []string{"info@example.com", "support@example.com"},
			expectedErr:      nil,
		},
		{
			name:             "No Emails",
			mockURL:          "/mock/noemails",
			mockResponseBody: "This is a sample text without any email addresses.",
			expectedEmails:   []string{},
			expectedErr:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new mock server
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == tc.mockURL {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(tc.mockResponseBody))
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer mockServer.Close()

			// Call the function with the mock server URL
			emails, err := FetchEmails(mockServer.URL + tc.mockURL)

			// Check the result against the expectations
			if err != nil && tc.expectedErr == nil {
				t.Errorf("Test case %q failed: unexpected error %v", tc.name, err)
			}
			if err == nil && tc.expectedErr != nil {
				t.Errorf("Test case %q failed: expected error %v, got nil", tc.name, tc.expectedErr)
			}
			if len(emails) != len(tc.expectedEmails) {
				t.Errorf("Test case %q failed: expected %d emails, got %d", tc.name, len(tc.expectedEmails), len(emails))
			}
			for i, expectedEmail := range tc.expectedEmails {
				if emails[i] != expectedEmail {
					t.Errorf("Test case %q failed: expected email %q, got %q", tc.name, expectedEmail, emails[i])
				}
			}
		})
	}
}

func TestMkTmpDir(t *testing.T) {
	// Test case: directory does not exist
	testDir := "test_dir"
	tmpDir, err := MkTmpDir(testDir)
	if err != nil {
		t.Errorf("Error creating temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Verify directory exists
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Errorf("Temporary directory was not created: %v", err)
	}

	// Test case: directory already exists
	tmpDir2, err := MkTmpDir(testDir)
	if err != nil {
		t.Errorf("Error creating temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir2)

	// Verify directory exists
	if _, err := os.Stat(tmpDir2); os.IsNotExist(err) {
		t.Errorf("Temporary directory was not created: %v", err)
	}
}

func TestTrimToMaxLen(t *testing.T) {
	// Test case: string is shorter than maxLen
	inputShort := "short string"
	expectedShort := "short string"
	resultShort := TrimToMaxLen(inputShort, 20)
	if resultShort != expectedShort {
		t.Errorf("TrimToMaxLen(%q, 20) = %q; want %q", inputShort, resultShort, expectedShort)
	}

	// Test case: string is longer than maxLen
	inputLong := "this is a long string that exceeds the max length"
	expectedLong := "this is a long strin"
	resultLong := TrimToMaxLen(inputLong, 20)
	if resultLong != expectedLong {
		t.Errorf("TrimToMaxLen(%q, 20) = %q; want %q", inputLong, resultLong, expectedLong)
	}

	// Test case: string is exactly maxLen
	inputExact := "exact length string"
	expectedExact := "exact length string"
	resultExact := TrimToMaxLen(inputExact, 20)
	if resultExact != expectedExact {
		t.Errorf("TrimToMaxLen(%q, 20) = %q; want %q", inputExact, resultExact, expectedExact)
	}
}
