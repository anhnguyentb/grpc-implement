package main

import (
	"testing"
	"github.com/anhnguyentb/grpc-implement/logging"
	"google.golang.org/grpc"
	"github.com/spf13/viper"
	"time"
	"golang.org/x/net/context"
	"github.com/anhnguyentb/grpc-implement/global"
	"github.com/anhnguyentb/grpc-implement/server"
	"github.com/anhnguyentb/grpc-implement/models"
)

func init() {

	createServer()
}

func createServer() {

	connChan := make(chan bool)
	go func() {
		global.LoadConfig()
		global.LoadLogger(true)

		select {
		case <-time.After(3 * time.Second):
			connChan <- true
			close(connChan)
		}
		server.InitServer()
	}()

	<- connChan
}

func populateTestData(t testing.TB, clearOnly bool) {

	//Delete all data with "testing" tags
	db, err := global.GetConnection()
	if err != nil {
		t.Fatalf("Cannot connect to database: %v", err)
	}

	_, err = db.Model(&models.Audit{}).Where("tags::jsonb ?& array['testing']").Delete()
	if err != nil {
		t.Fatalf("Cannot delete all of testing records: %v", err)
		return
	}

	if clearOnly {
		return
	}

	//Testing Data
	testData := []models.Audit{
		models.Audit{
			ServerIp: "187.2.169.212",
			Message: "Contrary to popular belief, Lorem Ipsum is not simply random text",
			Tags: []string{"testing", "brazil", "Contrary", "tag_a"},
		},
		models.Audit{
			ClientIp: "177.142.89.166",
			Message: "a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words",
			Tags: []string{"testing", "Contrary", "tag_a", "tag_b"},
		},
		models.Audit{
			ServerIp: "189.122.155.188",
			ClientIp: "184.1.112.179",
			Message: "The first line of Lorem Ipsum, \"Lorem ipsum dolor sit amet..\", comes from a line in section 1.10.32.",
			Tags: []string{"testing", "tag_a"},
		},
	}

	for _, data := range testData {
		err := db.Insert(&data)
		if err != nil {
			t.Fatalf("Failed to populate dummy data %v", err)
			return
		}
	}
}

func getConnectServer(t testing.TB) *grpc.ClientConn {

	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost" + viper.GetString("server.port"), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	return conn
}

func TestCreateNewRecord(t *testing.T) {
	conn := getConnectServer(t)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	c := logging.NewLoggingClient(conn)
	res, err := c.Create(ctx, &logging.LoggingRequest{
		ServerIp: "192.168.0.1",
		ClientIp: "192.168.1.1",
		Message: "Message from testing",
		Tags: []string{"testing", "go"},
	})
	if err != nil {
		t.Fatalf("Can not call Create method: %v", err)
	}

	if !res.Status {
		t.Error("Failed to call Create method with correct data")
	}
}

func TestCreateRecordShouldFailWithInCorrectData(t *testing.T) {
	conn := getConnectServer(t)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	c := logging.NewLoggingClient(conn)
	res, err := c.Create(ctx, &logging.LoggingRequest{
		ServerIp: "192.168.0.1",
		ClientIp: "192.168.1.1",
		Tags: []string{"testing", "go"},
	})
	if err != nil {
		t.Fatalf("Can not call Create method: %v", err)
	}

	if res.Status {
		t.Error("Failed to validate Create method request data")
	}

	if len(res.Errors) == 0 {
		t.Error("Create method should return errors with invalid request data")
	}
}

func TestFetchRecordWithServerIp(t *testing.T) {
	conn := getConnectServer(t)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	//Populate dummy data
	populateTestData(t, false)

	c := logging.NewLoggingClient(conn)
	res, err := c.Fetch(ctx, &logging.QueryRequest{
		ServerIp: "187.2.169.212",
	})
	if err != nil {
		t.Fatalf("Can not call Fetch method: %v", err)
	}

	if !res.Status {
		t.Errorf("Fail to fetch with server_ip, expected true but got %t", res.Status)
	}

	if len(res.Results) != 1 {
		t.Errorf("Fail to fetch with server_ip, expected 1 record but got %d", len(res.Results))
	}
}

func TestFetchRecordWithClientIp(t *testing.T) {
	conn := getConnectServer(t)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	//Populate dummy data
	populateTestData(t, false)

	c := logging.NewLoggingClient(conn)
	res, err := c.Fetch(ctx, &logging.QueryRequest{
		ClientIp: "177.142.89.166",
	})
	if err != nil {
		t.Fatalf("Can not call Fetch method: %v", err)
	}

	if !res.Status {
		t.Errorf("Fail to fetch with client_ip, expected true but got %t", res.Status)
	}

	if len(res.Results) != 1 {
		t.Errorf("Fail to fetch with client_ip, expected 1 record but got %d", len(res.Results))
	}
}

func TestFetchRecordWithTags(t *testing.T) {
	conn := getConnectServer(t)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	//Populate dummy data
	populateTestData(t, false)

	c := logging.NewLoggingClient(conn)
	res, err := c.Fetch(ctx, &logging.QueryRequest{
		Tags: []string{"testing", "tag_a"},
	})
	if err != nil {
		t.Fatalf("Can not call Fetch method: %v", err)
	}

	if !res.Status {
		t.Errorf("Fail to fetch with tags, expected true but got %t", res.Status)
	}

	if len(res.Results) != 3 {
		t.Errorf("Fail to fetch with tags, expected 3 record but got %d", len(res.Results))
	}
}

func TestFetchRecordWithAllParams(t *testing.T) {
	conn := getConnectServer(t)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	//Populate dummy data
	populateTestData(t, false)

	c := logging.NewLoggingClient(conn)
	res, err := c.Fetch(ctx, &logging.QueryRequest{
		ClientIp: "184.1.112.179",
		ServerIp: "189.122.155.188",
		Tags: []string{"testing", "tag_a"},
	})
	if err != nil {
		t.Fatalf("Can not call Fetch method: %v", err)
	}

	if !res.Status {
		t.Errorf("Fail to fetch with all of parameters, expected true but got %t", res.Status)
	}

	if len(res.Results) != 1 {
		t.Errorf("Fail to fetch with all of parameters, expected 1 record but got %d", len(res.Results))
	}
}

func TestFetchAllOfRecords(t *testing.T) {
	conn := getConnectServer(t)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	//Populate dummy data
	populateTestData(t, false)

	c := logging.NewLoggingClient(conn)
	res, err := c.Fetch(ctx, &logging.QueryRequest{})
	if err != nil {
		t.Fatalf("Can not call Fetch method: %v", err)
	}

	if !res.Status {
		t.Errorf("Fail to fetch all of records, expected true but got %t", res.Status)
	}

	if len(res.Results) < 3 {
		t.Errorf("Fail to fetch all of records, expected greater than 3 record but got %d", len(res.Results))
	}
}
