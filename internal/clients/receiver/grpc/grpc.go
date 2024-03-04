package grpc

import (
	"context"
	"fmt"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"grpc/internal/config"
	"grpc/internal/services/anomaly"
	"grpc/internal/storage"
	"grpc/pkg/logger/sl"
	tr "grpc/protos/gen/go/transmitter"
	"log/slog"
	"os"
	"sync"
	"time"
)

type Client struct {
	Api         tr.TransmitterClient
	Log         *slog.Logger
	Coefficient float64
	DB          *storage.DataBase
}

func New(cc *grpc.ClientConn, log *slog.Logger, cfg *config.Config) (*Client, error) {
	const op = "clients.receiver.grpc.New"

	db, err := storage.New(cfg.Postgres)
	if err != nil {
		log.Error("error occurred while connecting to database", op, sl.Err(err))
		return nil, err
	}

	return &Client{
		Api:         tr.NewTransmitterClient(cc),
		Log:         log,
		Coefficient: cfg.Coefficient,
		DB:          db,
	}, nil
}

func NewConnection(ctx context.Context, log *slog.Logger, addr string, retriesCount int, timeout time.Duration) (*grpc.ClientConn, error) {
	const op = "clients.receiver.grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cc, nil
}

func InterceptorLogger(log *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		log.Log(ctx, slog.Level(level), msg, fields...)
	})
}

func (c *Client) Receiver() {
	const op = "clients.receiver.grpc.Receiver"
	stream, err := c.Api.Transmit(context.Background(), &emptypb.Empty{})
	if err != nil {
		c.Log.Error("failed to start client", op, sl.Err(err))
		os.Exit(1)
	}

	pool := sync.Pool{
		New: func() interface{} {
			return &tr.TransmitResponse{}
		},
	}
	a := anomaly.New()
	var anomalyStage bool
	ctx := stream.Context()

	for {
		select {
		case <-ctx.Done():
			c.Log.Error(ctx.Err().Error())
			return
		default:
			resp := pool.Get().(*tr.TransmitResponse)
			resp, err = stream.Recv()
			if err != nil {
				c.Log.Error("failed to get request from client", op, err)
				return
			}
			c.Log.Info("received data",
				slog.String("id", resp.SessionId),
				slog.String("frequency", fmt.Sprintf("%.6f", resp.Frequency)),
				slog.String("timestamp", resp.Timestamp.AsTime().UTC().String()),
			)

			switch anomalyStage {
			case true:
				if a.IsAnomaly(resp.Frequency, c.Coefficient) {
					c.Log.Info("new anomaly detected",
						slog.String("value", fmt.Sprintf("%.6f", resp.Frequency)))
					c.DB.WriteAnomaly(
						resp.SessionId,
						resp.Frequency,
						resp.Timestamp.AsTime().UTC().String(),
					)
				}
			default:
				anomalyStage = a.UpdateMeanSTD(resp.Frequency)
				c.Log.Info("values processed",
					slog.String("count", fmt.Sprintf("%d", int(a.Count))),
					slog.String("predicted value of mean", fmt.Sprintf("%.6f", a.Mean)),
					slog.String("predicted value of std", fmt.Sprintf("%.6f", a.STD)),
				)
				if anomalyStage {
					c.Log.Info("starting Anomaly Detection stage")
				}
			}
			resp.Reset()
			pool.Put(resp)
		}
	}
}
