package server

import (
	"testing"
	"github.com/anhnguyentb/grpc-implement/global"
	"fmt"
	"time"
)

func TestInitServer(t *testing.T) {

	global.LoadConfig()
	global.LoadLogger(true)

	errChan := make(chan error)
	go func() {
		err := InitServer()
		defer close(errChan)

		if err != nil {
			errChan <- fmt.Errorf("Fatal error start gRPC server %s \n", err)
		}
	}()

	select {
		case err := <-errChan:
			t.Fatal(err)
		case <-time.After(time.Second):
			return
	}
}
