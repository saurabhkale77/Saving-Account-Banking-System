package account

import (
	"Saving-Account-Banking-System/app/specs"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func Deposit(accService Service) func(w http.ResponseWriter, r *http.Request) { //Post
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tknStr := r.Header.Get("Authorization")
		user_id, _, err := accService.Authenticate(tknStr)
		if err != nil {
			specs.ErrorUnauthorizedAccess(err, w)
			return
		}

		var req specs.Transaction
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			specs.ErrorInternalServer(err, w)
			return
		}

		err = req.ValidateTransaction()
		if err != nil {
			specs.ErrorBadRequest(err, w)
			return
		}

		//var result specs.Transaction
		result, err := accService.DepositMoney(ctx, req, user_id)
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

func Withdrawal(accService Service) func(w http.ResponseWriter, r *http.Request) { //PUT
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tknStr := r.Header.Get("Authorization")
		user_id, _, err := accService.Authenticate(tknStr)
		if err != nil {
			specs.ErrorUnauthorizedAccess(err, w)
			return
		}

		var req specs.Transaction
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			specs.ErrorInternalServer(err, w)
			return
		}

		err = req.ValidateTransaction()
		if err != nil {
			specs.ErrorBadRequest(err, w)
			return
		}

		result, err := accService.WithdrawalMoney(ctx, req, user_id)
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

func Delete(accService Service) func(w http.ResponseWriter, r *http.Request) { //Post
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tknStr := r.Header.Get("Authorization")
		user_id, role, err := accService.Authenticate(tknStr)
		if err != nil {
			specs.ErrorUnauthorizedAccess(err, w)
			return
		}

		var req specs.DeleteAccountReq
		// err = json.NewDecoder(r.Body).Decode(&req)
		queryParams := r.URL.Query()
		paramValue := queryParams.Get("acc_no")
		acc_no, err := strconv.Atoi(paramValue)
		req.Account_no = acc_no
		if err != nil {
			specs.ErrorInternalServer(err, w)
			return
		}

		if strings.EqualFold(role, "admin") {
			paramValue = queryParams.Get("user_id")
			user_id, err = strconv.Atoi(paramValue)
			if err != nil {
				specs.ErrorInternalServer(err, w)
				return
			}
		}

		err = req.ValidateDeleteReq()
		if err != nil {
			specs.ErrorBadRequest(err, w)
			return
		}

		response, err := accService.DeleteAccount(ctx, req, user_id)
		if err != nil {
			specs.ErrorBadRequest(err, w)
			return
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			specs.ErrorInternalServer(err, w)
			return
		}
	}
}

func ViewBalance(accService Service) func(w http.ResponseWriter, r *http.Request) { //GET pathparam
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tknStr := r.Header.Get("Authorization")
		user_id, _, err := accService.Authenticate(tknStr)
		if err != nil {
			specs.ErrorUnauthorizedAccess(err, w)
			return
		}

		var req specs.TransactionResponse
		// err = json.NewDecoder(r.Body).Decode(&req)
		queryParams := r.URL.Query()
		paramValue := queryParams.Get("acc_no")
		acc_no, err := strconv.Atoi(paramValue)
		req.Account_no = acc_no
		if err != nil {
			specs.ErrorInternalServer(err, w)
			return
		}

		response, err := accService.ViewBalance(ctx, req, user_id)
		if err != nil {
			specs.ErrorBadRequest(err, w)
			return
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			specs.ErrorInternalServer(err, w)
			return
		}
	}
}
