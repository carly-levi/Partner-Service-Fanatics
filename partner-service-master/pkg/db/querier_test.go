package db

import (
	"testing"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"jaxf-github.fanatics.corp/apparel/partner-service/pkg/db/dbconfig"
)

var testQuerier PartnerServiceQuerier
var testConn *pgx.Conn
var err error

// ServiceMethodsSuite allows us to attach setup and breakdown functions to multiple tests
type QuerierMethodsSuite struct {
	suite.Suite
}

func TestServiceMethods(t *testing.T) {
	suite.Run(t, new(QuerierMethodsSuite))
}

// SetupTest instantiates context and service fresh before every test
func (suite *QuerierMethodsSuite) SetupTest() {
	// gives us a blank context for each function
	testConn, err = pgx.Connect(dbconfig.ExtractConfig())
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		suite.T().Error(err)
	}

	testConn.Exec("DROP TABLE keys cascade;")
	testConn.Exec("CREATE TABLE keys (id serial primary key,name varchar);")
	testConn.Exec("INSERT INTO keys (name) VALUES ('Currency');")
	testConn.Exec("INSERT INTO keys (name) VALUES ('Type of Payment');")

	testConn.Exec("DROP TABLE groups cascade;")
	testConn.Exec("CREATE TABLE groups (id serial primary key, name varchar);")
	testConn.Exec("INSERT INTO groups (name) VALUES ('EDI');")
	testConn.Exec("INSERT INTO groups (name) VALUES ('Style');")
	testConn.Exec("INSERT INTO groups (name) VALUES ('Money');")

	testConn.Exec("DROP TABLE partners cascade;")
	testConn.Exec("CREATE TABLE partners (id serial primary key, name varchar, code varchar);")
	testConn.Exec("INSERT INTO partners (name, code) VALUES ('Kohls', 'KOH');")

	testConn.Exec("DROP TABLE groups_to_keys cascade;")
	testConn.Exec("CREATE TABLE groups_to_keys (id serial primary key, group_id int, key_id int, FOREIGN KEY(group_id) REFERENCES groups(id), FOREIGN KEY(key_id) REFERENCES keys(id));")
	testConn.Exec("INSERT INTO groups_to_keys (group_id, key_id) VALUES (3, 1);")
	testConn.Exec("INSERT INTO groups_to_keys (group_id, key_id) VALUES (3, 2);")

	testConn.Exec("DROP TABLE partner_mappings cascade;")
	testConn.Exec("CREATE TABLE partner_mappings (id serial primary key, partner_id int, key_id int, FOREIGN KEY(partner_id) REFERENCES keys(id), FOREIGN KEY(key_id) REFERENCES keys(id), value varchar);")
	testConn.Exec("INSERT INTO partner_mappings (partner_id, key_id, value) VALUES (1, 1, 'USD');")
	testConn.Exec("INSERT INTO partner_mappings (partner_id, key_id, value) VALUES (1, 2, 'Credit');")
	testQuerier = NewPartnerServiceQuerier(testConn)
}

//tests for FindPartnerDataFromKeyValue
func (suite *QuerierMethodsSuite) TestFindPartnerDataFromKeyValueHappy() {
	a := assert.New(suite.T())

	id, code, err := testQuerier.FindPartnerDataFromKeyValue("Currency", "USD")

	a.Nil(err)
	a.Equal(int32(1), id)
	a.Equal("KOH", code)
}

func (suite *QuerierMethodsSuite) TestFindPartnerDataFromKeyValueBadKey() {
	a := assert.New(suite.T())

	id, code, err := testQuerier.FindPartnerDataFromKeyValue("lakjkhg", "USD")
	a.Equal(int32(0), id)
	a.Equal("", code)
	a.NotNil(err)
}

func (suite *QuerierMethodsSuite) TestFindPartnerDataFromKeyValueBadValue() {
	a := assert.New(suite.T())

	id, code, err := testQuerier.FindPartnerDataFromKeyValue("Currency", "lks;kdhf")
	a.Equal(int32(0), id)
	a.Equal("", code)
	a.NotNil(err)
}

func (suite *QuerierMethodsSuite) TestFindPartnerDataFromKeyValueNilKey() {
	a := assert.New(suite.T())

	id, code, err := testQuerier.FindPartnerDataFromKeyValue("", "USD")
	a.Equal(int32(0), id)
	a.Equal("", code)
	a.NotNil(err)
}

func (suite *QuerierMethodsSuite) TestFindPartnerDataFromKeyValueNilValue() {
	a := assert.New(suite.T())

	id, code, err := testQuerier.FindPartnerDataFromKeyValue("Currency", "")
	a.Equal(int32(0), id)
	a.Equal("", code)
	a.NotNil(err)
}

//tests for FindPartnerDataByID
func (suite *QuerierMethodsSuite) TestFindPartnerDataByIDHappy() {
	a := assert.New(suite.T())

	id, code, err := testQuerier.FindPartnerDataByID(int32(1), "KOH")

	a.Nil(err)
	a.Equal(int32(1), id)
	a.Equal("KOH", code)
}

func (suite *QuerierMethodsSuite) TestFindPartnerDataByIDNilId() {
	a := assert.New(suite.T())

	id, code, err := testQuerier.FindPartnerDataByID(int32(0), "KOH")
	a.Nil(err)
	a.Equal(int32(1), id)
	a.Equal("KOH", code)
}

func (suite *QuerierMethodsSuite) TestFindPartnerDataByIDNilCode() {
	a := assert.New(suite.T())

	id, code, err := testQuerier.FindPartnerDataByID(int32(1), "")
	a.Nil(err)
	a.Equal(int32(1), id)
	a.Equal("KOH", code)
}

func (suite *QuerierMethodsSuite) TestFindPartnerDataByIDNilIdAndCode() {
	a := assert.New(suite.T())

	id, code, err := testQuerier.FindPartnerDataByID(int32(0), "")
	a.Equal(int32(0), id)
	a.Equal("", code)
	a.NotNil(err)
}

//tests for CheckPartnerIDEqualsPartnerCode
func (suite *QuerierMethodsSuite) TestCheckPartnerIDEqualsPartnerCodeHappy() {
	a := assert.New(suite.T())

	areEqual, err := testQuerier.CheckPartnerIDEqualsPartnerCode(int32(1), "KOH")
	a.Nil(err)
	a.Equal(true, areEqual)
}

func (suite *QuerierMethodsSuite) TestCheckPartnerIDEqualsPartnerCodeBadId() {
	a := assert.New(suite.T())

	areEqual, err := testQuerier.CheckPartnerIDEqualsPartnerCode(-8586, "KOH")
	a.Equal(false, areEqual)
	a.NotNil(err)
}

func (suite *QuerierMethodsSuite) TestCheckPartnerIDEqualsPartnerCodeBadCode() {
	a := assert.New(suite.T())

	areEqual, err := testQuerier.CheckPartnerIDEqualsPartnerCode(int32(1), "lahgk")
	a.Equal(false, areEqual)
	a.NotNil(err)
}

//tests for FindAllAttributesForPartner
func (suite *QuerierMethodsSuite) TestFindAllAttributesForPartnerHappy() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	attributes, err := testQuerier.FindAllAttributesForPartner(int32(1))
	a.Nil(err)
	a.Equal(wantedMap, attributes)
}

func (suite *QuerierMethodsSuite) TestFindAllAttributesForPartnerBadId() {
	a := assert.New(suite.T())

	attributes, err := testQuerier.FindAllAttributesForPartner(int32(-1))
	a.Equal(make(map[string]string), attributes)
	a.NotNil(err)
}

func (suite *QuerierMethodsSuite) TestFindAllAttributesForPartnerNilId() {
	a := assert.New(suite.T())

	attributes, err := testQuerier.FindAllAttributesForPartner(int32(0))
	a.Equal(make(map[string]string), attributes)
	a.NotNil(err)
}

func (suite *QuerierMethodsSuite) TestFindAllAttributesForPartnerNegativeId() {
	a := assert.New(suite.T())

	attributes, err := testQuerier.FindAllAttributesForPartner(int32(-1))
	a.Equal(make(map[string]string), attributes)
	a.NotNil(err)
}

//tests for FindPartnerAttribute
func (suite *QuerierMethodsSuite) TestFindPartnerAttributeHappy() {
	a := assert.New(suite.T())
	wantedMap := make(map[string]string)
	wantedMap["Currency"] = "USD"
	wantedMap["Type of Payment"] = "Credit"
	attributes, err := testQuerier.FindPartnerAttribute(int32(1), "Money")
	a.Nil(err)
	a.Equal(wantedMap, attributes)
}

func (suite *QuerierMethodsSuite) TestFindPartnerAttributeBadId() {
	a := assert.New(suite.T())

	attributes, err := testQuerier.FindPartnerAttribute(int32(-1), "Money")
	a.Equal(make(map[string]string), attributes)
	a.NotNil(err)
}

func (suite *QuerierMethodsSuite) TestFindPartnerAttributeNilId() {
	a := assert.New(suite.T())

	attributes, err := testQuerier.FindPartnerAttribute(int32(0), "Money")
	a.Equal(make(map[string]string), attributes)
	a.NotNil(err)
}

func (suite *QuerierMethodsSuite) TestFindPartnerAttributeNegativeId() {
	a := assert.New(suite.T())

	attributes, err := testQuerier.FindPartnerAttribute(int32(-1), "Money")
	a.Equal(make(map[string]string), attributes)
	a.NotNil(err)
}

func (suite *QuerierMethodsSuite) TestFindPartnerAttributeBadGroup() {
	a := assert.New(suite.T())

	attributes, err := testQuerier.FindPartnerAttribute(int32(1), "lshg")
	a.Equal(make(map[string]string), attributes)
	a.NotNil(err)
}

func (suite *QuerierMethodsSuite) TestFindPartnerAttributeNilGroup() {
	a := assert.New(suite.T())

	attributes, err := testQuerier.FindPartnerAttribute(int32(1), "")
	a.Equal(make(map[string]string), attributes)
	a.NotNil(err)
}
