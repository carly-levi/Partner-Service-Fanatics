package endpoints

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"jaxf-github.fanatics.corp/apparel/partner-service/pkg/service"
)

type mockQuerier struct {
	mock.Mock
}

func (m *mockQuerier) FindPartnerDataByID(partnerId int32, partnerCode string) (int32, string, error) {
	args := m.Called(partnerId, partnerCode)
	//have to cast to int32 from regular in because args.Int32 is not acceptable
	typeInt32 := args.Get(0).(int32)
	return typeInt32, args.String(1), args.Error(2)
}
func (m *mockQuerier) FindPartnerDataFromKeyValue(key, value string) (int32, string, error) {
	args := m.Called(key, value)
	typeInt32 := args.Get(0).(int32)
	return typeInt32, args.String(1), args.Error(2)
}
func (m *mockQuerier) FindAllAttributesForPartner(partnerId int32) (map[string]string, error) {
	args := m.Called(partnerId)
	typeMapStringString := args.Get(0).(map[string]string)
	return typeMapStringString, args.Error(1)
}
func (m *mockQuerier) FindPartnerAttribute(partnerId int32, group string) (map[string]string, error) {
	args := m.Called(partnerId, group)
	typeMapStringString := args.Get(0).(map[string]string)
	return typeMapStringString, args.Error(1)
}

func (m *mockQuerier) CheckPartnerIDEqualsPartnerCode(partnerId int32, partnerCode string) (bool, error) {
	args := m.Called(partnerId, partnerCode)
	return args.Bool(0), args.Error(1)
}

func TestMakeKeyValueEndpointHappy(t *testing.T) {
	a := assert.New(t)
	mq := new(mockQuerier)
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	mq.On("FindPartnerDataFromKeyValue", "Currency", "USD").Return(int32(1), "KOH", nil)
	mq.On("FindAllAttributesForPartner", int32(1)).Return(wantedMap, nil)
	mq.On("FindPartnerAttribute", int32(1), "Money").Return(wantedMap, nil)

	s := service.NewPartnerService(mq)

	req := &KeyValueRequest{
		Key:   "Currency",
		Value: "USD",
		Group: "Money",
	}

	ctx := context.Background()

	res, err := MakeKeyValueEndpoint(s)(ctx, *req)

	a.Equal(int32(1), int32(res.(PartnerDataReply).PartnerId))
	a.Equal("KOH", res.(PartnerDataReply).PartnerCode)
	a.Equal(wantedMap, res.(PartnerDataReply).Attributes)
	a.Nil(err)
}

func TestMakeKeyValueEndpointBadKey(t *testing.T) {
	a := assert.New(t)
	mq := new(mockQuerier)
	mq.On("FindPartnerDataFromKeyValue", "lkshdglk", "USD").Return(int32(0), "", errors.New("error finding partner data from key value"))
	mq.On("FindAllAttributesForPartner", int32(0)).Return((make(map[string]string)), errors.New("error finding all attributes for Partner"))
	mq.On("FindPartnerAttribute", int32(0), "Money").Return((make(map[string]string)), errors.New("error finding attributes for Partner & Group"))

	s := service.NewPartnerService(mq)

	req := &KeyValueRequest{
		Key:   "lkshdglk", //bad key
		Value: "USD",
		Group: "Money",
	}

	ctx := context.Background()

	res, _ := MakeKeyValueEndpoint(s)(ctx, *req)

	a.Equal(int32(0), res.(PartnerDataReply).PartnerId)
	a.Equal("", res.(PartnerDataReply).PartnerCode)
	a.Equal((make(map[string]string)), res.(PartnerDataReply).Attributes)
}
func TestMakeKeyValueEndpointBadValue(t *testing.T) {
	a := assert.New(t)
	mq := new(mockQuerier)
	mq.On("FindPartnerDataFromKeyValue", "Currency", "sgdsd").Return(int32(0), "", errors.New("error finding partner data from key value"))
	mq.On("FindAllAttributesForPartner", int32(0)).Return((make(map[string]string)), errors.New("error finding all attributes for Partner"))
	mq.On("FindPartnerAttribute", int32(0), "Money").Return((make(map[string]string)), errors.New("error finding attributes for Partner & Group"))

	s := service.NewPartnerService(mq)

	req := &KeyValueRequest{
		Key:   "Currency",
		Value: "sgdsd", //bad value
		Group: "Money",
	}

	ctx := context.Background()

	res, _ := MakeKeyValueEndpoint(s)(ctx, *req)

	a.Equal(int32(0), res.(PartnerDataReply).PartnerId)
	a.Equal("", res.(PartnerDataReply).PartnerCode)
	a.Equal((make(map[string]string)), res.(PartnerDataReply).Attributes)
}

func TestMakeKeyValueEndpointBadGroup(t *testing.T) {
	a := assert.New(t)
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	mq := new(mockQuerier)
	mq.On("FindPartnerDataFromKeyValue", "Currency", "USD").Return(int32(1), "KOH", nil)
	mq.On("FindAllAttributesForPartner", int32(1)).Return(wantedMap, nil)
	mq.On("FindPartnerAttribute", int32(1), "lksdhf").Return((make(map[string]string)), errors.New("error finding attributes for Partner & Group because bad group"))

	s := service.NewPartnerService(mq)

	req := &KeyValueRequest{
		Key:   "Currency",
		Value: "USD",
		Group: "lksdhf", //bad value
	}

	ctx := context.Background()

	res, _ := MakeKeyValueEndpoint(s)(ctx, *req)
	a.Equal(int32(1), res.(PartnerDataReply).PartnerId)
	a.Equal("KOH", res.(PartnerDataReply).PartnerCode)
	a.Equal(make(map[string]string), res.(PartnerDataReply).Attributes)
}

func TestMakeKeyValueEndpointNilKey(t *testing.T) {
	a := assert.New(t)
	mq := new(mockQuerier)
	mq.On("FindPartnerDataFromKeyValue", "", "USD").Return(int32(0), "", errors.New("error finding partner data from nil key value"))
	mq.On("FindAllAttributesForPartner", int32(0)).Return((make(map[string]string)), errors.New("error finding all attributes for Partner because nil key"))
	mq.On("FindPartnerAttribute", int32(0), "Money").Return((make(map[string]string)), errors.New("error finding attributes for Partner & Group because nil key"))

	s := service.NewPartnerService(mq)

	req := &KeyValueRequest{
		Key:   "", //nil key
		Value: "USD",
		Group: "Money",
	}

	ctx := context.Background()

	res, _ := MakeKeyValueEndpoint(s)(ctx, *req)

	a.Equal(int32(0), res.(PartnerDataReply).PartnerId)
	a.Equal("", res.(PartnerDataReply).PartnerCode)
	a.Equal((make(map[string]string)), res.(PartnerDataReply).Attributes)
}

func TestMakeKeyValueEndpointNilValue(t *testing.T) {
	a := assert.New(t)
	mq := new(mockQuerier)
	mq.On("FindPartnerDataFromKeyValue", "Currency", "").Return(int32(0), "", errors.New("error finding partner data from key value because nil value"))
	mq.On("FindAllAttributesForPartner", int32(0)).Return((make(map[string]string)), errors.New("error finding all attributes for Partner because nil value"))
	mq.On("FindPartnerAttribute", int32(0), "Money").Return((make(map[string]string)), errors.New("error finding attributes for Partner & Group because nil value"))

	s := service.NewPartnerService(mq)

	req := &KeyValueRequest{
		Key:   "Currency",
		Value: "", //nil value
		Group: "Money",
	}

	ctx := context.Background()

	res, _ := MakeKeyValueEndpoint(s)(ctx, *req)

	a.Equal(int32(0), res.(PartnerDataReply).PartnerId)
	a.Equal("", res.(PartnerDataReply).PartnerCode)
	a.Equal((make(map[string]string)), res.(PartnerDataReply).Attributes)
}

func TestMakeKeyValueEndpointNilGroup(t *testing.T) {
	a := assert.New(t)
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	mq := new(mockQuerier)
	mq.On("FindPartnerDataFromKeyValue", "Currency", "USD").Return(int32(1), "KOH", nil)
	mq.On("FindAllAttributesForPartner", int32(1)).Return(wantedMap, nil)
	mq.On("FindPartnerAttribute", int32(1), "").Return(wantedMap, nil)

	s := service.NewPartnerService(mq)

	req := &KeyValueRequest{
		Key:   "Currency",
		Value: "USD",
		Group: "", //nil value
	}

	ctx := context.Background()

	res, _ := MakeKeyValueEndpoint(s)(ctx, *req)
	a.Equal(int32(1), res.(PartnerDataReply).PartnerId)
	a.Equal("KOH", res.(PartnerDataReply).PartnerCode)
	a.Equal(wantedMap, res.(PartnerDataReply).Attributes)
}

//tests for the second endpoint MakeDataByIdEndpoint
func TestMakeDataByIdEndpointHappy(t *testing.T) {
	a := assert.New(t)
	mq := new(mockQuerier)
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	mq.On("FindPartnerDataByID", int32(1), "KOH").Return(int32(1), "KOH", nil)
	mq.On("FindAllAttributesForPartner", int32(1)).Return(wantedMap, nil)
	mq.On("FindPartnerAttribute", int32(1), "Money").Return(wantedMap, nil)
	mq.On("CheckPartnerIDEqualsPartnerCode", int32(1), "KOH").Return(true, nil)

	s := service.NewPartnerService(mq)

	req := &IdRequest{
		PartnerId:   int32(1),
		PartnerCode: "KOH",
		Group:       "Money",
	}

	ctx := context.Background()

	res, err := MakeGetDataByIdEndpoint(s)(ctx, *req)

	a.Equal(int32(1), res.(PartnerDataReply).PartnerId)
	a.Equal("KOH", res.(PartnerDataReply).PartnerCode)
	a.Equal(wantedMap, res.(PartnerDataReply).Attributes)
	a.Nil(err)
}

func TestMakeGetDataByIdEndpointBadId(t *testing.T) {
	a := assert.New(t)
	mq := new(mockQuerier)
	mq.On("FindPartnerDataByID", int32(-1), "KOH").Return(int32(0), "", errors.New("error finding partner data from id or Code because bad id"))
	mq.On("FindAllAttributesForPartner", int32(0)).Return((make(map[string]string)), errors.New("error finding all attributes for Partner because bad id"))
	mq.On("FindPartnerAttribute", int32(0), "Money").Return((make(map[string]string)), errors.New("error finding attributes for Partner & Group because bad id"))
	mq.On("CheckPartnerIDEqualsPartnerCode", int32(-1), "KOH").Return(false, errors.New("error checking if partnerId matches partnerCode because bad id"))

	s := service.NewPartnerService(mq)

	req := &IdRequest{
		PartnerId:   int32(-1),
		PartnerCode: "KOH",
		Group:       "Money",
	}

	ctx := context.Background()

	res, _ := MakeGetDataByIdEndpoint(s)(ctx, *req)

	a.Equal(int32(0), res.(PartnerDataReply).PartnerId)
	a.Equal("", res.(PartnerDataReply).PartnerCode)
	a.Equal((make(map[string]string)), res.(PartnerDataReply).Attributes)
}

func TestMakeGetDataByIdEndpointBadCode(t *testing.T) {
	a := assert.New(t)
	mq := new(mockQuerier)
	mq.On("FindPartnerDataByID", int32(1), "lhdfhg").Return(int32(0), "", errors.New("error finding partner data from id or Code because bad code"))
	mq.On("FindAllAttributesForPartner", int32(0)).Return((make(map[string]string)), errors.New("error finding all attributes for Partner because bad code"))
	mq.On("FindPartnerAttribute", int32(0), "Money").Return((make(map[string]string)), errors.New("error finding attributes for Partner & Group because bad code"))
	mq.On("CheckPartnerIDEqualsPartnerCode", int32(1), "lhdfhg").Return(false, errors.New("error checking if partnerId matches partnerCode because bad code"))

	s := service.NewPartnerService(mq)

	req := &IdRequest{
		PartnerId:   int32(1),
		PartnerCode: "lhdfhg",
		Group:       "Money",
	}

	ctx := context.Background()

	res, _ := MakeGetDataByIdEndpoint(s)(ctx, *req)

	a.Equal(int32(0), res.(PartnerDataReply).PartnerId)
	a.Equal("", res.(PartnerDataReply).PartnerCode)
	a.Equal((make(map[string]string)), res.(PartnerDataReply).Attributes)
}

func TestMakeGetDataByIdEndpointBadGroup(t *testing.T) {
	a := assert.New(t)
	mq := new(mockQuerier)
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	mq.On("FindPartnerDataByID", int32(1), "KOH").Return(int32(1), "KOH", nil)
	mq.On("FindAllAttributesForPartner", int32(1)).Return(wantedMap, nil)
	mq.On("FindPartnerAttribute", int32(1), "sldkgh").Return((make(map[string]string)), errors.New("error finding attributes for Partner & Group because bad group"))
	mq.On("CheckPartnerIDEqualsPartnerCode", int32(1), "KOH").Return(true, nil)

	s := service.NewPartnerService(mq)

	req := &IdRequest{
		PartnerId:   int32(1),
		PartnerCode: "KOH",
		Group:       "sldkgh",
	}

	ctx := context.Background()

	res, _ := MakeGetDataByIdEndpoint(s)(ctx, *req)

	a.Equal(int32(1), res.(PartnerDataReply).PartnerId)
	a.Equal("KOH", res.(PartnerDataReply).PartnerCode)
	a.Equal((make(map[string]string)), res.(PartnerDataReply).Attributes)
}

func TestMakeGetDataByIdEndpointNilId(t *testing.T) {
	a := assert.New(t)
	mq := new(mockQuerier)
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	mq.On("FindPartnerDataByID", int32(0), "KOH").Return(int32(1), "KOH", nil)
	mq.On("FindAllAttributesForPartner", int32(1)).Return(wantedMap, nil)
	mq.On("FindPartnerAttribute", int32(1), "Money").Return(wantedMap, nil)
	mq.On("CheckPartnerIDEqualsPartnerCode", int32(0), "KOH").Return(true, nil)
	s := service.NewPartnerService(mq)

	req := &IdRequest{
		PartnerId:   int32(0),
		PartnerCode: "KOH",
		Group:       "Money",
	}

	ctx := context.Background()

	res, _ := MakeGetDataByIdEndpoint(s)(ctx, *req)

	a.Equal(int32(1), res.(PartnerDataReply).PartnerId)
	a.Equal("KOH", res.(PartnerDataReply).PartnerCode)
	a.Equal(wantedMap, res.(PartnerDataReply).Attributes)
}

func TestMakeGetDataByIdEndpointNilCode(t *testing.T) {
	a := assert.New(t)
	mq := new(mockQuerier)

	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	mq.On("FindPartnerDataByID", int32(1), "").Return(int32(1), "KOH", nil)
	mq.On("FindAllAttributesForPartner", int32(1)).Return(wantedMap, nil)
	mq.On("FindPartnerAttribute", int32(1), "Money").Return(wantedMap, nil)
	mq.On("CheckPartnerIDEqualsPartnerCode", int32(1), "").Return(true, nil)

	s := service.NewPartnerService(mq)

	req := &IdRequest{
		PartnerId:   int32(1),
		PartnerCode: "",
		Group:       "Money",
	}

	ctx := context.Background()

	res, _ := MakeGetDataByIdEndpoint(s)(ctx, *req)

	a.Equal(int32(1), res.(PartnerDataReply).PartnerId)
	a.Equal("KOH", res.(PartnerDataReply).PartnerCode)
	a.Equal(wantedMap, res.(PartnerDataReply).Attributes)
}

func TestMakeGetDataByIdEndpointNilGroup(t *testing.T) {
	a := assert.New(t)
	mq := new(mockQuerier)
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	mq.On("FindPartnerDataByID", int32(1), "KOH").Return(int32(1), "KOH", nil)
	mq.On("FindAllAttributesForPartner", int32(1)).Return(wantedMap, nil)
	mq.On("FindPartnerAttribute", int32(1), "").Return(wantedMap, nil)
	mq.On("CheckPartnerIDEqualsPartnerCode", int32(1), "KOH").Return(true, nil)

	s := service.NewPartnerService(mq)

	req := &IdRequest{
		PartnerId:   int32(1),
		PartnerCode: "KOH",
		Group:       "",
	}

	ctx := context.Background()

	res, _ := MakeGetDataByIdEndpoint(s)(ctx, *req)

	a.Equal(int32(1), res.(PartnerDataReply).PartnerId)
	a.Equal("KOH", res.(PartnerDataReply).PartnerCode)
	a.Equal(wantedMap, res.(PartnerDataReply).Attributes)
}

func TestMakeGetDataByIdEndpointNegativeId(t *testing.T) {
	a := assert.New(t)
	mq := new(mockQuerier)
	mq.On("FindPartnerDataByID", int32(-1), "KOH").Return(int32(0), "", errors.New("error finding partner data from id or Code because negative id"))
	mq.On("FindAllAttributesForPartner", int32(-1)).Return((make(map[string]string)), errors.New("error finding all attributes for Partner because negative id"))
	mq.On("FindPartnerAttribute", int32(-1), "Money").Return((make(map[string]string)), errors.New("error finding attributes for Partner & Group because negative id"))
	mq.On("CheckPartnerIDEqualsPartnerCode", int32(-1), "KOH").Return(false, errors.New("error checking if partnerId matches partnerCode because negative id"))

	s := service.NewPartnerService(mq)

	req := &IdRequest{
		PartnerId:   int32(-1),
		PartnerCode: "KOH",
		Group:       "Money",
	}

	ctx := context.Background()

	res, _ := MakeGetDataByIdEndpoint(s)(ctx, *req)

	a.Equal(int32(0), res.(PartnerDataReply).PartnerId)
	a.Equal("", res.(PartnerDataReply).PartnerCode)
	a.Equal((make(map[string]string)), res.(PartnerDataReply).Attributes)
}
