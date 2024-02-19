package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var _users_storage = make([]User, 0, 1000000)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Users struct {
	Users []User `json:"users"`
}

type Error struct {
	Error string `json:"errors"`
}

type Response struct {
	Data       interface{} `json:"data"`
	StatusCode int         `json:"statusCode"`
	ServerTime int64       `json:"serverTime"`
}

func logg(ltype string, request *http.Request) {
	ansiColor := "\033[1;34m" // default INFO
	ansiColorClear := "\033[0m"

	rfcTimestamp := time.Now().Format(time.RFC3339)

	userHost := request.Host
	userPath := request.URL.Path
	userMethod := request.Method

	switch ltype {
	case "DEBUG":
		ansiColor = "\033[0;36m"
		// ansiColor = "\033[0;36m%s"
	case "INFO":
		ansiColor = "\033[0;34m"
	case "WARNING":
		ansiColor = "\033[0;33m"
	case "NOTICE":
		ansiColor = "\033[1;36m"
	case "ERROR":
		ansiColor = "\033[0;31m"
	}

	fmt.Printf(
		"\r%s | %s%s%s | %s%s %s\n",
		rfcTimestamp,
		ansiColor,
		ltype,
		ansiColorClear,
		userHost,
		userPath,
		userMethod,
	)
}

func usersEndpoints(writer http.ResponseWriter, request *http.Request) {
	logg("INFO", request)

	writer.Header().Set("Content-Type", "application/json")

	_serverTime := time.Now().UnixMilli()
	_statusCode := http.StatusOK

	switch request.Method {
	case "GET":

		_response := &Response{
			Data:       _users_storage,
			StatusCode: _statusCode,
			ServerTime: _serverTime,
		}

		response, err := json.Marshal(_response)

		if err != nil {
			logg("ERROR", request)
			return
		}

		writer.Write(response)
		return

	case "POST":
		_raw_amount := request.URL.Query().Get("amount")

		if _raw_amount == "" {
			_statusCode := http.StatusBadRequest
			_error := &Error{
				Error: "empty parameter",
			}

			_response := &Response{
				Data:       _error,
				StatusCode: _statusCode,
				ServerTime: _serverTime,
			}

			response, err := json.Marshal(_response)

			if err != nil {
				return
			}

			writer.Write(response)
			return
		}

		_amount, err := strconv.Atoi(_raw_amount)

		if err != nil {
			_statusCode := http.StatusBadRequest
			_error := &Error{
				Error: "amount must be integer",
			}

			_response := &Response{
				Data:       _error,
				StatusCode: _statusCode,
				ServerTime: _serverTime,
			}

			response, err := json.Marshal(_response)

			if err != nil {
				logg("ERROR", request)
				return
			}

			writer.Write(response)
			return
		}

		_users_local := make([]User, _amount)

		for amount := 0; amount < _amount; amount++ {
			_id := rand.Intn(999999999)
			_name := fmt.Sprintf("testuser%d", _id)

			_user := &User{
				Id:   _id,
				Name: _name,
			}

			_users_local = append(_users_local, *_user)
		}

		_users_storage := append(_users_storage, _users_local...)

		_response := &Response{
			Data:       _users_storage,
			StatusCode: _statusCode,
			ServerTime: _serverTime,
		}

		response, err := json.Marshal(_response)

		if err != nil {
			logg("ERROR", request)
			return
		}

		writer.Write(response)
		return

	case "DELETE":

	default:
		_statusCode := http.StatusNotFound
		_error := &Error{
			Error: "whaaat?",
		}

		_response := &Response{
			Data:       _error,
			StatusCode: _statusCode,
			ServerTime: _serverTime,
		}

		response, err := json.Marshal(_response)

		if err != nil {
			logg("ERROR", request)
			return
		}

		writer.Write(response)
		return
	}
}

func main() {
	http.HandleFunc("/users/", usersEndpoints)
	http.ListenAndServe(":8090", nil)
}
