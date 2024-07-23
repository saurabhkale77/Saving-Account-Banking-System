package enduser

import (
	"Saving-Account-Banking-System/app/specs"
	"encoding/json"
	"net/http"
)

func Login(userService Service) func(w http.ResponseWriter, r *http.Request) { //POST
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req specs.CreateLoginRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			specs.ErrorInternalServer(err, w)
			return
		}

		err = req.Validate()
		if err != nil {
			specs.ErrorBadRequest(err, w)
			return
		}

		response, err := userService.CreateLogin(ctx, req)
		if err != nil {
			specs.ErrorBadRequest(err, w)
			return
		}

		res := specs.LoginToken{
			Issuespecsken: response,
		}

		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			specs.ErrorInternalServer(err, w)
			return
		}
	}
}

func Signup(userService Service) func(w http.ResponseWriter, r *http.Request) { //Post
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req specs.CreateUser

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			specs.ErrorInternalServer(err, w)
			return
		}

		err = req.ValidateUser()
		if err != nil {
			specs.ErrorBadRequest(err, w)
			return
		}

		result, err := userService.CreateSignup(ctx, req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			specs.ErrorBadRequest(err, w)
			return
		}

		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			specs.ErrorInternalServer(err, w)
			return
		}
	}
}

func Update(userService Service) func(w http.ResponseWriter, r *http.Request) { //PUT
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tknStr := r.Header.Get("Authorization")

		user_id, _, err := userService.Authenticate(tknStr)
		if err != nil {
			specs.ErrorUnauthorizedAccess(err, w)
			return
		}

		//Updating User Info
		var req specs.UpdateUser
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			specs.ErrorBadRequest(err, w)
			return
		}

		err = req.ValidateUpdate()
		if err != nil {
			specs.ErrorBadRequest(err, w)
			return
		}

		result, err := userService.UpdateUser(ctx, req, user_id)
		if err != nil {
			specs.ErrorBadRequest(err, w)
			return
		}

		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			specs.ErrorInternalServer(err, w)
			return
		}
	}
}

func GetUser(userService Service) func(w http.ResponseWriter, r *http.Request) { //PUT
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tknStr := r.Header.Get("Authorization")

		user_id, _, err := userService.Authenticate(tknStr)
		if err != nil {
			specs.ErrorUnauthorizedAccess(err, w)
			return
		}

		result, err := userService.GetUser(ctx, user_id)
		if err != nil {
			specs.ErrorBadRequest(err, w)
			return
		}

		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			specs.ErrorInternalServer(err, w)
			return
		}
	}
}

func GetMyAccounts(userService Service) func(w http.ResponseWriter, r *http.Request) { //PUT
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tknStr := r.Header.Get("Authorization")

		user_id, _, err := userService.Authenticate(tknStr)
		if err != nil {
			specs.ErrorUnauthorizedAccess(err, w)
			return
		}

		result, err := userService.GetMyAccounts(ctx, user_id)
		if err != nil {
			specs.ErrorBadRequest(err, w)
			return
		}

		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			specs.ErrorInternalServer(err, w)
			return
		}
	}
}
