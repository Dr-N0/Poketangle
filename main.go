package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// RequestExpect represents the expected request body structure
type RequestExpect struct {
	Pokemon  string `json:"pokemon"`
	Question string `json:"question"`
}

type PokemonData map[string][]interface{}

func loadPokemonData() (PokemonData, error) {
	// Read the JSON file
	jsonData, err := os.ReadFile("pokedex.json")
	if err != nil {
		return nil, err
	}

	// Define a variable of type PokemonData to hold the parsed data
	var pokemonData PokemonData

	// Parse the JSON data into the pokemonData variable
	err = json.Unmarshal(jsonData, &pokemonData)
	if err != nil {
		return nil, err
	}

	return pokemonData, nil
}

func main() {
	// Set the secret token as an environment variable
	secretToken := os.Getenv("MY_SECRET_TOKEN")
	if secretToken == "" {
		fmt.Println("Please set the MY_SECRET_TOKEN environment variable.")
		return
	}

	// Define your routes and corresponding handler functions
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/guess", authMiddleware(guessHandler, secretToken))

	// Start the server on port 8080
	port := 8080
	fmt.Printf("Server is running on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

// Handler for the home route "/"
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page!")
}

// Handler for the protected "/guess" route
func guessHandler(w http.ResponseWriter, r *http.Request) {

	pokemonData, pokemonError := loadPokemonData()
	if pokemonError != nil {
		log.Fatalf("Failed to load Pokemon data: %v", pokemonError)
	}

	// Decode the request body into a RequestExpect struct
	var requestBody RequestExpect
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	// Access the Pokemon and Question fields from the request body
	pokemon := requestBody.Pokemon
	question := requestBody.Question

	// Retrieve Pokemon details from the provided JSON data
	pokemonDetails, found := pokemonData[pokemon]
	if !found {
		http.Error(w, "Pokemon not found", http.StatusNotFound)
		return
	}

	// Process different types of questions
	var response string
	switch question {
	case "is it a monotype?":
		response = fmt.Sprintf("%s is %s", pokemon, monotypeCheck(pokemonDetails))
	case "is it a combination [physical, special][physical, special] type":
		response = fmt.Sprintf("%s is %s", pokemon, combinationCheck(pokemonDetails, "physical", "special"))
	case "does it have a [special, physical] type":
		response = fmt.Sprintf("%s %s", pokemon, typeCheck(pokemonDetails, "special", "physical"))
	// Add more cases for other questions...

	default:
		response = fmt.Sprintf("Sorry, I don't understand the question: %s", question)
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"response": response})
}

// Check if the Pokemon is monotype
func monotypeCheck(details []interface{}) string {
	if details[1] == details[2] {
		return "a monotype"
	}
	return "not a monotype"
}

// Check if the Pokemon has a combination of types
func combinationCheck(details []interface{}, type1, type2 string) string {
	if details[1] == type1 && details[2] == type1 && details[1] == type2 && details[2] == type2 {
		return fmt.Sprintf("a combination [%s, %s][%s, %s] type", type1, type2, type1, type2)
	}
	return fmt.Sprintf("not a combination [%s, %s][%s, %s] type", type1, type2, type1, type2)
}

// Check if the Pokemon has a specific type
func typeCheck(details []interface{}, type1, type2 string) string {
	if details[1] == type1 || details[2] == type1 {
		if details[1] == type2 || details[2] == type2 {
			return fmt.Sprintf("a %s %s type", type1, type2)
		}
		return fmt.Sprintf("a %s type", type1)
	} else if details[1] == type2 || details[2] == type2 {
		return fmt.Sprintf("a %s type", type2)
	}
	return fmt.Sprintf("not a %s or %s type", type1, type2)
}

func authMiddleware(next http.HandlerFunc, token string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check the Authorization header for the token
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer "+token {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Call the next handler if authentication is successful
		next(w, r)
	}
}
