package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"gangsaur.com/billing-exercise/internal/repository/db/psql"
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

	LoanPayments []LoanPaymentResponse `json:"loan_payments,omitempty"`
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

func toLoanResponseWithLoanPayments(loan psql.Loan, loanPayments []psql.LoanPayment) LoanResponse {
	loanPaymentResponses := make([]LoanPaymentResponse, 0, len(loanPayments))
	for _, lp := range loanPayments {
		loanPaymentResponses = append(loanPaymentResponses, toLoanPaymentResponse(lp))
	}

	return LoanResponse{
		Id:                loan.Id,
		Duration:          loan.Duration,
		PrincipalAmount:   loan.PrincipalAmount,
		OutstandingAmount: loan.OutstandingAmount,
		Status:            loan.Status,
		InterestRate:      loan.InterestRate,
		UserId:            loan.UserId,

		LoanPayments: loanPaymentResponses,
	}
}

type LoanPaymentResponse struct {
	Id      int
	Period  int
	Amount  int
	DueDate time.Time
	PaidAt  time.Time
	Status  int
	LoanId  int
}

func toLoanPaymentResponse(loanPayment psql.LoanPayment) LoanPaymentResponse {
	return LoanPaymentResponse{
		Id:      loanPayment.Id,
		Period:  loanPayment.Period,
		Amount:  loanPayment.Amount,
		DueDate: loanPayment.DueDate,
		PaidAt:  loanPayment.PaidAt,
		Status:  loanPayment.Status,
		LoanId:  loanPayment.LoanId,
	}
}

// Interface

type LoanService interface {
	GetLoan(ctx context.Context, id int) (psql.Loan, error)
	GetLoanAndLoanPayments(ctx context.Context, id int) (psql.Loan, []psql.LoanPayment, error)
	PayLoan(ctx context.Context, id int, amount int) (psql.Loan, error)
}

// Handler

type LoanHandler struct {
	loanService LoanService
}

func NewLoanHandler(l LoanService) *LoanHandler {
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

func (l *LoanHandler) GetLoanAndLoanPayments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			WriteErrorResponse(r.Context(), w, r, err)
			return
		}

		loan, loanPayments, err := l.loanService.GetLoanAndLoanPayments(r.Context(), id)
		if err != nil {
			WriteErrorResponse(r.Context(), w, r, err)
			return
		}

		res, err := json.Marshal(toLoanResponseWithLoanPayments(loan, loanPayments))
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
