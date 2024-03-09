package create_transaction

import (
	"github.com/fabioods/fc-ms-wallet/internal/entity"
	"github.com/fabioods/fc-ms-wallet/internal/gateway"
	"github.com/fabioods/fc-ms-wallet/pkg/events"
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
	TransactionCreated events.EventInterface
	EventsDispatcher   events.EventDispatcherInterface
}

func NewCreateTransactionUseCase(tg gateway.TransactionGateway, ag gateway.AccountGateway, tc events.EventInterface, ed events.EventDispatcherInterface) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		TransactionGateway: tg,
		AccountGateway:     ag,
		TransactionCreated: tc,
		EventsDispatcher:   ed,
	}
}

func (uc *CreateTransactionUseCase) Execute(input CreateTransactionInputDto) (*CreateTransactionOutputDto, error) {
	accountFrom, err := uc.AccountGateway.FindByID(input.AccountIdFrom)
	if err != nil {
		return nil, err
	}
	accountTo, err := uc.AccountGateway.FindByID(input.AccountIdTo)
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

	output := &CreateTransactionOutputDto{Id: transaction.ID}

	uc.TransactionCreated.SetPayload(output)
	err = uc.EventsDispatcher.Dispatch(uc.TransactionCreated)
	if err != nil {
		return nil, err
	}

	return output, nil
}
