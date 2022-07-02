package main

import (
	"gRPC-Chat/pkg/api"
	"gRPC-Chat/pkg/service"
	"google.golang.org/grpc"
	glog "google.golang.org/grpc/grpclog"
	"log"
	"net"
	"os"
)

var grpcLog glog.LoggerV2

func init() {
	grpcLog = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

func main() {
	var connections []*service.Connection

	server := &service.Server{Connection: connections}

	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":9180")
	if err != nil {
		log.Fatalf("error creating the server %v", err)
	}

	grpcLog.Info("Starting server at port :9180")

	api.RegisterBroadcastServer(grpcServer, server)
	grpcServer.Serve(listener)
}
