package config

import "database/sql"

import "screamingvortex/utilities"

type Inspiration struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  PerterbationId sql.NullInt64 `sql:"perterbation_id"`
  InspirationRolls []*Roll
  NestedInspirations []*NestedInspiration
}

func LoadInspirationFrom(client utilities.ClientInterface, id int64) *Inspiration {
  inspiration := new(Inspiration)
  client.Fetch(inspiration, "", id)
  FetchAllRolls(client, &inspiration.InspirationRolls, id, inspiration.TableName(""), "rolls")
  exampleNestedInspiration := &NestedInspiration{}
  client.FetchMany(&inspiration.NestedInspirations, id, exampleNestedInspiration.TableName(""), inspiration.TableName(""), "inspirations", "", true)
  totalWeightedInspirations := 0
  for _, nestedInspiration := range inspiration.NestedInspirations {
    FetchAllRolls(client, &nestedInspiration.CountRolls, nestedInspiration.Id, nestedInspiration.TableName(""), "count")
    FetchAllWeightedInspirations(client, &nestedInspiration.WeightedInspirations, nestedInspiration.Id, nestedInspiration.TableName(""), "weighted_inspirations")
    totalWeightedInspirations += len(nestedInspiration.WeightedInspirations)
  }

  return inspiration
}

func FetchManyInspirationIds(client utilities.ClientInterface, ids *[]int64, parentId int64, tableName string, valueName string) {
  exampleInspiration := new(Inspiration)
  client.FetchManyToManyChildIds(ids, parentId, tableName, exampleInspiration.TableName(""), valueName, "", false)
}

func (inspiration *Inspiration) TableName(inspirationType string) string {
  return "plan_inspiration"
}

func (inspiration *Inspiration) GetId() *int64 {
  panic("GetId() not implemented. Config should not be editted.")
}

type NestedInspiration struct {
  Id int64 `sql:"id"`
  CountRolls []*Roll
  WeightedInspirations []*WeightedValue
}

func (nestedInspiration *NestedInspiration) TableName(nestedInspirationType string) string {
  return "plan_inspiration_nested"
}

func (nestedInspiration *NestedInspiration) GetId() *int64 {
  panic("GetId() not implemented. Config should not be editted.")
}
