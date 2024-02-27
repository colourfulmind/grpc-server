package grpcserver

import (
	"cmd/grpc/main.go/internal/services"
	tr "cmd/grpc/main.go/protos/gen/go/transmitter"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Transmitter interface {
	Transmit(ctx context.Context, req *tr.TransmitRequest) (*tr.TransmitResponse, error)
}

type Server struct {
	tr.UnimplementedTransmitterServer
	tr Transmitter
}

func Register(gRPC *grpc.Server) {
	tr.RegisterTransmitterServer(gRPC, &Server{})
}

func (s *Server) Transmit(ctx context.Context, req *tr.TransmitRequest) (*tr.TransmitResponse, error) {
	d := services.New()
	sessionID, frequency := d.GenerateData(ctx)
	return &tr.TransmitResponse{
		SessionId: sessionID,
		Frequency: frequency,
		Timestamp: timestamppb.Now(),
	}, nil
}
