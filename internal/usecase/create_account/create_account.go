package create_account

import (
	"github.com/fabioods/fc-ms-wallet/internal/entity"
	"github.com/fabioods/fc-ms-wallet/internal/gateway"
	"github.com/fabioods/fc-ms-wallet/pkg/events"
)

type CreateAccountInputDto struct {
	ClientId string
}

type CreateAccountOutputDto struct {
	Id string
}

// variaveis da injeção de dependencia
type CreateAccountUseCase struct {
	AccountGateway     gateway.AccountGateway
	ClientGateway      gateway.ClientGateway
	TransactionCreated events.EventInterface
	EventsDispatcher   events.EventDispatcherInterface
}

// construtor
func NewCreateAccountUseCase(ag gateway.AccountGateway, cg gateway.ClientGateway,
	tc events.EventInterface, ed events.EventDispatcherInterface) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		AccountGateway:     ag,
		ClientGateway:      cg,
		TransactionCreated: tc,
		EventsDispatcher:   ed,
	}
}

func (uc *CreateAccountUseCase) Execute(input CreateAccountInputDto) (*CreateAccountOutputDto, error) {
	client, err := uc.ClientGateway.Get(input.ClientId)

	if err != nil {
		return nil, err
	}

	account, err := entity.NewAccount(client)
	if err != nil {
		return nil, err
	}

	err = uc.AccountGateway.Save(account)

	if err != nil {
		return nil, err
	}

	output := &CreateAccountOutputDto{
		Id: account.ID,
	}
	uc.TransactionCreated.SetPayload(output)
	err = uc.EventsDispatcher.Dispatch(uc.TransactionCreated)
	if err != nil {
		return nil, err
	}

	return output, nil

}
