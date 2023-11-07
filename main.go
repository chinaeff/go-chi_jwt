package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/docgen"
	"github.com/go-chi/jwtauth"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

type User struct {
	Username string
	Password string
}

var users = map[string]User{}

type RequestAddressSearch struct {
	Query string `json:"query"`
}

type ResponseAddress struct {
	Addresses []*Address `json:"addresses"`
}

type RequestAddressGeocode struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type Address struct {
	City string `json:"city"`
}

func SearchAddressHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	username, _ := claims["username"].(string)

	if username != "" {
		var request RequestAddressSearch
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Wrong data", http.StatusBadRequest)
			return
		}

		response := &ResponseAddress{
			Addresses: []*Address{
				{
					City: "Moscow",
				},
			},
		}

		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	} else {
		http.Error(w, "Error", http.StatusForbidden)
	}
}

func GeocodeAddressHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	username, _ := claims["username"].(string)

	if username != "" {
		var request RequestAddressGeocode
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Wrong data", http.StatusBadRequest)
			return
		}

		response := &ResponseAddress{
			Addresses: []*Address{
				{
					City: "Moscow",
				},
			},
		}

		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	} else {
		http.Error(w, "Error", http.StatusForbidden)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Wrong data", http.StatusBadRequest)
		return
	}

	if _, exists := users[user.Username]; exists {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)

	users[user.Username] = user
	w.WriteHeader(http.StatusCreated)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Wrong data", http.StatusBadRequest)
		return
	}

	storedUser, exists := users[user.Username]
	if !exists || bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)) != nil {
		http.Error(w, "Error", http.StatusUnauthorized)
		return
	}

	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"username": user.Username})
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func main() {
	r := chi.NewRouter()

	r.Use(jwtauth.Verifier(tokenAuth))
	r.Use(jwtauth.Authenticator)

	r.Post("/api/login", LoginHandler)
	r.Post("/api/register", RegisterHandler)

	r.Post("/api/address/search", SearchAddressHandler)
	r.Post("/api/address/geocode", GeocodeAddressHandler)

	docgen.PrintRoutes(r)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
