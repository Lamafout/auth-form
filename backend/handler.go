package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func main() {
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/show", showHandler)

	// любименький cors
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowCredentials: true,
	})
	handler := c.Handler(http.DefaultServeMux)

	http.ListenAndServe(":8080", handler)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	// подключаемся с базе данных
	connectStr := "user=postgres dbname=users_web sslmode=disable password=1234"
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		log.Fatal(err)
	}

	// декодирем пришедший от клиента json
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	var result string
	if isRegistered(user.Login, db) {
		if checkPassword(user, db) {
			result = "Password is correct"
		} else {
			result = "Password is incorrect"
		}
	}
	if !isRegistered(user.Login, db) {
		result = "User created"
		createUser(user, db)
	}

	// возвращаем результат в виде json
	json.NewEncoder(w).Encode(struct {
		Result string `json:"result"`
	}{
		Result: result,
	})
}

func showHandler(w http.ResponseWriter, r *http.Request) {
	// подключаемся с базе данных
	connectStr := "user=postgres dbname=users_web sslmode=disable password=1234"
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		log.Fatal(err)
	}

	// получаем массив пользователей из БД
	result, err := getUsers(db)
	if err != nil {
		log.Fatal(err)
	}

	// преобразуем результат в формат JSON
	usersJSON, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Failed to encode users to JSON", http.StatusInternalServerError)
		return
	}

	// устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

	// отправляем JSON-данные клиенту
	w.Write(usersJSON)
}

func isRegistered(login string, db *sql.DB) bool {
	return findUser(login, db).Login != ""
}

func createUser(data User, db *sql.DB) {
	// Хешируем пароль
	hashedPassword := hashPassword(data.Password)

	_, err := db.Exec("INSERT INTO users (id, login, password) VALUES(default, $1, $2)", data.Login, hashedPassword)
	if err != nil {
		log.Fatal(err)
	}
}

func checkPassword(user User, db *sql.DB) bool {
	userInDB := findUser(user.Login, db)
	return hashPassword(user.Password) == userInDB.Password // Сравниваем хеши
}

func findUser(login string, db *sql.DB) User {
	var user User

	// Выполнение SQL-запроса
	row := db.QueryRow("SELECT login, password FROM users WHERE login=$1;", login)

	// Сканирование результата
	err := row.Scan(&user.Login, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			// Пользователь не найден, возвращаем пустого пользователя
			return User{Login: "", Password: ""}
		} else {
			log.Fatal(err)
		}
	}

	return user
}

func hashPassword(password string) string {
	hashedPassword := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hashedPassword[:]) // Возвращаем хеш в шестнадцатеричном формате
}

func getUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT login, password FROM users;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Login, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
