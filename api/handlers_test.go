package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dfroese-korewireless/continuous-demo/api"
	"github.com/dfroese-korewireless/continuous-demo/messages"
	"github.com/gorilla/mux"
)

var nextID uint64 = 1

type mockDatabase struct {
	data map[uint64]messages.Message
}

// GetAllMessages mocked implementation
func (db *mockDatabase) GetAllMessages() ([]messages.Message, error) {
	var msgs []messages.Message
	for _, v := range db.data {
		msgs = append(msgs, v)
	}

	return msgs, nil
}

// GetMessage mocked implementation
func (db *mockDatabase) GetMessage(id uint64) (messages.Message, error) {
	msg := db.data[id]
	if (msg == messages.Message{}) {
		return msg, fmt.Errorf("looking up value for %d return nil", id)
	}

	return msg, nil
}

// StoreMessage mocked implementation
func (db *mockDatabase) StoreMessage(msg messages.Message) (uint64, error) {
	msg.ID = nextID
	db.data[nextID] = msg

	nextID += 1
	return msg.ID, nil
}

var apiCtx api.Accessor

func setup() {
	testData := []struct {
		RequestBody   string
		RequestMethod string
		RequestURL    string
	}{{
		RequestBody:   `{"Text":"Hello, World", "Username":"Banksy"}`,
		RequestMethod: "POST",
		RequestURL:    "/api/v1/messages",
	}, {
		RequestBody:   `{"Text":"Banksy was here", "Username":"not banksy"}`,
		RequestMethod: "POST",
		RequestURL:    "/api/v1/messages",
	}}

	nextID = 1
	db := &mockDatabase{data: make(map[uint64]messages.Message)}

	apiCtx = api.New(db)

	for _, d := range testData {
		req, err := http.NewRequest(d.RequestMethod, d.RequestURL, bytes.NewBuffer([]byte(d.RequestBody)))
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(apiCtx.CreateMessage)

		handler.ServeHTTP(rr, req)
	}
}

func TestCreateMessage(t *testing.T) {
	testCases := []struct {
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

	setup()

	for _, tc := range testCases {
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

func TestGetMessages(t *testing.T) {
	testCase := struct {
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

	setup()

	req, err := http.NewRequest(testCase.RequestMethod, testCase.RequestURL, nil)
	if err != nil {
		t.Fatalf("unable to create request for method %v url %v", testCase.RequestMethod, testCase.RequestURL)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apiCtx.GetMessages)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != testCase.ExpectedStatusCode {
		t.Errorf("handler return wrong status code: got %v want %v", status, testCase.ExpectedStatusCode)
	}

	if content := rr.Header().Get("Content-Type"); content != testCase.ExpectedContentType {
		t.Errorf("handler returned wrong content-type: got %v want %v", content, testCase.ExpectedContentType)
	}

	var msgs []messages.Message
	err = json.Unmarshal([]byte(rr.Body.String()), &msgs)
	if err != nil {
		t.Fatalf("unable to parse the response into messages array")
	}

	if len(msgs) != testCase.ExpectedCount {
		t.Errorf("handler returned the wrong count of messages: got %v want %v", len(msgs), testCase.ExpectedCount)
	}
}

func TestGetMessage(t *testing.T) {
	testCases := []struct {
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

	setup()
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/messages/{id}", apiCtx.GetMessage)

	for _, tc := range testCases {
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
