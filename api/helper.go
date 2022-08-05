package api

import (
	api2 "therealbroker/api/proto/src/broker/api/proto"
	"therealbroker/pkg/broker"
)

func CreatePublishResponse(id int) *api2.PublishResponse {
	return &api2.PublishResponse{Id: int32(id)}
}

func CreateMessageResponse(body broker.Message) *api2.MessageResponse {
	return &api2.MessageResponse{Body: []byte(body.Body)}
}
