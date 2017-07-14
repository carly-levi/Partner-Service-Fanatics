package service

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

var ctx context.Context
var service PartnerService

type mockQuerier struct {
	mock.Mock
}

func (m *mockQuerier) FindPartnerDataByID(partnerId int32, partnerCode string) (int32, string, error) {
	args := m.Called(partnerId, partnerCode)
	typeInt32 := (args.Get(0).(int32))
	return typeInt32, args.String(1), args.Error(2)
}
func (m *mockQuerier) FindPartnerDataFromKeyValue(key, value string) (int32, string, error) {
	args := m.Called(key, value)
	typeInt32 := (args.Get(0).(int32))
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

// ServiceMethodsSuite allows us to attach setup and breakdown functions to multiple tests
type ServiceMethodsSuite struct {
	suite.Suite
}

func TestServiceMethods(t *testing.T) {
	suite.Run(t, new(ServiceMethodsSuite))
}

// SetupTest instantiates context and service fresh before every test
func (suite *ServiceMethodsSuite) SetupTest() {
	// gives us a blank context for each function
	ctx = context.Background()
	mq := new(mockQuerier)

	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	//when all goes well...
	mq.On("FindPartnerDataFromKeyValue", "Currency", "USD").Return(int32(1), "KOH", nil)
	mq.On("FindAllAttributesForPartner", int32(1)).Return(wantedMap, nil)
	mq.On("FindPartnerAttribute", int32(1), "Money").Return(wantedMap, nil)
	mq.On("FindPartnerDataByID", int32(1), "KOH").Return(int32(1), "KOH", nil)
	mq.On("CheckPartnerIDEqualsPartnerCode", int32(1), "KOH").Return(true, nil)
	//cases for when things are missing/bad inputs...
	mq.On("FindPartnerDataFromKeyValue", "", "USD").Return(int32(0), "", errors.New("error finding partner data from key value because empty key"))
	mq.On("FindPartnerDataFromKeyValue", "Currency", "").Return(int32(0), "", errors.New("error finding partner data from key value because empty value"))
	mq.On("FindPartnerDataFromKeyValue", "Currency", "asdfjkl").Return(int32(0), "", errors.New("error finding partner data from key value because bad value"))
	mq.On("FindPartnerDataFromKeyValue", "asdfjkl", "USD").Return(int32(0), "", errors.New("error finding partner data from key value because bad key"))
	mq.On("FindPartnerDataFromKeyValue", "", "").Return(int32(0), "", errors.New("error finding partner data from key value because both empty"))

	mq.On("FindPartnerDataByID", int32(1), "").Return(int32(1), "KOH", nil)
	mq.On("FindPartnerDataByID", int32(0), "KOH").Return(int32(1), "KOH", nil)
	mq.On("FindPartnerDataByID", int32(-1), "KOH").Return(int32(0), "", errors.New("error finding partner data from id or code because bad/negative id"))
	mq.On("FindPartnerDataByID", int32(1), "asdfjkl").Return(int32(0), "", errors.New("error finding partner data from id or code because bad code"))
	mq.On("FindPartnerDataByID", int32(0), "").Return(int32(0), "", errors.New("error finding partner data from id or code because both empty"))

	mq.On("FindAllAttributesForPartner", int32(0)).Return((make(map[string]string)), errors.New("error finding all attributes for Partner because nil id"))
	mq.On("FindAllAttributesForPartner", int32(-1)).Return((make(map[string]string)), errors.New("error finding all attributes for Partner because negative id"))

	mq.On("FindPartnerAttribute", int32(1), "").Return(wantedMap, nil)
	mq.On("FindPartnerAttribute", int32(-1), "Money").Return((make(map[string]string)), errors.New("error finding attributes for Partner & Group because bad/negative id"))
	mq.On("FindPartnerAttribute", int32(1), "asdfjkl").Return((make(map[string]string)), errors.New("error finding attributes for Partner & Group because bad group"))
	mq.On("FindPartnerAttribute", int32(0), "Money").Return((make(map[string]string)), errors.New("error finding attributes for Partner & Group because both empty"))

	mq.On("CheckPartnerIDEqualsPartnerCode", int32(1), "").Return(true, nil)
	mq.On("CheckPartnerIDEqualsPartnerCode", int32(0), "KOH").Return(true, nil)
	mq.On("CheckPartnerIDEqualsPartnerCode", int32(1), "asdfjkl").Return(false, errors.New("error checking if partnerId matches partnerCode because bad code"))
	mq.On("CheckPartnerIDEqualsPartnerCode", int32(-1), "KOH").Return(false, errors.New("error checking if partnerId matches partnerCode because bad/negative id"))
	mq.On("CheckPartnerIDEqualsPartnerCode", int32(0), "").Return(false, errors.New("error checking if partnerId matches partnerCode because both empty"))

	service = NewPartnerService(mq)
}

//test findPartnerDataFromKeyValue which only uses the GetPartnerDataByKeyValue
func (suite *ServiceMethodsSuite) TestFindPartnerDataFromKeyValueHappy() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetPartnerDataByKeyValue(ctx, "Currency", "USD", "Money")
	a.Nil(err)
	a.Equal(int32(1), partnerId)
	a.Equal("KOH", partnerCode)
	a.Equal(wantedMap, attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerDataFromKeyValueNilKey() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetPartnerDataByKeyValue(ctx, "", "USD", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerDataFromKeyValueNilValue() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetPartnerDataByKeyValue(ctx, "Currency", "", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerDataFromKeyValueBadKey() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetPartnerDataByKeyValue(ctx, "asdfjkl", "USD", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerDataFromKeyValueBadValue() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetPartnerDataByKeyValue(ctx, "Currency", "asdfjkl", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

//test FindAllAttributesForPartner when using GetPartnerDataByKeyValue
func (suite *ServiceMethodsSuite) TestFindAllAttributesForPartnerKVHappy() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetPartnerDataByKeyValue(ctx, "Currency", "USD", "Money")
	a.Nil(err)
	a.Equal(int32(1), partnerId)
	a.Equal("KOH", partnerCode)
	a.Equal(wantedMap, attributes)
}

func (suite *ServiceMethodsSuite) TestFindAllAttributesForPartnerNilKey() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetPartnerDataByKeyValue(ctx, "", "USD", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestFindAllAttributesForPartnerNilValue() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetPartnerDataByKeyValue(ctx, "Currency", "", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestFindAllAttributesForPartnerBadKey() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetPartnerDataByKeyValue(ctx, "asdfjkl", "USD", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestFindAllAttributesForPartnerBadValue() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetPartnerDataByKeyValue(ctx, "Currency", "asdfjkl", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

//test FindAllAttributesForPartner when using GetDataById

func (suite *ServiceMethodsSuite) TestFindAllAttributesForPartnerIDHappy() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(1), "KOH", "Money")
	a.Nil(err)
	a.Equal(int32(1), partnerId)
	a.Equal("KOH", partnerCode)
	a.Equal(wantedMap, attributes)
}

func (suite *ServiceMethodsSuite) TestFindAllAttributesForPartnerNilIdAndNilCode() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(0), "", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestFindAllAttributesForPartnerIDNilId() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(0), "KOH", "Money")
	a.Nil(err)
	a.Equal(int32(1), partnerId)
	a.Equal("KOH", partnerCode)
	a.Equal(wantedMap, attributes)
}

func (suite *ServiceMethodsSuite) TestFindAllAttributesForPartnerIDNilCode() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(1), "", "Money")
	a.Nil(err)
	a.Equal(int32(1), partnerId)
	a.Equal("KOH", partnerCode)
	a.Equal(wantedMap, attributes)
}

func (suite *ServiceMethodsSuite) TestFindAllAttributesForPartnerNegativeId() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(-1), "KOH", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestFindAllAttributesForPartnerBadId() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(-1), "KOH", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestFindAllAttributesForPartnerBadCode() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(1), "asdfjkl", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

//test FindPartnerAttribute when using GetPartnerDataByKeyValue
func (suite *ServiceMethodsSuite) TestFindPartnerAttributeHappyKV() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetPartnerDataByKeyValue(ctx, "Currency", "USD", "Money")
	a.Nil(err)
	a.Equal(int32(1), partnerId)
	a.Equal("KOH", partnerCode)
	a.Equal(wantedMap, attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerAttributeNilKey() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetPartnerDataByKeyValue(ctx, "", "USD", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerAttributeNilValue() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetPartnerDataByKeyValue(ctx, "Currency", "", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerAttributeNilGroup() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetPartnerDataByKeyValue(ctx, "Currency", "USD", "")
	a.Nil(err)
	a.Equal(int32(1), partnerId)
	a.Equal("KOH", partnerCode)
	a.Equal(wantedMap, attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerAttributeBadKey() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetPartnerDataByKeyValue(ctx, "asdfjkl", "USD", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerAttributeBadValue() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetPartnerDataByKeyValue(ctx, "Currency", "asdfjkl", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerAttributeBadGroup() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetPartnerDataByKeyValue(ctx, "Currency", "USD", "asdfjkl")
	a.NotNil(err)
	a.Equal(int32(1), partnerId)
	a.Equal("KOH", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

//test FindPartnerAttribute when using GetDataById
func (suite *ServiceMethodsSuite) TestFindPartnerAttributeIDHappy() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(1), "KOH", "Money")
	a.Nil(err)
	a.Equal(int32(1), partnerId)
	a.Equal("KOH", partnerCode)
	a.Equal(wantedMap, attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerAttributeNilIdAndNilCode() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(0), "", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerAttributeNilId() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(0), "KOH", "Money")
	a.Nil(err)
	a.Equal(int32(1), partnerId)
	a.Equal("KOH", partnerCode)
	a.Equal(wantedMap, attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerAttributeNilCode() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(1), "", "Money")
	a.Nil(err)
	a.Equal(int32(1), partnerId)
	a.Equal("KOH", partnerCode)
	a.Equal(wantedMap, attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerAttributeIDNilGroup() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(1), "KOH", "")
	a.Nil(err)
	a.Equal(int32(1), partnerId)
	a.Equal("KOH", partnerCode)
	a.Equal(wantedMap, attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerAttributeNegativeId() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(-1), "KOH", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerAttributeBadId() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(-1), "KOH", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerAttributeBadCode() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(1), "asdfjkl", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerAttributeIDBadGroup() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(1), "KOH", "asdfjkl")
	a.NotNil(err)
	a.Equal(int32(1), partnerId)
	a.Equal("KOH", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

//test FindPartnerDataByID which only uses the GetDataById
func (suite *ServiceMethodsSuite) TestFindPartnerDataByIDHappy() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(1), "KOH", "Money")
	a.Nil(err)
	a.Equal(int32(1), partnerId)
	a.Equal("KOH", partnerCode)
	a.Equal(wantedMap, attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerDataByIDNilId() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(0), "KOH", "Money")
	a.Nil(err)
	a.Equal(int32(1), partnerId)
	a.Equal("KOH", partnerCode)
	a.Equal(wantedMap, attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerDataByIDNilCode() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(1), "", "Money")
	a.Nil(err)
	a.Equal(int32(1), partnerId)
	a.Equal("KOH", partnerCode)
	a.Equal(wantedMap, attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerDataByNilIdAndNilCode() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(0), "", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerDataByNegativeId() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(-1), "KOH", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerDataByIDBadId() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(-1), "KOH", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestFindPartnerDataByIDBadCode() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(1), "asdfjkl", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

//test CheckPartnerIDEqualsPartnerCode which only uses the GetDataById
func (suite *ServiceMethodsSuite) TestCheckPartnerIDEqualsPartnerCodeHappy() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(1), "KOH", "Money")
	a.Nil(err)
	a.Equal(int32(1), partnerId)
	a.Equal("KOH", partnerCode)
	a.Equal(wantedMap, attributes)
}

func (suite *ServiceMethodsSuite) TestCheckPartnerIDEqualsPartnerCodeNilId() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(0), "KOH", "Money")
	a.Nil(err)
	a.Equal(int32(1), partnerId)
	a.Equal("KOH", partnerCode)
	a.Equal(wantedMap, attributes)
}

func (suite *ServiceMethodsSuite) TestCheckPartnerIDEqualsPartnerCodeNilCode() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(1), "", "Money")
	a.Nil(err)
	a.Equal(int32(1), partnerId)
	a.Equal("KOH", partnerCode)
	a.Equal(wantedMap, attributes)
}

func (suite *ServiceMethodsSuite) TestCheckPartnerIDEqualsPartnerCodeNilIdAndNilCode() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(0), "", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestCheckPartnerIDEqualsPartnerCodeNegativeId() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(-1), "KOH", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestCheckPartnerIDEqualsPartnerCodeBadId() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(-1), "KOH", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}

func (suite *ServiceMethodsSuite) TestCheckPartnerIDEqualsPartnerCodeBadCode() {
	a := assert.New(suite.T())
	partnerId, partnerCode, attributes, err := service.GetDataById(ctx, int32(1), "asdfjkl", "Money")
	a.NotNil(err)
	a.Equal(int32(0), partnerId)
	a.Equal("", partnerCode)
	a.Equal(make(map[string]string), attributes)
}
