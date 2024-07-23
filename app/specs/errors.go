package specs

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	Error_code int    `json:"error_code"`
	Error_msg  string `json:"error_description"`
}

var e Error

func ErrorInternalServer(err error, w http.ResponseWriter) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		e.Error_code = 500
		e.Error_msg = "error, while decoding data into json, plz provide valid credentails"
		return
	}
}

func ErrorBadRequest(err error, w http.ResponseWriter) {
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		e.Error_code = 400
		e.Error_msg = err.Error()
		_ = json.NewEncoder(w).Encode(e)
		return
	}
}

func ErrorUnauthorizedAccess(err error, w http.ResponseWriter) {
	fmt.Println(err.Error())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		e.Error_code = 401
		e.Error_msg = "Plz, Do Login First"
		_ = json.NewEncoder(w).Encode(e)
		return
	}
}
