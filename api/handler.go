package api

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	api "therealbroker/api/proto/src/broker/api/proto"
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

func (h *Handler) Publish(ctx context.Context, request *api.PublishRequest) (*api.PublishResponse, error) {
	status := "success"
	method := "publish"
	start := time.Now()
	response, err := h.broker.Publish(ctx, request.Subject, module.CreateBrokerMessage(request.Body, request.ExpirationSeconds))
	if err != nil {
		logrus.Errorln("broker.Publish", err)
		status = "fail"
	}
	duration := time.Since(start).Nanoseconds()
	metrics.MethodDuration.WithLabelValues(method, status).Observe(float64(duration))
	metrics.MethodCount.WithLabelValues(method, status).Inc()
	return CreatePublishResponse(response), err
}

func (h *Handler) Subscribe(request *api.SubscribeRequest, stream api.Broker_SubscribeServer) error {
	streamChannel, err := h.broker.Subscribe(stream.Context(), request.Subject)
	if err != nil {
		logrus.Errorln("broker.Subscribe", err)
		return err
	}

	var wg sync.WaitGroup
	metrics.ActiveSubscribers.Inc()
	wg.Add(1)
	go func() {
		for {
			select {
			case body := <-streamChannel:
				err = stream.Send(CreateMessageResponse(body))
				if err != nil {
					logrus.Errorln("stream.Send", err)
				}
			case <-stream.Context().Done():
				err = stream.Context().Err()
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()
	metrics.ActiveSubscribers.Dec()
	return err
}

func (h *Handler) Fetch(ctx context.Context, request *api.FetchRequest) (*api.MessageResponse, error) {
	status := "success"
	method := "fetch"
	start := time.Now()
	response, err := h.broker.Fetch(ctx, request.Subject, int(request.Id))
	if err != nil {
		logrus.Errorln("broker.Fetch", err)
		status = "fail"
	}
	duration := time.Since(start).Nanoseconds()
	metrics.MethodDuration.WithLabelValues(method, status).Observe(float64(duration))
	metrics.MethodCount.WithLabelValues(method, status).Inc()
	return CreateMessageResponse(response), err
}
