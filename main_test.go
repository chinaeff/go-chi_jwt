package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
	user := User{
		Username: "testuser",
		Password: "testpassword",
	}
	userJSON, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "/api/register", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	RegisterHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected %d, but %d", http.StatusCreated, rr.Code)
	}

	if _, exists := users[user.Username]; !exists {
		t.Error("User error")
	}

	if users[user.Username].Password == user.Password {
		t.Error("Password error")
	}
}

func TestLoginHandler(t *testing.T) {
	user := User{
		Username: "testuser",
		Password: "testpassword",
	}
	userJSON, _ := json.Marshal(user)

	reqRegister, err := http.NewRequest("POST", "/api/register", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatal(err)
	}
	rrRegister := httptest.NewRecorder()
	RegisterHandler(rrRegister, reqRegister)

	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	LoginHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected %d, but %d", http.StatusOK, rr.Code)
	}

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	token, ok := response["token"]
	if !ok {
		t.Error("Token error")
	}

	reqWithToken, err := http.NewRequest("POST", "/api/address/search", nil)
	reqWithToken.Header.Set("Authorization", "Bearer "+token)
	if err != nil {
		t.Fatal(err)
	}

	rrWithToken := httptest.NewRecorder()
	SearchAddressHandler(rrWithToken, reqWithToken)

}
