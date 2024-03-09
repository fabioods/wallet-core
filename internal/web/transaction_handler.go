package web

import (
	"encoding/json"
	"github.com/fabioods/fc-ms-wallet/internal/usecase/create_transaction"
	"net/http"
)

type TransactionHandlerWeb struct {
	CreateTransactionUseCase create_transaction.CreateTransactionUseCase
}

func NewTransactionHandlerWeb(createTransactionUseCase create_transaction.CreateTransactionUseCase) *TransactionHandlerWeb {
	return &TransactionHandlerWeb{
		CreateTransactionUseCase: createTransactionUseCase,
	}
}

func (t *TransactionHandlerWeb) CreateTransactionHandlerWeb(w http.ResponseWriter, r *http.Request) {
	var dto create_transaction.CreateTransactionInputDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	output, err := t.CreateTransactionUseCase.Execute(dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
