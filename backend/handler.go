package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func main() {
	http.HandleFunc("/auth", handlerFunc)
	handler := cors.Default().Handler(http.HandlerFunc(handlerFunc))
	http.ListenAndServe(":8080", handler)
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	var result string
	if isRegistered(user.Login) {
		if checkPassword(user) {
			result = "Password is correct"
		} else {
			result = "Password is incorrect"
		}
	}
	if !isRegistered(user.Login) {
		result = "User created"
		createUser(user)
	}
	json.NewEncoder(w).Encode(struct {
		Result string `json:"result"`
	}{
		Result: result,
	})
}

func isRegistered(login string) bool {
	return findUser(login).Login != ""
}

func createUser(data User) {
	// Хешируем пароль
	hashedPassword := hashPassword(data.Password)

	var usersFile, err = os.ReadFile("users.json")
	if err != nil {
		panic(err)
	}

	var users []User
	err = json.Unmarshal(usersFile, &users)
	if err != nil {
		panic(err)
	}

	users = append(users, User{data.Login, hashedPassword}) // Используем хеш пароля

	usersFile, err = json.Marshal(users)
	if err != nil {
		panic(err)
	}
	os.WriteFile("users.json", usersFile, 0644)
}

func checkPassword(user User) bool {
	userInFile := findUser(user.Login)
	return hashPassword(user.Password) == userInFile.Password // Сравниваем хеши
}

func findUser(login string) User {
	var usersFile, err = os.ReadFile("users.json")
	if err != nil {
		log.Fatal(err)
	}

	var users []User
	err = json.Unmarshal(usersFile, &users)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(usersFile, &users)
	if err != nil {
		log.Fatal(err)
	}

	for _, value := range users {
		if value.Login == login {
			return value
		}
	}
	return User{"", ""}
}

func hashPassword(password string) string {
	hashedPassword := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hashedPassword[:]) // Возвращаем хеш в шестнадцатеричном формате
}
