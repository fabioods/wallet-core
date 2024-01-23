package create_transaction

import (
	"github.com/fabioods/fc-ms-wallet/internal/entity"
	"github.com/fabioods/fc-ms-wallet/internal/gateway"
)

type CreateTransactionInputDto struct {
	AccountIdFrom string
	AccountIdTo   string
	Amount        float64
}

type CreateTransactionOutputDto struct {
	Id string
}

// variaveis para DI
type CreateTransactionUseCase struct {
	TransactionGateway gateway.TransactionGateway
	AccountGateway     gateway.AccountGateway
}

func NewCreateTransactionUseCase(tg gateway.TransactionGateway, ag gateway.AccountGateway) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		TransactionGateway: tg,
		AccountGateway:     ag,
	}
}

func (uc *CreateTransactionUseCase) Execute(input CreateTransactionInputDto) (*CreateTransactionOutputDto, error) {
	accountFrom, err := uc.AccountGateway.FindById(input.AccountIdFrom)
	if err != nil {
		return nil, err
	}
	accountTo, err := uc.AccountGateway.FindById(input.AccountIdTo)
	if err != nil {
		return nil, err
	}
	transaction, err := entity.NewTransaction(accountFrom, accountTo, input.Amount)

	if err != nil {
		return nil, err
	}
	err = uc.TransactionGateway.Create(transaction)

	if err != nil {
		return nil, err
	}

	return &CreateTransactionOutputDto{Id: transaction.ID}, nil
}
