package service

import (
	"testing"
	"transmitter-artemis/entity"
	mocksRepo "transmitter-artemis/mocks/repository"

	"github.com/stretchr/testify/assert"
)

func TestGetAllClientData_All(t *testing.T) {
	type want struct {
		res []entity.ClientData
		err error
	}

	type test struct {
		name               string
		clientDataRepoMock func() (clientDataRepoMock *mocksRepo.ClientRepository)
		want               want
	}

	tests := []test{
		{
			name: "Get all client data",
			clientDataRepoMock: func() (clientDataRepoMock *mocksRepo.ClientRepository) {
				clientDataRepoMock = mocksRepo.NewClientRepository(t)
				clientDataRepoMock.On("GetAllClientData").Return([]entity.ClientData{
					{ClientName: "lino",
						Token:         "abc",
						PhoneNumberID: "123",
						WAHost:        "https://graph.facebook.com",
					},
					{ClientName: "thomas",
						Token:         "abc",
						PhoneNumberID: "123",
						WAHost:        "https://graph.facebook.com"},
				}, nil)
				return
			},
			want: want{
				res: []entity.ClientData{
					{ClientName: "lino",
						Token:         "abc",
						PhoneNumberID: "123",
						WAHost:        "https://graph.facebook.com",
					},
					{ClientName: "thomas",
						Token:         "abc",
						PhoneNumberID: "123",
						WAHost:        "https://graph.facebook.com"},
				},
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clientDataRepoMock := test.clientDataRepoMock()
			service := NewClientService(clientDataRepoMock)

			res, err := service.GetAllClientData()
			assert.NoError(t, err)
			assert.Equal(t, test.want.res, res)
			assert.Equal(t, test.want.err, err)
			clientDataRepoMock.AssertExpectations(t)
		})
	}
}
