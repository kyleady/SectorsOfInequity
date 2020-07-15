package config

type Inspiration struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  PerterbationIds []int64
  InspirationRolls []*InspirationTable
  InspirationTables []*InspirationTable
}

func LoadInspiration(manager *ConfigManager, inspiration *Inspiration) {
  exampleInspirationTable := &InspirationTable{}
  manager.Client.FetchMany(&inspiration.InspirationTables, inspiration.Id, inspiration.TableName(""), exampleInspirationTable.TableName(""), "inspiration_tables", "", false)
  manager.Client.FetchMany(&inspiration.InspirationRolls, inspiration.Id, inspiration.TableName(""), exampleInspirationTable.TableName(""), "roll_groups", "", false)
  inspiration.PerterbationIds = FetchManyPerterbationIds(manager, inspiration.Id, inspiration.TableName(""), "perterbations")
  for _, inspirationTable := range inspiration.InspirationTables {
    inspirationTable.CountRolls = FetchManyRolls(manager, inspirationTable.Id, inspirationTable.TableName(""), "count")
    inspirationTable.WeightedInspirations = FetchManyWeightedInspirations(manager, inspirationTable.Id, inspirationTable.TableName(""), "weighted_inspirations")
    inspirationTable.ConstituentParts = []*InspirationTable{inspirationTable}
  }

  for _, rollGroups := range inspiration.InspirationRolls {
    rollGroups.CountRolls = FetchManyRolls(manager, rollGroups.Id, rollGroups.TableName(""), "count")
    rollGroups.ConstituentParts = []*InspirationTable{rollGroups}
  }
}

func (inspiration *Inspiration) TableName(inspirationType string) string {
  return "plan_inspiration"
}

func (inspiration *Inspiration) GetId() *int64 {
  return &inspiration.Id
}

type InspirationTable struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  CountRolls []*Roll
  ConstituentParts []*InspirationTable
  WeightedInspirations []*WeightedValue
}

func (inspirationTable *InspirationTable) TableName(inspirationTableType string) string {
  return "plan_inspiration_table"
}

func (inspirationTable *InspirationTable) GetId() *int64 {
  return &inspirationTable.Id
}

func FetchManyInspirationIds(manager *ConfigManager, parentId int64, tableName string, valueName string) []int64 {
  ids := make([]int64, 0)
  exampleInspiration := new(Inspiration)
  manager.Client.FetchManyToManyChildIds(&ids, parentId, tableName, exampleInspiration.TableName(""), valueName, "", false)
  return ids
}

func (inspirationTable *InspirationTable) Clone() *InspirationTable {
  newInspirationTable := new(InspirationTable)
  newInspirationTable.Id = inspirationTable.Id
  newInspirationTable.Name = inspirationTable.Name
  newInspirationTable.CountRolls = make([]*Roll, len(inspirationTable.CountRolls))
  copy(newInspirationTable.CountRolls, inspirationTable.CountRolls)
  newInspirationTable.WeightedInspirations = make([]*WeightedValue, len(inspirationTable.WeightedInspirations))
  copy(newInspirationTable.WeightedInspirations, inspirationTable.WeightedInspirations)
  newInspirationTable.ConstituentParts = make([]*InspirationTable, len(inspirationTable.ConstituentParts))
  copy(newInspirationTable.ConstituentParts, inspirationTable.ConstituentParts)
  return newInspirationTable
}

func StackInspirationTables(firstInspirationTables []*InspirationTable, secondInspirationTables []*InspirationTable) []*InspirationTable {
  newInspirationTables := make([]*InspirationTable, len(firstInspirationTables))
  for i, firstInspirationTable := range firstInspirationTables {
    newInspirationTables[i] = firstInspirationTable.Clone()
  }

  for _, secondInspirationTable := range secondInspirationTables {
    inspirationTableStacked := false
    for _, newInspirationTable := range newInspirationTables {
      if newInspirationTable.Name == secondInspirationTable.Name {
        inspirationTableStacked = true
        newInspirationTable.CountRolls = append(newInspirationTable.CountRolls, secondInspirationTable.CountRolls...)
        newInspirationTable.WeightedInspirations = StackWeightedInspirations(newInspirationTable.WeightedInspirations, secondInspirationTable.WeightedInspirations)
        newInspirationTable.ConstituentParts = append(newInspirationTable.ConstituentParts, secondInspirationTable.ConstituentParts...)
        break
      }
    }

    if !inspirationTableStacked {
      newInspirationTables = append(newInspirationTables, secondInspirationTable.Clone())
    }
  }

  return newInspirationTables
}
