package config

import "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"

type Inspiration struct {
  PerterbationId int64 `sql:"perterbation_id"`
  InspirationRolls []*Roll
}

func LoadInspirationFrom(client utilities.ClientInterface, inspirationType string, id int64) *Inspiration {
  inspiration := new(Inspiration)
  client.Fetch(inspiration, inspirationType, id)
  FetchAllRolls(client, &inspiration.InspirationRolls, inspirationType, id)
  return inspiration
}

func (inspiration *Inspiration) TableName(inspirationType string) string {
  switch inspirationType {
  case InspirationSystemFeatureTag():
    return "plan_inspiration_system_feature"
  default:
    panic("Unexpected inspirationType.")
  }
}

func (inspiration *Inspiration) GetId() *int64 {
  panic("GetId() not implemented. Config should not be editted.")
}

func InspirationSystemFeatureTag() string { return "system feature" }
