package api

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"

	"therealbroker/api/proto/src/broker/api/proto"
	module "therealbroker/internal/broker"
	"therealbroker/internal/config"
	"therealbroker/internal/metrics"
	"therealbroker/pkg/broker"
)

type Handler struct {
	broker broker.Broker
	conf   *config.Config
}

func New(broker broker.Broker, conf *config.Config) *Handler {
	return &Handler{
		broker: broker,
		conf:   conf,
	}
}

func (h *Handler) Publish(ctx context.Context, request *proto.PublishRequest) (*proto.PublishResponse, error) {
	response, err := h.broker.Publish(ctx, request.Subject, module.CreateBrokerMessage(request.Body, request.ExpirationSeconds))
	if err != nil {
		return nil, err
	}

	return &proto.PublishResponse{Id: int32(response)}, nil
}

func (h *Handler) Subscribe(request *proto.SubscribeRequest, stream proto.Broker_SubscribeServer) error {
	streamChannel, err := h.broker.Subscribe(stream.Context(), request.Subject)
	if err != nil {
		return err
	}

	metrics.ActiveSubscribers.Inc()
	defer metrics.ActiveSubscribers.Dec()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case body := <-streamChannel:
				if err = stream.Send(&proto.MessageResponse{Body: []byte(body.Body)}); err != nil {
					logrus.WithError(err).Error("failed to send message to stream")
					return
				}
			case <-stream.Context().Done():
				logrus.Info("stream context done")
				return
			}
		}
	}()
	wg.Wait()

	return nil
}

func (h *Handler) Fetch(ctx context.Context, request *proto.FetchRequest) (*proto.MessageResponse, error) {
	response, err := h.broker.Fetch(ctx, request.Subject, int(request.Id))
	if err != nil {
		return nil, err
	}

	return &proto.MessageResponse{Body: []byte(response.Body)}, nil
}
