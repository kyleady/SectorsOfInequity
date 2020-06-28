package config

import "database/sql"

type Inspiration struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  PerterbationId sql.NullInt64 `sql:"perterbation_id"`
  InspirationRolls []*Roll
  NestedInspirations []*NestedInspiration
}

func LoadInspiration(manager *ConfigManager, inspiration *Inspiration) {
  inspiration.InspirationRolls = FetchManyRolls(manager, inspiration.Id, inspiration.TableName(""), "rolls")
  exampleNestedInspiration := &NestedInspiration{}
  manager.Client.FetchMany(&inspiration.NestedInspirations, inspiration.Id, inspiration.TableName(""), exampleNestedInspiration.TableName(""), "inspirations", "", true)
  totalWeightedInspirations := 0
  for _, nestedInspiration := range inspiration.NestedInspirations {
    nestedInspiration.CountRolls = FetchManyRolls(manager, nestedInspiration.Id, nestedInspiration.TableName(""), "count")
    nestedInspiration.WeightedInspirations = FetchManyWeightedInspirations(manager, nestedInspiration.Id, nestedInspiration.TableName(""), "weighted_inspirations")
    totalWeightedInspirations += len(nestedInspiration.WeightedInspirations)
  }
}

func (inspiration *Inspiration) TableName(inspirationType string) string {
  return "plan_inspiration"
}

func (inspiration *Inspiration) GetId() *int64 {
  return &inspiration.Id
}

type NestedInspiration struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  CountRolls []*Roll
  WeightedInspirations []*WeightedValue
}

func (nestedInspiration *NestedInspiration) TableName(nestedInspirationType string) string {
  return "plan_inspiration_nested"
}

func (nestedInspiration *NestedInspiration) GetId() *int64 {
  panic("GetId() not implemented. Config should not be editted.")
}

func FetchManyInspirationIds(manager *ConfigManager, parentId int64, tableName string, valueName string) []int64 {
  ids := make([]int64, 0)
  exampleInspiration := new(Inspiration)
  manager.Client.FetchManyToManyChildIds(&ids, parentId, tableName, exampleInspiration.TableName(""), valueName, "", false)
  return ids
}

func (nestedInspiration *NestedInspiration) Clone() *NestedInspiration {
  newNestedInspiration := new(NestedInspiration)
  newNestedInspiration.Id = nestedInspiration.Id
  newNestedInspiration.Name = nestedInspiration.Name
  newNestedInspiration.CountRolls = make([]*Roll, len(nestedInspiration.CountRolls))
  copy(newNestedInspiration.CountRolls, nestedInspiration.CountRolls)
  newNestedInspiration.WeightedInspirations = make([]*WeightedValue, len(nestedInspiration.WeightedInspirations))
  copy(newNestedInspiration.WeightedInspirations, nestedInspiration.WeightedInspirations)
  return newNestedInspiration
}

func StackNestedInspirations(firstNestedInspirations []*NestedInspiration, secondNestedInspirations []*NestedInspiration) []*NestedInspiration {
  newNestedInspirations := make([]*NestedInspiration, len(firstNestedInspirations))
  for i, firstNestedInspiration := range firstNestedInspirations {
    newNestedInspirations[i] = firstNestedInspiration.Clone()
  }

  for _, secondNestedInspiration := range secondNestedInspirations {
    nestedInspirationStacked := false
    for _, newNestedInspiration := range newNestedInspirations {
      if newNestedInspiration.Name == secondNestedInspiration.Name {
        nestedInspirationStacked = true
        newNestedInspiration.CountRolls = append(newNestedInspiration.CountRolls, secondNestedInspiration.CountRolls...)
        newNestedInspiration.WeightedInspirations = StackWeightedValues(newNestedInspiration.WeightedInspirations, secondNestedInspiration.WeightedInspirations)
        break
      }
    }

    if !nestedInspirationStacked {
      newNestedInspirations = append(newNestedInspirations, secondNestedInspiration.Clone())
    }
  }

  return newNestedInspirations
}
