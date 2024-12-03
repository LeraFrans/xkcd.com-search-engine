package handler

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"task7/pkg/service"
// )

// func TestHandleLogin(t *testing.T) {
// 	// Create a new HTTP request with POST method
// 	req, err := http.NewRequest("POST", "/login", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Create a new ResponseRecorder to record the response
// 	rr := httptest.NewRecorder()

// 	// Create a new instance of Handler
// 	services := &service.Service{}
// 			handler := Handler{services}

// 	// handler := &Handler{
// 	// 	services: services,// initialize your services here,
// 	// }

// 	// Call the handleLogin function with the ResponseRecorder and the test request
// 	handler.handleLogin(rr, req)

// 	// Check the response status code
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}

// 	// Check the response body
// 	expected := `{"token":"your_token_string"}`
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
// 	}

// 	// Check the Authorization header
// 	if rr.Header().Get("Authorization") != "your_token_string" {
// 		t.Errorf("handler returned unexpected Authorization header: got %v want %v", rr.Header().Get("Authorization"), "your_token_string")
// 	}

// 	// Add more tests here as needed
// }