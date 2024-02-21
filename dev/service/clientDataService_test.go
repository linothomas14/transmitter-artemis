package service

import (
	"testing"
	"transmitter-artemis/entity"
	mocks "transmitter-artemis/mocks/repository"

	"github.com/stretchr/testify/assert"
)

func TestGetAllClientData(t *testing.T) {

	mockClientDataRepo := mocks.NewClientRepository(t)
	app := NewClientService(mockClientDataRepo)

	t.Run("Test GetAllClientData", func(t *testing.T) {

		expectedClientDataResp := []entity.ClientData{
			{ClientName: "lino",
				Token:         "abc",
				PhoneNumberID: "123",
				WAHost:        "https://graph.facebook.com",
			},
			{ClientName: "thomas",
				Token:         "abc",
				PhoneNumberID: "123",
				WAHost:        "https://graph.facebook.com"},
		}

		mockClientDataRepo.On("GetAllClientData").Return(expectedClientDataResp, nil)

		actualClientDataResp, err := app.GetAllClientData()

		assert.NoError(t, err)
		assert.NotEmpty(t, actualClientDataResp)
		assert.Equal(t, expectedClientDataResp, actualClientDataResp)
	})
}
