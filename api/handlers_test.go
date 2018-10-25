package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dfroese-korewireless/continuous-demo/api"
	"github.com/dfroese-korewireless/continuous-demo/messages"
	"github.com/dfroese-korewireless/continuous-demo/storage"
	"github.com/gorilla/mux"
)

const (
	dbLocation = "unit_test.database"
)

var apiCtx api.Accessor

func setup() {
	fmt.Println("Running setup...")

	db, err := storage.New(dbLocation)
	if err != nil {
		panic(err)
	}
	apiCtx = api.New(db)
}

func teardown() {
	fmt.Println("Running teardown...")
	err := os.Remove(dbLocation)
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	setup()

	retCode := m.Run()

	teardown()
	os.Exit(retCode)
}

var createTests = []struct {
	ExpectedStatusCode  int
	ExpectedContentType string
	RequestBody         string
	RequestMethod       string
	RequestURL          string
}{{
	ExpectedStatusCode:  http.StatusCreated,
	ExpectedContentType: "application/json",
	RequestBody:         `{"Text":"Hello, World", "Username":"Banksy"}`,
	RequestMethod:       "POST",
	RequestURL:          "/api/v1/messages",
}, {
	ExpectedStatusCode:  http.StatusCreated,
	ExpectedContentType: "application/json",
	RequestBody:         `{"Text":"Banksy was here", "Username":"not banksy"}`,
	RequestMethod:       "POST",
	RequestURL:          "/api/v1/messages",
}}

func TestCreateMessage(t *testing.T) {
	for _, tc := range createTests {
		req, err := http.NewRequest(tc.RequestMethod, tc.RequestURL, bytes.NewBuffer([]byte(tc.RequestBody)))
		if err != nil {
			t.Fatalf("unable to create request for method %v url %v with body %v", tc.RequestMethod, tc.RequestURL, tc.RequestBody)

		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(apiCtx.CreateMessage)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != tc.ExpectedStatusCode {
			t.Errorf("handler returned wrong status code: got %v want %v", status, tc.ExpectedStatusCode)
		}

		if content := rr.Header().Get("Content-Type"); content != tc.ExpectedContentType {
			t.Errorf("handler returned wrong content-type: got %v want %v", content, tc.ExpectedContentType)
		}
	}
}

var getAllTest = struct {
	ExpectedStatusCode  int
	ExpectedContentType string
	ExpectedCount       int
	RequestMethod       string
	RequestURL          string
}{
	ExpectedStatusCode:  http.StatusOK,
	ExpectedContentType: "application/json",
	ExpectedCount:       2,
	RequestMethod:       "GET",
	RequestURL:          "/api/v1/messages",
}

func TestGetMessages(t *testing.T) {
	req, err := http.NewRequest(getAllTest.RequestMethod, getAllTest.RequestURL, nil)
	if err != nil {
		t.Fatalf("unable to create request for method %v url %v", getAllTest.RequestMethod, getAllTest.RequestURL)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apiCtx.GetMessages)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != getAllTest.ExpectedStatusCode {
		t.Errorf("handler return wrong status code: got %v want %v", status, getAllTest.ExpectedStatusCode)
	}

	if content := rr.Header().Get("Content-Type"); content != getAllTest.ExpectedContentType {
		t.Errorf("handler returned wrong content-type: got %v want %v", content, getAllTest.ExpectedContentType)
	}

	var msgs []messages.Message
	err = json.Unmarshal([]byte(rr.Body.String()), &msgs)
	if err != nil {
		t.Fatalf("unable to parse the response into messages array")
	}

	if len(msgs) != getAllTest.ExpectedCount {
		t.Errorf("handler returned the wrong count of messages: got %v want %v", len(msgs), getAllTest.ExpectedCount)
	}
}

var getMessageTests = []struct {
	ExpectedStatusCode  int
	ExpectedContentType string
	ExpectedUsername    string
	ExpectedID          uint64
	ExpectedText        string
	RequestMethod       string
	RequestURL          string
}{{
	ExpectedStatusCode:  200,
	ExpectedContentType: "application/json",
	ExpectedUsername:    "Banksy",
	ExpectedID:          1,
	ExpectedText:        "Hello, World",
	RequestMethod:       "GET",
	RequestURL:          "/api/v1/messages/1",
}, {
	ExpectedStatusCode:  200,
	ExpectedContentType: "application/json",
	ExpectedUsername:    "not banksy",
	ExpectedID:          2,
	ExpectedText:        "Banksy was here",
	RequestMethod:       "GET",
	RequestURL:          "/api/v1/messages/2",
}, {
	ExpectedStatusCode:  404,
	ExpectedContentType: "text/plain; charset=utf-8",
	RequestMethod:       "GET",
	RequestURL:          "/api/v1/messages/3",
}}

func TestGetMessage(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/messages/{id}", apiCtx.GetMessage)

	for _, tc := range getMessageTests {
		req, err := http.NewRequest(tc.RequestMethod, tc.RequestURL, nil)
		if err != nil {
			t.Fatalf("unable to create request for method %v url %v", tc.RequestMethod, tc.RequestURL)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if content := rr.Header().Get("Content-Type"); content != tc.ExpectedContentType {
			t.Errorf("handler returned wrong content-type: got %v want %v", content, tc.ExpectedContentType)
		}
		if status := rr.Code; status != tc.ExpectedStatusCode {
			t.Errorf("handler return wrong status code: got %v want %v", status, tc.ExpectedStatusCode)
		}

		if rr.Code != http.StatusOK {
			continue
		}

		var msg messages.Message
		err = json.Unmarshal([]byte(rr.Body.String()), &msg)
		if err != nil {
			t.Fatalf("unable to parse the response into message: %v", err)
		}

		if msg.Username != tc.ExpectedUsername {
			t.Errorf("handler returned the wrong Username: got %v want %v", msg.Username, tc.ExpectedUsername)
		}
		if msg.ID != tc.ExpectedID {
			t.Errorf("handler returned the wrong ID: got %v want %v", msg.ID, tc.ExpectedID)
		}
		if msg.Text != tc.ExpectedText {
			t.Errorf("handler returned the wrong Text: got %v want %v", msg.Text, tc.ExpectedText)
		}
	}
}
