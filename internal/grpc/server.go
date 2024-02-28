package grpcserver

import (
	"cmd/grpc/main.go/internal/services"
	tr "cmd/grpc/main.go/protos/gen/go/transmitter"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	tr.UnimplementedTransmitterServer
}

func Register(gRPC *grpc.Server) {
	tr.RegisterTransmitterServer(gRPC, &Server{})
}

func (s *Server) Transmit(ctx context.Context, _ *emptypb.Empty) (*tr.TransmitResponse, error) {
	d := services.New()
	sessionID, frequency := d.GenerateData()
	return &tr.TransmitResponse{
		SessionId: sessionID,
		Frequency: frequency,
		Timestamp: timestamppb.Now(),
	}, nil
}
