package models

import (
	"testing"

	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
)

func TestPartnerWithAllValues(t *testing.T) {
	ns := pgx.NullString{String: "Dicks", Valid: true}
  cs := pgx.NullString{String: "DIC", Valid: true}
  is := pgx.NullInt32{Int32: 5, Valid: true}
  m := map[string]string{
    "Currency": "USD",
    "Type of Payment": "Debit",
    "Color": "Black",
    "Gender": "Women",
    "Sleeves": "Short"}
	partnerModel := &Partner{
		Name: ns,
    Code: cs,
    Id: is,
    Attributes: m,
	}

	partner := partnerModel.Gen(m)
	assert.Equal(t, "Dicks", partner.Name)
  assert.Equal(t, "DIC", partner.Code)
  assert.Equal(t, int32(5), partner.Id)
  assert.Equal(t, m, partner.Attributes)
}

func TestPartnerValueWithNameNil(t *testing.T) {
  ns := pgx.NullString{String: "", Valid: false}
  cs := pgx.NullString{String: "DIC", Valid: true}
  is := pgx.NullInt32{Int32: 5, Valid: true}
  m := map[string]string{
    "Currency": "USD",
    "Type of Payment": "Debit",
    "Color": "Black",
    "Gender": "Women",
    "Sleeves": "Short"}
	partnerModel := &Partner{
    Name: ns,
    Code: cs,
    Id: is,
    Attributes: m,
	}

	partner := partnerModel.Gen(m)
	assert.Equal(t, "", partner.Name)
  assert.Equal(t, "DIC", partner.Code)
  assert.Equal(t, int32(5), partner.Id)
  assert.Equal(t, m, partner.Attributes)
}

func TestPartnerValueWithCodeNil(t *testing.T) {
  ns := pgx.NullString{String: "Dicks", Valid: false}
  cs := pgx.NullString{String: "", Valid: true}
  is := pgx.NullInt32{Int32: 5, Valid: true}
  m := map[string]string{
    "Currency": "USD",
    "Type of Payment": "Debit",
    "Color": "Black",
    "Gender": "Women",
    "Sleeves": "Short"}
	partnerModel := &Partner{
    Name: ns,
    Code: cs,
    Id: is,
    Attributes: m,
	}

	partner := partnerModel.Gen(m)
	assert.Equal(t, "Dicks", partner.Name)
  assert.Equal(t, "", partner.Code)
  assert.Equal(t, int32(5), partner.Id)
  assert.Equal(t, m, partner.Attributes)
}

func TestPartnerValueWithIDNil(t *testing.T) {
  ns := pgx.NullString{String: "Dicks", Valid: true}
  cs := pgx.NullString{String: "DIC", Valid: true}
  is := pgx.NullInt32{Int32: -1, Valid: false}
  m := map[string]string{
    "Currency": "USD",
    "Type of Payment": "Debit",
    "Color": "Black",
    "Gender": "Women",
    "Sleeves": "Short"}
	partnerModel := &Partner{
    Name: ns,
    Code: cs,
    Id: is,
    Attributes: m,
	}

	partner := partnerModel.Gen(m)
	assert.Equal(t, "Dicks", partner.Name)
  assert.Equal(t, "DIC", partner.Code)
  assert.Equal(t, int32(-1), partner.Id)
  assert.Equal(t, m, partner.Attributes)
}

func TestPartnerValueWithAttributeNil(t *testing.T) {
  ns := pgx.NullString{String: "Dicks", Valid: true}
  cs := pgx.NullString{String: "DIC", Valid: true}
  is := pgx.NullInt32{Int32: 5, Valid: true}
  m := map[string]string{}
	partnerModel := &Partner{
    Name: ns,
    Code: cs,
    Id: is,
    Attributes: m,
	}

	partner := partnerModel.Gen(m)
	assert.Equal(t, "Dicks", partner.Name)
  assert.Equal(t, "DIC", partner.Code)
  assert.Equal(t, int32(5), partner.Id)
  assert.Equal(t, m, partner.Attributes)
}

func TestPartnerValueWithAllNil(t *testing.T) {
  ns := pgx.NullString{String: "", Valid: false}
  cs := pgx.NullString{String: "", Valid: false}
  is := pgx.NullInt32{Int32: -1, Valid: false}
  m := map[string]string{}
	partnerModel := &Partner{
    Name: ns,
    Code: cs,
    Id: is,
    Attributes: m,
	}

	partner := partnerModel.Gen(m)
	assert.Equal(t, "", partner.Name)
  assert.Equal(t, "", partner.Code)
  assert.Equal(t, int32(-1), partner.Id)
  assert.Equal(t, m, partner.Attributes)
}
