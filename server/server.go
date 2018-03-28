package server

import (
	"context"
	pb "github.com/anhnguyentb/grpc-implement/logging"
	"google.golang.org/grpc"
	"net"
	"github.com/spf13/viper"
	"github.com/anhnguyentb/grpc-implement/global"
	"go.uber.org/zap"
	"fmt"
	"github.com/anhnguyentb/grpc-implement/models"
	"time"
	"strconv"
	"strings"
)

type loggingServer struct {}

func (l *loggingServer) Create(ctx context.Context, data *pb.LoggingRequest) (*pb.LoggingResponse, error) {

	global.Log.Infow(
		"Client requested to Create method",
		zap.String("server_ip", data.GetServerIp()),
		zap.String("client_ip", data.GetClientIp()),
		zap.String("message", data.GetMessage()),
		zap.Any("tags", data.GetTags()),
	)

	response := &pb.LoggingResponse{Status: false}
	if len(data.GetMessage()) == 0 {

		response.Errors = []string{"Message is required"}
		return response, nil
	}

	auditRecord := models.Audit{
		ClientIp: data.GetClientIp(),
		ServerIp: data.GetServerIp(),
		Message:  data.GetMessage(),
		Tags:     data.GetTags(),
	}
	err := global.Db.Insert(&auditRecord)
	if err != nil {
		global.Log.Errorw(
			"Failed to create new audit record",
			zap.String("server_ip", data.GetServerIp()),
			zap.String("client_ip", data.GetClientIp()),
			zap.String("message", data.GetMessage()),
			zap.Any("tags", data.GetTags()),
		)
		return response, nil
	}

	global.Log.Infow(
		"Created new audit record",
		zap.Int64("id", auditRecord.Id),
	)

	response.Status = true
	response.Message = "New audit record has been added with ID " + strconv.FormatInt(auditRecord.Id, 10)

	return response, nil
}

func (l *loggingServer) Fetch(ctx context.Context, data *pb.QueryRequest) (*pb.QueryResponse, error) {

	global.Log.Infow(
		"Client requested to Fetch method",
		zap.String("server_ip", data.GetServerIp()),
		zap.String("client_ip", data.GetClientIp()),
		zap.Any("tags", data.GetTags()),
	)

	var auditRecords []models.Audit
	response := &pb.QueryResponse{Status: false}

	query := global.Db.Model(&auditRecords)

	//Check if server_ip has been defined
	if len(data.GetServerIp()) > 0 {
		query.Where("server_ip = ?", data.GetServerIp())
	}

	//Check if client_ip has been defined
	if len(data.GetClientIp()) > 0 {
		query.Where("client_ip = ?", data.GetClientIp())
	}

	//Check if tags has been defined
	if len(data.GetTags()) > 0 {

		stmt := "tags::jsonb ?& array["
		for _, value := range data.GetTags() {
			stmt += "'" + value + "',"
		}
		stmt = strings.TrimRight(stmt, ",") + "]"
		query.Where(stmt)
	}

	err := query.Select()
	if err != nil {

		global.Log.Errorw(
			"Failed to fetch audit record",
			zap.String("server_ip", data.GetServerIp()),
			zap.String("client_ip", data.GetClientIp()),
			zap.Any("tags", data.GetTags()),
		)
		response.Errors = []string{err.Error()}
		return response, nil
	}

	global.Log.Infow(
		fmt.Sprintf("Fetched %d audit record", len(auditRecords)),
		zap.String("server_ip", data.GetServerIp()),
		zap.String("client_ip", data.GetClientIp()),
		zap.Any("tags", data.GetTags()),
	)
	var auditResponses []*pb.AuditRecord
	for _, record := range auditRecords {
		auditResponses = append(auditResponses, &pb.AuditRecord{
			ClientIp: record.ClientIp,
			ServerIp: record.ServerIp,
			Message: record.Message,
			Tags: record.Tags,
		})
	}

	response.Status = true
	response.Results = auditResponses

	return response, nil
}

func InitServer() error {

	port := viper.GetString("server.port")
	if !viper.IsSet("server.port") || port == "" {
		err := fmt.Errorf("Port is not defined")
		global.Log.Error(
			"Port is not defined",
			zap.Error(err),
		)
		return err
	}

	global.Log.Infow(
		"Server starting ...",
		zap.String("port", port),
	)
	listen, err := net.Listen("tcp", port)
	if err != nil {

		global.Log.Errorw(
			"Failed to listen tcp",
			zap.String("port", port),
			zap.Error(err),
		)
		return err
	}

	//Load database
	global.LoadDatabase()
	//defer global.Db.Close()

	//Create default schema if it's not exists
	global.CreateSchema()

	server := grpc.NewServer()
	pb.RegisterLoggingServer(server, &loggingServer{})
	go func() {
		select {
			case <-time.After(time.Second):
				global.Log.Infof("Server is running under port: %s", port)
		}
	}()
	if err := server.Serve(listen); err != nil {

		global.Log.Errorw(
			"Failed to serve from gRPC server",
			zap.Error(err),
		)
		return err
	}

	return nil
}