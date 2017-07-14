package transport_grpc

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"jaxf-github.fanatics.corp/apparel/partner-service/pkg/endpoints"
	"jaxf-github.fanatics.corp/apparel/partner-service/pkg/pb"
)

// Test err2str
func TestErr2StrValue(t *testing.T) {
	err := errors.New("an error")

	val := err2str(err)

	assert.Equal(t, "an error", val)
}

func TestErr2StrNil(t *testing.T) {
	var err error

	val := err2str(err)

	assert.Equal(t, "", val)
}

// Test KeyValue decode function
func TestDecodeGRPCKeyValueRequest(t *testing.T) { //Do we need separate tests for group and no group?
	ctx := context.Background()
	hr := &pb.KeyValueRequest{
		Key:   "Currency",
		Value: "USD",
		Group: "Money",
	}

	decReq, err := DecodeGRPCKeyValueRequest(ctx, hr)

	assert.Equal(t, "Currency", decReq.(endpoints.KeyValueRequest).Key)
	assert.Equal(t, "USD", decReq.(endpoints.KeyValueRequest).Value)
	assert.Equal(t, "Money", decReq.(endpoints.KeyValueRequest).Group)

	assert.Nil(t, err)
}

func TestDecodeGRPCKeyValueRequestEmpty(t *testing.T) {
	ctx := context.Background()
	hr := &pb.KeyValueRequest{
		Key:   "",
		Value: "",
		Group: "",
	}

	decReq, err := DecodeGRPCKeyValueRequest(ctx, hr)

	assert.Equal(t, "", decReq.(endpoints.KeyValueRequest).Key)
	assert.Equal(t, "", decReq.(endpoints.KeyValueRequest).Value)
	assert.Equal(t, "", decReq.(endpoints.KeyValueRequest).Group)
	assert.Nil(t, err)
}

// Test ID or Code decode function
func TestDecodeGRPCDataByIdRequestID(t *testing.T) {
	ctx := context.Background()
	hr := &pb.IdRequest{
		PartnerId:   1,
		PartnerCode: "KOH",
		Group:       "Money",
	}

	decReq, err := DecodeGRPCDataByIdRequest(ctx, hr)

	assert.Equal(t, int32(1), decReq.(endpoints.IdRequest).PartnerId)
	assert.Equal(t, "KOH", decReq.(endpoints.IdRequest).PartnerCode)
	assert.Equal(t, "Money", decReq.(endpoints.IdRequest).Group)

	assert.Nil(t, err)
}

func TestDecodeGRPCDataByIdRequestEmpty(t *testing.T) {
	ctx := context.Background()
	hr := &pb.IdRequest{
		PartnerId:   0,
		PartnerCode: "",
		Group:       "",
	}

	decReq, err := DecodeGRPCDataByIdRequest(ctx, hr)

	assert.Equal(t, int32(0), decReq.(endpoints.IdRequest).PartnerId)
	assert.Equal(t, "", decReq.(endpoints.IdRequest).PartnerCode)
	assert.Equal(t, "", decReq.(endpoints.IdRequest).Group)

	assert.Nil(t, err)
}

// Test encode function
func TestEncodeGRPCResponses(t *testing.T) {
	ctx := context.Background()
	hr := &endpoints.PartnerDataReply{
		PartnerId:   1,
		PartnerCode: "KOH",
		Attributes:  map[string]string{"Currency": "USD", "Type of Payment": "Credit"},
		Error:       "",
	}

	encRep, err := EncodeGRPCResponse(ctx, *hr)

	assert.Equal(t, int32(1), encRep.(*pb.PartnerDataReply).PartnerId)
	assert.Equal(t, "KOH", encRep.(*pb.PartnerDataReply).PartnerCode)
	assert.Equal(t, map[string]string{"Currency": "USD", "Type of Payment": "Credit"}, encRep.(*pb.PartnerDataReply).Attributes)
	assert.Equal(t, "", encRep.(*pb.PartnerDataReply).Error)
	assert.Nil(t, err)
}

func TestEncodeGRPCResponseErr(t *testing.T) {
	ctx := context.Background()
	err := errors.New("test error")
	hr := &endpoints.PartnerDataReply{
		PartnerId:   0,
		PartnerCode: "",
		Attributes:  make(map[string]string),
		Error:       err2str(err),
	}

	encRep, err := EncodeGRPCResponse(ctx, *hr)

	assert.Equal(t, int32(0), encRep.(*pb.PartnerDataReply).PartnerId)
	assert.Equal(t, "", encRep.(*pb.PartnerDataReply).PartnerCode)
	assert.Equal(t, make(map[string]string), encRep.(*pb.PartnerDataReply).Attributes)
	assert.Equal(t, "test error", encRep.(*pb.PartnerDataReply).Error)
	assert.Nil(t, err)
}
