package queries

import (
	"testing"
	"jaxf-github.fanatics.corp/apparel/partner-service/pkg/db/dbconfig"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func CreateDB() {
  conn, err := pgx.Connect(dbconfig.ExtractConfig())
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		//t.Error(err)
	}
	defer conn.Close()

  conn.Exec("DROP TABLE partner_mappings;")
  conn.Exec("DROP TABLE groups_to_keys;")
  conn.Exec("DROP TABLE partners;")
  conn.Exec("DROP TABLE keys;")
  conn.Exec("DROP TABLE groups;")

  conn.Exec("CREATE TABLE keys (id serial primary key, name varchar(255) not null);")
  conn.Exec("CREATE TABLE groups (id serial primary key, name varchar(255) not null);")
  conn.Exec("CREATE TABLE partners (id serial primary key, name varchar(255) not null, code varchar(255) not null);")
  conn.Exec("CREATE TABLE partner_mappings (id serial primary key, partner_id int not null, key_id int not null, FOREIGN KEY(partner_id) REFERENCES keys(id), FOREIGN KEY(key_id) REFERENCES keys(id), value varchar(255) not null);")
  conn.Exec("CREATE TABLE groups_to_keys (id serial primary key, group_id int not null, key_id int not null, FOREIGN KEY(group_id) REFERENCES groups(id), FOREIGN KEY(key_id) REFERENCES keys(id));")

  conn.Exec("INSERT INTO keys (name) VALUES ('Currency');")
  conn.Exec("INSERT INTO keys (name) VALUES ('Type of Payment');")
  conn.Exec("INSERT INTO keys (name) VALUES ('Color');")

  conn.Exec("INSERT INTO partners (name, code) VALUES ('Kohls', 'KOH');")

  conn.Exec("INSERT INTO partner_mappings (partner_id, key_id, value) VALUES (1, 1, 'USD');")
  conn.Exec("INSERT INTO partner_mappings (partner_id, key_id, value) VALUES (1, 2, 'Credit');")
  conn.Exec("INSERT INTO partner_mappings (partner_id, key_id, value) VALUES (1, 3, 'Blue');")

  conn.Exec("INSERT INTO groups (name) VALUES ('Style');")
  conn.Exec("INSERT INTO groups (name) VALUES ('Money');")

  conn.Exec("INSERT INTO groups_to_keys (group_id, key_id) VALUES (1, 3);")
  conn.Exec("INSERT INTO groups_to_keys (group_id, key_id) VALUES (2, 1);")
  conn.Exec("INSERT INTO groups_to_keys (group_id, key_id) VALUES (2, 2);")
}

//GetPartnerDataFromKeyValue
func TestGetPartnerDataFromKeyValueHappy(t *testing.T) { 
	conn, err := pgx.Connect(dbconfig.ExtractConfig())
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		t.Error(err)
	}
	defer conn.Close()

  id, c, err := GetPartnerDataFromKeyValue("Currency", "USD", conn)

	assert.Equal(t, 1, int(id))
  assert.Equal(t, "KOH", c)
}

func TestGetPartnerDataFromKeyValueBadKey(t *testing.T) {
	conn, err := pgx.Connect(dbconfig.ExtractConfig())
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		t.Error(err)
	}
	defer conn.Close()

  CreateDB()

  id, c, err := GetPartnerDataFromKeyValue("Type of Payment", "USD", conn)
	assert.Equal(t, 0, int(id))
	//assert.NotNil(t, err)
  assert.Equal(t, "", c) //Is this right?
  //assert.NotNil(t, err)
}

func TestGetPartnerDataFromKeyValueBadValue(t *testing.T) {
	conn, err := pgx.Connect(dbconfig.ExtractConfig())
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		t.Error(err)
	}
	defer conn.Close()

  CreateDB()

  id, c, err := GetPartnerDataFromKeyValue("Currency", "CAD", conn)

  assert.Equal(t, 0, int(id))
	//assert.NotNil(t, err)
  assert.Equal(t, "", c) //Is this right?
  //assert.NotNil(t, err)
}

//GetPartnerDataByIDOrCode
func TestGetPartnerDataByIDHappy(t *testing.T) {
	conn, err := pgx.Connect(dbconfig.ExtractConfig())
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		t.Error(err)
	}
	defer conn.Close()

  CreateDB()

  id, c, err := GetPartnerDataByIDOrCode(1, "", conn)

	assert.Equal(t, 1, int(id))
  assert.Equal(t,"KOH", c)
}

func TestGetPartnerDataByCodeHappy(t *testing.T) {
	conn, err := pgx.Connect(dbconfig.ExtractConfig())
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		t.Error(err)
	}
	defer conn.Close()

  CreateDB()

  id, c, err := GetPartnerDataByIDOrCode(0, "KOH", conn)

	assert.Equal(t, 1, int(id))
  assert.Equal(t, "KOH", c)
}

func TestGetPartnerDataByIDAndCodeHappy(t *testing.T) {
	conn, err := pgx.Connect(dbconfig.ExtractConfig())
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		t.Error(err)
	}
	defer conn.Close()

  CreateDB()

  id, c, err := GetPartnerDataByIDOrCode(1, "KOH", conn)

	assert.Equal(t, 1, int(id))
  assert.Equal(t, "KOH", c)
}

func TestGetPartnerDataByIDBadID(t *testing.T) { //DONE
	conn, err := pgx.Connect(dbconfig.ExtractConfig())
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		t.Error(err)
	}
	defer conn.Close()

  CreateDB()

  id, c, err := GetPartnerDataByIDOrCode(2, "", conn)
  assert.Equal(t, 0, int(id))
	assert.NotNil(t, err)
  assert.Equal(t, "", c) //Is this right?
  assert.NotNil(t, err)
}

func TestGetPartnerDataByIDBadCode(t *testing.T) { // Do i need this one? Look at Code of func
	conn, err := pgx.Connect(dbconfig.ExtractConfig())
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		t.Error(err)
	}
	defer conn.Close()

  CreateDB()

  id, c, err := GetPartnerDataByIDOrCode(1, "DIC", conn)
  assert.Equal(t, 1, int(id)) //what do i return
  assert.Equal(t, "KOH", c) //Is this right?
}

func TestGetPartnerDataNoIDAndCode(t *testing.T) { //Do i need this one? Look at Code of func
	conn, err := pgx.Connect(dbconfig.ExtractConfig())
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		t.Error(err)
	}
	defer conn.Close()

  CreateDB()

  id, c, err := GetPartnerDataByIDOrCode(0, "", conn)

  assert.Equal(t, 0, int(id)) //what do i return
	assert.NotNil(t, err)
  assert.Equal(t, "", c) //Is this right?
  assert.NotNil(t, err)
}

func TestGetPartnerDataByCodeBadCode(t *testing.T) { //Done
	conn, err := pgx.Connect(dbconfig.ExtractConfig())
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		t.Error(err)
	}
	defer conn.Close()

  CreateDB()

  id, c, err := GetPartnerDataByIDOrCode(0, "DIC", conn)

  assert.Equal(t, 0, int(id))
	assert.NotNil(t, err)
  assert.Equal(t, "", c) //Is this right?
  assert.NotNil(t, err)
}

func TestGetPartnerDataByIDAndCodeBadID(t *testing.T) { //May be the same thing as id bad Code and Code bad ID
	conn, err := pgx.Connect(dbconfig.ExtractConfig())
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		t.Error(err)
	}
	defer conn.Close()

  CreateDB()

  id, c, err := GetPartnerDataByIDOrCode(3, "DIC", conn)

  assert.Equal(t, 0, int(id))
	assert.NotNil(t, err)
  assert.Equal(t, "", c) //Is this right?
  assert.NotNil(t, err)
}

//GetAllAttributesForPartner
func TestGetAllAttributesForPartnerHappy(t *testing.T) { //DONE
	conn, err := pgx.Connect(dbconfig.ExtractConfig())
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		t.Error(err)
	}
	defer conn.Close()

  CreateDB()

  att, err := GetAllAttributesForPartner(1, conn)

	assert.Equal(t, map[string]string{"Currency": "USD", "Type of Payment" : "Credit", "Color": "Blue"}, att)
}

func TestGetAllAttributesForPartnerBadId(t *testing.T) { //DONE
	conn, err := pgx.Connect(dbconfig.ExtractConfig())
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		t.Error(err)
	}
	defer conn.Close()

  CreateDB()

  att, err := GetAllAttributesForPartner(2, conn)

	assert.Equal(t, make(map[string]string), att) //nmeed empty map
	//assert.NotNil(t, err)
}

func TestGetGroupAttributesForPartnerHappy(t *testing.T) { //DONE
	conn, err := pgx.Connect(dbconfig.ExtractConfig())
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		t.Error(err)
	}
	defer conn.Close()

  CreateDB()

  att, err := GetGroupAttributesForPartner(1, "Money", conn)

	assert.Equal(t, map[string]string{"Currency" : "USD", "Type of Payment": "Credit"}, att)
}

func TestGetGroupAttributesForPartnerBadId(t *testing.T) { //DONE
	conn, err := pgx.Connect(dbconfig.ExtractConfig())
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		t.Error(err)
	}
	defer conn.Close()

  CreateDB()

  att, err := GetGroupAttributesForPartner(2, "Money", conn)

	assert.Equal(t, make(map[string]string), att)
	//assert.NotNil(t, err)
}

func TestGetGroupAttributesForPartnerBadGroup(t *testing.T) { //DONE
	conn, err := pgx.Connect(dbconfig.ExtractConfig())
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		t.Error(err)
	}
	defer conn.Close()

  CreateDB()

  att, err := GetGroupAttributesForPartner(1, "EDI", conn)

  assert.Equal(t, make(map[string]string), att) //need empty map
  //assert.NotNil(t, err)
}
