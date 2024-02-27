package testing

import (
	"context"
	"fmt"
	"log"
	"net/http/httptest"
	"testing"
	"transmitter-artemis/config"
	"transmitter-artemis/consumer"
	"transmitter-artemis/entity"
	"transmitter-artemis/platform"
	"transmitter-artemis/provider"
	"transmitter-artemis/repository"
	"transmitter-artemis/service"

	"github.com/go-stomp/stomp/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const ConstClientName = "test-client-1"

type testSvc struct {
	suite.Suite
	ctx         context.Context
	log         provider.ILogger
	artemis     *stomp.Conn
	mongoClient *mongo.Client
	metaServer  *httptest.Server
	stop        chan struct{}
}

func (suite *testSvc) SetupSuite() {

	suite.loadConfig()

	suite.ctx = context.Background()
	suite.log = provider.NewLogger()
	suite.stop = make(chan struct{})

	suite.metaServer = MetaServer()

	suite.initArtemis()
	suite.initMongoTest()
	suite.initClientData()
	suite.initAPP()
}

func (suite *testSvc) loadConfig() {

	if err := config.LoadConfig("."); err != nil {
		log.Fatal(err)
	}

	provider.InitLogDir()
}

func (suite *testSvc) initAPP() {

	clientRepository := repository.NewClientRepository(suite.mongoClient)
	outboundRepository := repository.NewOutboundRepository(suite.mongoClient)
	drRepository := repository.NewDRRepository(suite.artemis)

	metaPlatform := platform.NewMetaClient()

	clientService := service.NewClientService(clientRepository)
	queueService := service.NewQueueService(outboundRepository, drRepository, metaPlatform, suite.log)

	clients, _ := clientService.GetAllClientData()

	go func() {
		// Membuat goroutine untuk setiap queue
		for _, clientData := range clients {
			go func(ctx context.Context, clientData entity.ClientData) {
				// Membuat listener untuk queue
				listener := consumer.NewQueueListener(suite.artemis, queueService, clientData, suite.log)
				listener.Start(ctx)
			}(suite.ctx, clientData)
		}
	}()

}

func (suite *testSvc) initMongoTest() {
	t := suite.T()

	mongoC, err := mongodb.RunContainer(
		suite.ctx,
		testcontainers.WithImage("mongo:4.4.23"),
	)
	assert.NoError(t, err)

	uri, err := mongoC.ConnectionString(suite.ctx)
	assert.NoError(t, err)

	mongoClient, err := mongo.Connect(suite.ctx, options.Client().ApplyURI(uri))
	assert.NoError(t, err)
	suite.mongoClient = mongoClient
	t.Log("SUCCESS CREATE MongoDB")
}

func (suite *testSvc) initClientData() {
	t := suite.T()

	docs := []interface{}{
		entity.ClientData{
			ClientName:    ConstClientName,
			Token:         "valid_token",
			PhoneNumberID: "valid_phone_number_id",
			WAHost:        suite.metaServer.URL,
		},
		entity.ClientData{
			ClientName:    "client-2",
			Token:         "invalid_token",
			PhoneNumberID: "valid_phone_number_id",
			WAHost:        suite.metaServer.URL,
		},
	}

	collName := "client-info"
	coll := suite.mongoClient.Database(config.Configuration.MongoDB.Database).Collection(collName)
	result, err := coll.InsertMany(context.TODO(), docs)

	assert.NoError(t, err)

	fmt.Printf("ClientData inserted: %v\n", len(result.InsertedIDs))
	for _, id := range result.InsertedIDs {
		fmt.Printf("Inserted CLientData with _id: %v\n", id)
	}
}

func (suite *testSvc) initArtemis() {
	ctx := context.Background()

	t := suite.T()

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "docker.io/apache/activemq-artemis:2.30.0-alpine",
			Env: map[string]string{
				"ARTEMIS_USER":     "artemis",
				"ARTEMIS_PASSWORD": "artemis",
				"AMQ_USER":         "artemis",
				"AMQ_PASSWORD":     "artemis",
			},
			ExposedPorts: []string{"61616/tcp", "8161/tcp"},
			WaitingFor: wait.ForAll(
				wait.ForLog("Server is now live"),
				wait.ForLog("REST API available"),
			),
		},
		Started: true,
	}

	artemisContainer, err := testcontainers.GenericContainer(ctx, req)
	require.NoError(t, err)

	host, err := artemisContainer.Host(ctx)
	if err != nil {
		require.NoError(t, err)
	}

	port, err := artemisContainer.MappedPort(suite.ctx, "61616")
	assert.NoError(t, err)

	host = fmt.Sprintf("%s:%s", host, port.Port())

	conn, err := stomp.Dial("tcp", host, stomp.ConnOpt.Login("artemis", "artemis"))
	if err != nil {
		require.NoError(t, err)
	}

	suite.artemis = conn

	t.Log("SUCCESS CREATE ARTEMIS")
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(testSvc))
}
