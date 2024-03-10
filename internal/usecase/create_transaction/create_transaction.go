package create_transaction

import (
	"context"
	"fmt"
	"github.com/fabioods/fc-ms-wallet/internal/entity"
	"github.com/fabioods/fc-ms-wallet/internal/gateway"
	"github.com/fabioods/fc-ms-wallet/pkg/events"
	"github.com/fabioods/fc-ms-wallet/pkg/uow"
)

type CreateTransactionInputDto struct {
	AccountIdFrom string  `json:"account_id_from"`
	AccountIdTo   string  `json:"account_id_to"`
	Amount        float64 `json:"amount"`
}

type CreateTransactionOutputDto struct {
	Id          string  `json:"id"`
	AccountFrom string  `json:"account_from"`
	AccountTo   string  `json:"account_to"`
	Amount      float64 `json:"amount"`
}

// variaveis para DI
type CreateTransactionUseCase struct {
	UnitOfWork         uow.UowInterface
	TransactionCreated events.EventInterface
	EventsDispatcher   events.EventDispatcherInterface
}

func NewCreateTransactionUseCase(uow uow.UowInterface, tc events.EventInterface, ed events.EventDispatcherInterface) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		UnitOfWork:         uow,
		TransactionCreated: tc,
		EventsDispatcher:   ed,
	}
}

func (uc *CreateTransactionUseCase) Execute(ctx context.Context, input CreateTransactionInputDto) (*CreateTransactionOutputDto, error) {
	output := &CreateTransactionOutputDto{}
	err := uc.UnitOfWork.Do(ctx, func(uow *uow.Uow) error {
		accountFrom, err := uc.getAccountRepository(ctx).FindByID(input.AccountIdFrom)
		if err != nil {
			fmt.Println("Error to find account from")
			return err
		}
		accountTo, err := uc.getAccountRepository(ctx).FindByID(input.AccountIdTo)
		if err != nil {
			fmt.Println("Error to find account to")
			return err
		}
		transaction, err := entity.NewTransaction(accountFrom, accountTo, input.Amount)

		if err != nil {
			fmt.Println("Error to create transaction")
			return err
		}

		err = uc.getAccountRepository(ctx).UpdateBalance(accountFrom)
		if err != nil {
			fmt.Println("Error to update account from")
			return err
		}

		err = uc.getAccountRepository(ctx).UpdateBalance(accountTo)
		if err != nil {
			fmt.Println("Error to update account to")
			return err
		}

		err = uc.getTransactionRepository(ctx).Create(transaction)

		if err != nil {
			fmt.Println("Error to execute transaction")
			return err
		}
		output.Id = transaction.ID
		output.AccountFrom = transaction.AccountFrom.ID
		output.AccountTo = transaction.AccountTo.ID
		output.Amount = transaction.Amount
		return nil

	})
	uc.TransactionCreated.SetPayload(output)
	err = uc.EventsDispatcher.Dispatch(uc.TransactionCreated)
	if err != nil {
		fmt.Println("Error to dispatch transaction created")
		return nil, err
	}

	return output, nil
}

func (uc *CreateTransactionUseCase) getAccountRepository(ctx context.Context) gateway.AccountGateway {
	repo, err := uc.UnitOfWork.GetRepository(ctx, "AccountDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.AccountGateway)
}

func (uc *CreateTransactionUseCase) getTransactionRepository(ctx context.Context) gateway.TransactionGateway {
	repo, err := uc.UnitOfWork.GetRepository(ctx, "TransactionDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.TransactionGateway)
}
