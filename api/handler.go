package api

import (
	"context"
	"google.golang.org/grpc"
	"net"
	"sync"
	api "therealbroker/api/proto/src/broker/api/proto"
	module "therealbroker/internal/broker"
	"therealbroker/internal/prometheus"
	"therealbroker/pkg/broker"
	"therealbroker/pkg/config"
	"therealbroker/pkg/log"
	"time"
)

type Handler struct {
	broker  broker.Broker
	log     *log.Logger
	conf    *config.Config
	metrics *prometheus.APIMetrics
}

func New(broker broker.Broker, log *log.Logger, conf *config.Config, metrics *prometheus.APIMetrics) *Handler {
	server := grpc.NewServer()

	handler := Handler{
		broker:  broker,
		log:     log,
		conf:    conf,
		metrics: metrics,
	}

	api.RegisterBrokerServer(server, &handler)

	lis, err := net.Listen("tcp", conf.Grpc.Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Info("Start listening on address: ", conf.Grpc.Address)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return &handler
}

func (h *Handler) Publish(ctx context.Context, request *api.PublishRequest) (*api.PublishResponse, error) {
	status := "success"
	method := "publish"
	start := time.Now()
	response, err := h.broker.Publish(ctx, request.Subject, module.CreateBrokerMessage(request.Body, request.ExpirationSeconds))
	if err != nil {
		h.log.Errorln("broker.Publish", err)
		status = "fail"
	}
	duration := time.Since(start).Nanoseconds()
	h.metrics.MethodDuration.WithLabelValues(method, status).Observe(float64(duration))
	h.metrics.MethodCount.WithLabelValues(method, status).Inc()
	return CreatePublishResponse(response), err
}

func (h *Handler) Subscribe(request *api.SubscribeRequest, stream api.Broker_SubscribeServer) error {
	streamChannel, err := h.broker.Subscribe(stream.Context(), request.Subject)
	if err != nil {
		h.log.Errorln("broker.Subscribe", err)
		return err
	}

	var wg sync.WaitGroup
	h.metrics.ActiveSubscribers.Inc()
	wg.Add(1)
	go func() {
		for {
			select {
			case body := <-streamChannel:
				err = stream.Send(CreateMessageResponse(body))
				if err != nil {
					h.log.Errorln("stream.Send", err)
				}
			case <-stream.Context().Done():
				err = stream.Context().Err()
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()
	h.metrics.ActiveSubscribers.Dec()
	return err
}

func (h *Handler) Fetch(ctx context.Context, request *api.FetchRequest) (*api.MessageResponse, error) {
	status := "success"
	method := "fetch"
	start := time.Now()
	response, err := h.broker.Fetch(ctx, request.Subject, int(request.Id))
	if err != nil {
		h.log.Errorln("broker.Fetch", err)
		status = "fail"
	}
	duration := time.Since(start).Nanoseconds()
	h.metrics.MethodDuration.WithLabelValues(method, status).Observe(float64(duration))
	h.metrics.MethodCount.WithLabelValues(method, status).Inc()
	return CreateMessageResponse(response), err
}
