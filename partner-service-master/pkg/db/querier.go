package db

import (
	//"fmt"

	"fmt"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"jaxf-github.fanatics.corp/apparel/partner-service/pkg/db/queries"
)

type PartnerServiceQuerier interface {
	FindPartnerDataFromKeyValue(string, string) (int32, string, error) //KeyValue
	FindAllAttributesForPartner(int32) (map[string]string, error)      //Used by KeyValue and Id/code
	FindPartnerAttribute(int32, string) (map[string]string, error)     //Used by KeyValue and Id/code
	FindPartnerDataByID(int32, string) (int32, string, error)          //Id or code
	CheckPartnerIDEqualsPartnerCode(int32, string) (bool, error)       //check that the id and code correspond to same data
}

func NewPartnerServiceQuerier(c *pgx.Conn) PartnerServiceQuerier {
	return querier{
		conn: c,
	}
}

type querier struct {
	conn *pgx.Conn
}

//DB query for PartnerID from partners table
func (q querier) FindPartnerDataFromKeyValue(key, value string) (int32, string, error) { //DB query for PartnerID from partners table
	id, code, err := queries.GetPartnerDataFromKeyValue(key, value, q.conn)

	if err != nil || id == 0 {
		if err == nil {
			err = errors.Wrap(errors.New(""), fmt.Sprintf("error finding PartnerID from key: %s and value: %s in FindParnterID", key, value))
		} else {
			err = errors.Wrap(err, fmt.Sprintf("error finding PartnerID from key: %s and value: %s in FindParnterID", key, value))
		}
		return 0, "", err
	}
	return id, code, nil
}

func (q querier) FindAllAttributesForPartner(id int32) (map[string]string, error) { //DB query for PartnerAttribute from partner_mappings table
	attribute, err := queries.GetAllAttributesForPartner(id, q.conn)
	if err != nil {
		err = errors.Wrap(errors.New(""), fmt.Sprintf("error finding attributes in FindAllAttributes"))
		return make(map[string]string), err
	}
	return attribute, nil
}
func (q querier) FindPartnerAttribute(id int32, group string) (map[string]string, error) { //DB query for PartnerAttribute from partner_mappings table
	attribute, err := queries.GetGroupAttributesForPartner(id, group, q.conn)
	if err != nil {
		err = errors.Wrap(errors.New(""), fmt.Sprintf("error finding attributes in FindPartnerAttributes (by group)"))
		return make(map[string]string), err
	}
	return attribute, nil
}

func (q querier) FindPartnerDataByID(partnerId int32, code string) (int32, string, error) {
	id, code, err := queries.GetPartnerDataByIDOrCode(partnerId, code, q.conn)

	if err != nil {
		err = errors.Wrap(err, "error finding partnerData in FindPartnerIDbyID")
		return 0, "", err
	}
	return id, code, nil
}

func (q querier) CheckPartnerIDEqualsPartnerCode(partnerId int32, code string) (bool, error) {
	areEqual, err := queries.GetCheckPartnerIDEqualsPartnerCode(partnerId, code, q.conn)

	if err != nil {
		err = errors.Wrap(err, "partnerId and partnerCode correspond to different rows")
		return false, err
	}
	return areEqual, nil
}
