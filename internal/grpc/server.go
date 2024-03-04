package grpcserver

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"grpc/internal/services/frequency"
	tr "grpc/protos/gen/go/transmitter"
	"log/slog"
	"time"
)

type Server struct {
	tr.UnimplementedTransmitterServer
	log *slog.Logger
}

func Register(gRPC *grpc.Server, log *slog.Logger) {
	tr.RegisterTransmitterServer(gRPC, &Server{log: log})
}

func (s *Server) Transmit(_ *emptypb.Empty, stream tr.Transmitter_TransmitServer) error {
	f := frequency.New()
	ctx := stream.Context()

	s.log.Info("starting stream",
		slog.String("id", f.SessionID),
		slog.String("mean", fmt.Sprintf("%.6f", f.Mean)),
		slog.String("std", fmt.Sprintf("%.6f", f.STD)))

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			freq := f.GenerateFrequency()
			msg := &tr.TransmitResponse{
				SessionId: f.SessionID,
				Frequency: freq,
				Timestamp: timestamppb.Now(),
			}
			s.log.Info("received data",
				slog.String("id", msg.SessionId),
				slog.String("frequency", fmt.Sprintf("%.6f", msg.Frequency)),
				slog.String("timestamp", msg.Timestamp.AsTime().UTC().String()),
			)
			if err := stream.SendMsg(msg); err != nil {
				return err
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}
}
