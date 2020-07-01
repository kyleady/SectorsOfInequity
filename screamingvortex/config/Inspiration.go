package config

import "database/sql"

type Inspiration struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  PerterbationId sql.NullInt64 `sql:"perterbation_id"`
  InspirationRolls []*NestedInspiration
  NestedInspirations []*NestedInspiration
}

func LoadInspiration(manager *ConfigManager, inspiration *Inspiration) {
  exampleNestedInspiration := &NestedInspiration{}
  manager.Client.FetchMany(&inspiration.NestedInspirations, inspiration.Id, inspiration.TableName(""), exampleNestedInspiration.TableName(""), "nested_inspirations", "", false)
  manager.Client.FetchMany(&inspiration.InspirationRolls, inspiration.Id, inspiration.TableName(""), exampleNestedInspiration.TableName(""), "roll_groups", "", false)
  for _, nestedInspiration := range inspiration.NestedInspirations {
    nestedInspiration.CountRolls = FetchManyRolls(manager, nestedInspiration.Id, nestedInspiration.TableName(""), "count")
    nestedInspiration.WeightedInspirations = FetchManyWeightedInspirations(manager, nestedInspiration.Id, nestedInspiration.TableName(""), "weighted_inspirations")
    nestedInspiration.ConstituentParts = []*NestedInspiration{nestedInspiration}
  }

  for _, rollGroups := range inspiration.InspirationRolls {
    rollGroups.CountRolls = FetchManyRolls(manager, rollGroups.Id, rollGroups.TableName(""), "count")
    rollGroups.ConstituentParts = []*NestedInspiration{rollGroups}
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
  ConstituentParts []*NestedInspiration
  WeightedInspirations []*WeightedValue
}

func (nestedInspiration *NestedInspiration) TableName(nestedInspirationType string) string {
  return "plan_inspiration_nested"
}

func (nestedInspiration *NestedInspiration) GetId() *int64 {
  return &nestedInspiration.Id
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
  newNestedInspiration.ConstituentParts = make([]*NestedInspiration, len(nestedInspiration.ConstituentParts))
  copy(newNestedInspiration.ConstituentParts, nestedInspiration.ConstituentParts)
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
        newNestedInspiration.WeightedInspirations = StackWeightedInspirations(newNestedInspiration.WeightedInspirations, secondNestedInspiration.WeightedInspirations)
        newNestedInspiration.ConstituentParts = append(newNestedInspiration.ConstituentParts, secondNestedInspiration.ConstituentParts...)
        break
      }
    }

    if !nestedInspirationStacked {
      newNestedInspirations = append(newNestedInspirations, secondNestedInspiration.Clone())
    }
  }

  return newNestedInspirations
}
