package endpoints

import (
	"context"
	//"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"jaxf-github.fanatics.corp/apparel/partner-service/pkg/service"
)

func New(svc service.PartnerService, logger log.Logger) Endpoints {
	var keyValueEndpoint endpoint.Endpoint
	{
		keyValueEndpoint = MakeKeyValueEndpoint(svc)
		keyValueEndpoint = LoggingMiddleware(log.With(logger, "method", "Get data by Key/Value"))(keyValueEndpoint)
	}

	var getDataByIdEndpoint endpoint.Endpoint
	{
		getDataByIdEndpoint = MakeGetDataByIdEndpoint(svc)
		getDataByIdEndpoint = LoggingMiddleware(log.With(logger, "method", "Get Data By Id"))(getDataByIdEndpoint)
	}

	return Endpoints{
		KeyValueEndpoint:    keyValueEndpoint,
		GetDataByIdEndpoint: getDataByIdEndpoint,
	}
}

type Endpoints struct {
	KeyValueEndpoint    endpoint.Endpoint
	GetDataByIdEndpoint endpoint.Endpoint
}

// MakeSayHelloEndpoint returns an endpoint that invokes SayHello on the service.
// Primarily useful in a server.

func MakeKeyValueEndpoint(service service.PartnerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		keyValueReq := request.(KeyValueRequest)
		partnerIdReply, partnerCodeReply, attributes, err := service.GetPartnerDataByKeyValue(ctx, keyValueReq.Key, keyValueReq.Value, keyValueReq.Group)

		return PartnerDataReply{
			PartnerId:   partnerIdReply,
			PartnerCode: partnerCodeReply,
			Attributes:  attributes,
			Error:       err2str(err),
		}, nil
	}
}

// These types are unexported because they only exist to serve the endpoint
// domain, which is totally encapsulated in this package. They are otherwise
// opaque to all callers.

func MakeGetDataByIdEndpoint(service service.PartnerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		getDataByIdReq := request.(IdRequest)
		partnerIdReply, partnerCodeReply, attributes, err := service.GetDataById(ctx, getDataByIdReq.PartnerId, getDataByIdReq.PartnerCode, getDataByIdReq.Group)

		return PartnerDataReply{
			PartnerId:   partnerIdReply,
			PartnerCode: partnerCodeReply,
			Attributes:  attributes,
			Error:       err2str(err),
		}, nil
	}
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

type KeyValueRequest struct {
	Key   string
	Value string
	Group string
}

type IdRequest struct {
	PartnerId   int32
	PartnerCode string
	Group       string
}

type PartnerDataReply struct {
	PartnerId   int32
	PartnerCode string
	Attributes  map[string]string
	Error       string
}
