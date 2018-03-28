# grpc-implement

## Installation guide
1. Checkout the code with command `go get github.com/anhnguyentb/grpc-implement`
2. Go to the project directory `cd $GOPATH/src/github.com/anhnguyentb/grpc-implement/`
3. Install `go-dep` at `https://golang.github.io/dep/`
4. Run command `dep ensure` to install all of dependencies
5. Update config in `configurations.yml`
6. Run command `go run main.go` to start RPC server

## gRPC protocol buffers definition
The gRPC protocol buffers definition is in `logging/logging.proto`. You can use this file to automatic generate SDK for another platforms 

## Testing
Run command `go test ./... -v` under the directory, will get the results

```$xslt
=== RUN   TestCreateNewRecord
--- PASS: TestCreateNewRecord (0.00s)
=== RUN   TestCreateRecordShouldFailWithInCorrectData
--- PASS: TestCreateRecordShouldFailWithInCorrectData (0.00s)
=== RUN   TestFetchRecordWithServerIp
--- PASS: TestFetchRecordWithServerIp (0.01s)
=== RUN   TestFetchRecordWithClientIp
--- PASS: TestFetchRecordWithClientIp (0.01s)
=== RUN   TestFetchRecordWithTags
--- PASS: TestFetchRecordWithTags (0.01s)
=== RUN   TestFetchRecordWithAllParams
--- PASS: TestFetchRecordWithAllParams (0.01s)
=== RUN   TestFetchAllOfRecords
--- PASS: TestFetchAllOfRecords (0.02s)
PASS
ok  	github.com/anhnguyentb/grpc-implement	3.098s
=== RUN   TestConfigLoaded
--- PASS: TestConfigLoaded (0.00s)
=== RUN   TestLoadDatabase
--- PASS: TestLoadDatabase (0.01s)
=== RUN   TestCreateSchema
--- PASS: TestCreateSchema (0.01s)
=== RUN   TestLoadLogger
--- PASS: TestLoadLogger (0.00s)
PASS
ok  	github.com/anhnguyentb/grpc-implement/global	0.038s
?   	github.com/anhnguyentb/grpc-implement/logging	[no test files]
?   	github.com/anhnguyentb/grpc-implement/mocks	[no test files]
?   	github.com/anhnguyentb/grpc-implement/models	[no test files]
=== RUN   TestInitServer
--- PASS: TestInitServer (1.00s)
PASS
ok  	github.com/anhnguyentb/grpc-implement/server	1.023s
```

## Testing with benchmark
Run command `go test -bench=. -v` under the root directory, you will get the results
```$xslt
=== RUN   TestCreateNewRecord
--- PASS: TestCreateNewRecord (0.02s)
=== RUN   TestCreateRecordShouldFailWithInCorrectData
--- PASS: TestCreateRecordShouldFailWithInCorrectData (0.00s)
=== RUN   TestFetchRecordWithServerIp
--- PASS: TestFetchRecordWithServerIp (0.01s)
=== RUN   TestFetchRecordWithClientIp
--- PASS: TestFetchRecordWithClientIp (0.01s)
=== RUN   TestFetchRecordWithTags
--- PASS: TestFetchRecordWithTags (0.01s)
=== RUN   TestFetchRecordWithAllParams
--- PASS: TestFetchRecordWithAllParams (0.01s)
=== RUN   TestFetchAllOfRecords
--- PASS: TestFetchAllOfRecords (0.01s)
goos: darwin
goarch: amd64
pkg: github.com/anhnguyentb/grpc-implement
BenchmarkClientServer-8   	  200000	      5366 ns/op
PASS
ok  	github.com/anhnguyentb/grpc-implement	4.649s
``` 