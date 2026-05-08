package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gangsaur.com/billing-exercise/internal/service"
)

type LoangHandler struct {
	loanService service.LoanService
}

func NewLoanHandler(l service.LoanService) *LoangHandler {
	return &LoangHandler{
		loanService: l,
	}
}

func (l *LoangHandler) GetLoan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			WriteGenericError(r.Context(), w, r, err)
			return
		}

		loanData, err := l.loanService.GetLoan(r.Context(), id)
		if err != nil {
			WriteGenericError(r.Context(), w, r, err)
			return
		}

		res, err := json.Marshal(loanData)
		if err != nil {
			WriteGenericError(r.Context(), w, r, err)
			return
		}

		w.Write(res)
	}
}
