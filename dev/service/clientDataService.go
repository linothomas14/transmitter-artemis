package service

import (
	"transmitter-artemis/entity"
	"transmitter-artemis/repository"
)

type ClientService interface {
	GetAllClientData() ([]entity.ClientData, error)
}

type clientService struct {
	clientRepository repository.ClientRepository
}

func NewClientService(clientRepository repository.ClientRepository) *clientService {
	return &clientService{
		clientRepository: clientRepository,
	}
}

func (cs *clientService) GetAllClientData() ([]entity.ClientData, error) {
	return cs.clientRepository.GetAllClientData()
}
