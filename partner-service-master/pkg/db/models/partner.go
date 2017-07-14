package models

import (
	"github.com/jackc/pgx"
	"jaxf-github.fanatics.corp/apparel/partner-service/pkg/pb"
)

type Partner struct {
	Name       pgx.NullString
	Code       pgx.NullString
	Id         pgx.NullInt32
	Attributes map[string]string
}

type Attribute struct {
	Name  pgx.NullString
	Value pgx.NullString
}

func (p Partner) Gen(attrs map[string]string) *pb.Partner {
	return &pb.Partner{
		Name:       p.Name.String,
		Code:       p.Code.String,
		Id:         p.Id.Int32,
		Attributes: attrs,
	}
}
