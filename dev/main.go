package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"transmitter-artemis/config"
	"transmitter-artemis/consumer"
	"transmitter-artemis/entity"
	"transmitter-artemis/platform"
	"transmitter-artemis/provider"
	"transmitter-artemis/repository"
	"transmitter-artemis/service"
)

func init() {
	if err := config.LoadConfig("."); err != nil {
		log.Fatal(err)
	}

	provider.InitLogDir()
}

func main() {
	logger := provider.NewLogger()

	mongodb, err := provider.NewMongoDBClient()

	if err != nil {
		logger.Errorf(provider.AppLog, "Error Connect DB")
	}
	logger.Infof(provider.AppLog, "Success Connect DB")

	// Create conn Artemis
	artemisConn, err := provider.NewArtemis()

	if err != nil {
		logger.Errorf(provider.AppLog, "Error Connect to Artemis")
	}
	logger.Infof(provider.AppLog, "Success Connect to Artemis")

	// Create instance repository and service
	clientRepository := repository.NewClientRepository(mongodb)
	outboundRepository := repository.NewOutboundRepository(mongodb)
	drRepository := repository.NewDRRepository(artemisConn)
	metaPlatform := platform.NewMetaClient()

	clientService := service.NewClientService(clientRepository)
	queueService := service.NewQueueService(outboundRepository, drRepository, metaPlatform, logger)

	// Get All Client Data
	clients, err := clientService.GetAllClientData()

	if err != nil {
		logger.Errorf(provider.AppLog, "Error Get Client data")
	}
	logger.Infof(provider.AppLog, "Success Get Client Data")

	// Membuat goroutine untuk setiap queue
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {

		for _, clientData := range clients {
			go func(ctx context.Context, clientData entity.ClientData) {
				// Membuat listener untuk queue
				listener := consumer.NewQueueListener(artemisConn, queueService, clientData, logger)
				listener.Start(ctx)
			}(ctx, clientData)
		}
	}()

	// Tunggu sampai aplikasi diberhentikan
	waitForShutdown(logger)
}

func waitForShutdown(logger provider.ILogger) {
	// Membuat channel untuk menerima sinyal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Menunggu sinyal untuk keluar
	<-stop

	logger.Infof(provider.AppLog, "Received signal to stop. Exiting...")
}
