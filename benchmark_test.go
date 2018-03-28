package main

import (
	"testing"
	"time"
	"github.com/anhnguyentb/grpc-implement/logging"
	"golang.org/x/net/context"
)

func BenchmarkClientServer(t *testing.B) {

	for i := 0; i < t.N; i++ {

		go func(idx int) {

			conn := getConnectServer(t)
			defer conn.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			c := logging.NewLoggingClient(conn)
			switch {
			case idx%3 == 0:
				res, err := c.Create(ctx, &logging.LoggingRequest{
					ServerIp: "192.168.0.1",
					ClientIp: "192.168.1.1",
					Message:  "Message from testing",
					Tags:     []string{"testing", "go"},
				})
				if err != nil {
					t.Fatalf("Can not call Create method: %v", err)
				}

				if !res.Status {
					t.Error("Failed to call Create method with correct data")
				}
			default:
				res, err := c.Fetch(ctx, &logging.QueryRequest{})
				if err != nil {
					t.Fatalf("Can not call Fetch method: %v", err)
				}

				if !res.Status {
					t.Errorf("Fail to fetch with tags, expected true but got %t", res.Status)
				}

				if len(res.Results) == 0 {
					t.Errorf("Fail to fetch with tags, expected more than 0 record but got %d", len(res.Results))
				}

			}
		}(i)
	}
}
