package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gangsaur.com/billing-exercise/internal/repository/db/psql"
	"gangsaur.com/billing-exercise/internal/service"
)

// Request

type PayLoanRequest struct {
	Amount int `json:"amount"`
}

func DecodeJSON(w http.ResponseWriter, r *http.Request, v any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return err
	}
	return nil
}

// Response

type LoanResponse struct {
	Id                int     `json:"id"`
	Duration          int     `json:"duration"`
	PrincipalAmount   int     `json:"principal_amount"`
	OutstandingAmount int     `json:"outstanding_amount"`
	Status            int     `json:"status"`
	InterestRate      float32 `json:"interest"`
	UserId            int     `json:"user_id"`
}

func toLoanResponse(loan psql.Loan) LoanResponse {
	return LoanResponse{
		Id:                loan.Id,
		Duration:          loan.Duration,
		PrincipalAmount:   loan.PrincipalAmount,
		OutstandingAmount: loan.OutstandingAmount,
		Status:            loan.Status,
		InterestRate:      loan.InterestRate,
		UserId:            loan.UserId,
	}
}

// Handler

type LoanHandler struct {
	loanService *service.LoanService
}

func NewLoanHandler(l *service.LoanService) *LoanHandler {
	return &LoanHandler{
		loanService: l,
	}
}

func (l *LoanHandler) GetLoan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			WriteErrorResponse(r.Context(), w, r, err)
			return
		}

		loan, err := l.loanService.GetLoan(r.Context(), id)
		if err != nil {
			WriteErrorResponse(r.Context(), w, r, err)
			return
		}

		res, err := json.Marshal(toLoanResponse(loan))
		if err != nil {
			WriteErrorResponse(r.Context(), w, r, err)
			return
		}

		w.Write(res)
	}
}

func (l *LoanHandler) PayLoan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			WriteErrorResponse(r.Context(), w, r, err)
			return
		}

		var payLoanRequest PayLoanRequest
		err = DecodeJSON(w, r, &payLoanRequest)
		if err != nil {
			WriteErrorResponse(r.Context(), w, r, err)
			return
		}

		loan, err := l.loanService.PayLoan(r.Context(), id, payLoanRequest.Amount)
		if err != nil {
			WriteErrorResponse(r.Context(), w, r, err)
			return
		}

		res, err := json.Marshal(toLoanResponse(loan))
		if err != nil {
			WriteErrorResponse(r.Context(), w, r, err)
			return
		}

		w.Write(res)
	}
}
