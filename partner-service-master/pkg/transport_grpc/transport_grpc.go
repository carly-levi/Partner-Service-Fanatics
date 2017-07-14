package transport_grpc

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/pkg/errors"
	// oldcontext is necessary because transport_grpc still uses the experimental context rather than stdlib context
	oldcontext "golang.org/x/net/context"
	"jaxf-github.fanatics.corp/apparel/partner-service/pkg/endpoints"
	"jaxf-github.fanatics.corp/apparel/partner-service/pkg/pb"
)

func MakeGRPCServer(endpoints endpoints.Endpoints, logger log.Logger) pb.PartnerServiceServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}

	return &grpcServer{
		keyValue: grpctransport.NewServer(
			endpoints.KeyValueEndpoint,
			DecodeGRPCKeyValueRequest,
			EncodeGRPCResponse,
			options...,
		),
		dataById: grpctransport.NewServer(
			endpoints.GetDataByIdEndpoint,
			DecodeGRPCDataByIdRequest,
			EncodeGRPCResponse,
			options...,
		),
	}
}

type grpcServer struct {
	keyValue grpctransport.Handler
	dataById grpctransport.Handler
}

func (s *grpcServer) GetPartnerDataByKeyValue(ctx oldcontext.Context, req *pb.KeyValueRequest) (*pb.PartnerDataReply, error) {

	_, rep, err := s.keyValue.ServeGRPC(ctx, req)

	if err != nil {

		err = errors.Wrap(err, fmt.Sprintf("Error serving transport_grpc in KeyValue"))
		return nil, err
	}
	return rep.(*pb.PartnerDataReply), nil
}

func (s *grpcServer) GetDataById(ctx oldcontext.Context, req *pb.IdRequest) (*pb.PartnerDataReply, error) {
	_, rep, err := s.dataById.ServeGRPC(ctx, req)
	if err != nil {
		err = errors.Wrap(err, "error serving transport_grpc in DataById")
		return nil, err
	}
	return rep.(*pb.PartnerDataReply), nil
}

func DecodeGRPCKeyValueRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.KeyValueRequest)

	return endpoints.KeyValueRequest{Key: req.Key, Value: req.Value, Group: req.Group}, nil
}

func DecodeGRPCDataByIdRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.IdRequest)
	return endpoints.IdRequest{PartnerId: req.PartnerId, PartnerCode: req.PartnerCode, Group: req.Group}, nil
}

func EncodeGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoints.PartnerDataReply)
	return &pb.PartnerDataReply{PartnerId: resp.PartnerId, PartnerCode: resp.PartnerCode, Attributes: resp.Attributes, Error: resp.Error}, nil
}

// This helper function is required to translate Go error types to a string.
func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
