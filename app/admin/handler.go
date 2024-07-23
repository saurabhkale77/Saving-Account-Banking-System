package admin

import (
	"Saving-Account-Banking-System/app/specs"
	"encoding/json"
	"net/http"
)

func ListUsers(AdminService Service) func(w http.ResponseWriter, r *http.Request) { //GET
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tknStr := r.Header.Get("Authorization")
		_, err := AdminService.Authenticate(tknStr)
		if err != nil {
			specs.ErrorUnauthorizedAccess(err, w)
			return
		}

		resp, err := AdminService.ListUsers(ctx)
		if err != nil {
			specs.ErrorBadRequest(err, w)
			return
		}

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			specs.ErrorInternalServer(err, w)
			return
		}
	}
}

func Update(adminService Service) func(w http.ResponseWriter, r *http.Request) { //PUT
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tknStr := r.Header.Get("Authorization")

		_, err := adminService.Authenticate(tknStr)
		if err != nil {
			specs.ErrorUnauthorizedAccess(err, w)
			return
		}

		//Updating User Info
		var req specs.UpdateUserInfo
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			specs.ErrorInternalServer(err, w)
			return
		}

		err = req.ValidateUpdate()
		if err != nil {
			specs.ErrorBadRequest(err, w)
			return
		}

		result, err := adminService.UpdateUser(ctx, req)
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

func CreateAccount(adminService Service) func(w http.ResponseWriter, r *http.Request) { //Post
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tknStr := r.Header.Get("Authorization")
		_, err := adminService.Authenticate(tknStr)
		if err != nil {
			specs.ErrorUnauthorizedAccess(err, w)
			return
		}

		var req specs.CreateAccountReq
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			specs.ErrorInternalServer(err, w)
			return
		}

		err = req.Validate()
		if err != nil {
			specs.ErrorBadRequest(err, w)
			return
		}

		result, err := adminService.CreateAccount(ctx, req)
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

func ListBranches(AdminService Service) func(w http.ResponseWriter, r *http.Request) { //GET
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tknStr := r.Header.Get("Authorization")
		_, err := AdminService.Authenticate(tknStr)
		if err != nil {
			specs.ErrorUnauthorizedAccess(err, w)
			return
		}

		resp, err := AdminService.ListBranches(ctx)
		if err != nil {
			specs.ErrorBadRequest(err, w)
			return
		}

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			specs.ErrorInternalServer(err, w)
			return
		}
	}
}
