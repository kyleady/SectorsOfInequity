package config

import "database/sql"

import "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"

type Inspiration struct {
  Id int64 `sql:"id"`
  PerterbationId sql.NullInt64 `sql:"perterbation_id"`
  InspirationRolls []*Roll
}

func LoadInspirationFrom(client utilities.ClientInterface, id int64) *Inspiration {
  inspiration := new(Inspiration)
  client.Fetch(inspiration, "", id)
  FetchAllRolls(client, &inspiration.InspirationRolls, id, inspiration.TableName(""), "rolls")
  return inspiration
}

func (inspiration *Inspiration) TableName(inspirationType string) string {
  return "plan_inspiration"
}

func (inspiration *Inspiration) GetId() *int64 {
  panic("GetId() not implemented. Config should not be editted.")
}
