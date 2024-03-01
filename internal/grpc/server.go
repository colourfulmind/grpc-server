package grpcserver

import (
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"

	//"google.golang.org/protobuf/types/known/timestamppb"
	"grpc/internal/services"
	tr "grpc/protos/gen/go/transmitter"
)

type Server struct {
	tr.UnimplementedTransmitterServer
}

func Register(gRPC *grpc.Server) {
	tr.RegisterTransmitterServer(gRPC, &Server{})
}

func (s *Server) Transmit(_ *emptypb.Empty, stream tr.Transmitter_TransmitServer) error {
	d := services.New()
	ctx := stream.Context()

	log.Printf(
		"Starting stream with ID: %s\tMean: %f\tSTD: %f\n",
		d.SessionID,
		d.Mean,
		d.STD,
	)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			frequency := d.GenerateFrequency()
			msg := &tr.TransmitResponse{
				SessionId: d.SessionID,
				Frequency: frequency,
				Timestamp: timestamppb.Now(),
			}
			log.Printf(
				"Recieve data ID: %s\tFrequency: %f\tTimestamp: %s\n",
				msg.SessionId,
				msg.Frequency,
				msg.Timestamp.AsTime().UTC(),
			)
			if err := stream.SendMsg(msg); err != nil {
				return err
			}
			time.Sleep(3 * time.Second)
		}
	}
}
