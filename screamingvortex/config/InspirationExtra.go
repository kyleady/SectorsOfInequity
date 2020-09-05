package config

type InspirationExtra struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  TypeId int64 `sql:"type_id"`
  CountRolls []*Roll
  InspirationTables []*WeightedValue
  Address []*InspirationKey
}

func (inspirationExtra *InspirationExtra) TableName(inspirationExtraType string) string {
  return "plan_inspiration_extra"
}

func (inspirationExtra *InspirationExtra) GetId() *int64 {
  return &inspirationExtra.Id
}

func (inspirationExtra *InspirationExtra) AddPerterbation(perterbationInspirationExtra *InspirationExtra) *InspirationExtra {
  newInspirationExtra := new(InspirationExtra)
  newInspirationExtra.Name = inspirationExtra.Name
  newInspirationExtra.TypeId = inspirationExtra.TypeId
  newInspirationExtra.CountRolls = append(inspirationExtra.CountRolls, perterbationInspirationExtra.CountRolls...)
  newInspirationExtra.InspirationTables = StackWeightedValues(inspirationExtra.InspirationTables, perterbationInspirationExtra.InspirationTables)
  return newInspirationExtra
}

func (inspirationExtra *InspirationExtra) Clone() *InspirationExtra {
  newInspirationExtra := new(InspirationExtra)
  newInspirationExtra.Id = inspirationExtra.Id
  newInspirationExtra.TypeId = inspirationExtra.TypeId
  newInspirationExtra.CountRolls = make([]*Roll, len(inspirationExtra.CountRolls))
  copy(newInspirationExtra.CountRolls, inspirationExtra.CountRolls)
  newInspirationExtra.InspirationTables = make([]*WeightedValue, len(inspirationExtra.InspirationTables))
  copy(newInspirationExtra.InspirationTables, inspirationExtra.InspirationTables)
  return newInspirationExtra
}

func (inspirationExtra *InspirationExtra) FetchChildren(manager *ConfigManager) {
  inspirationExtra.CountRolls = FetchManyRolls(manager, inspirationExtra.Id, inspirationExtra.TableName(""), "count")
  inspirationExtra.InspirationTables = FetchManyWeightedTables(manager, inspirationExtra.Id, inspirationExtra.TableName(""), "inspiration_tables")
}

func (inspirationExtra *InspirationExtra) SetAddress(address []*InspirationKey) {
  inspirationExtra.Address = append(address, &InspirationKey{Type: "InspirationExtra", Key: inspirationExtra.Name, Index: inspirationExtra.TypeId})
}

func (inspirationExtra *InspirationExtra) GetInspirationTableNames(perterbation *Perterbation) []string {
  tableNames := []string{}
  SortWeightedValues(inspirationExtra.InspirationTables, perterbation)
  for _, inspirationTable := range inspirationExtra.InspirationTables {
    tableNames = append(tableNames, inspirationTable.ValueName)
  }

  return tableNames
}

func (inspirationExtra *InspirationExtra) GetInspirationTable(inspirationTableName string, perterbation *Perterbation) *InspirationTable {
  for _, inspirationTable := range inspirationExtra.InspirationTables {
    if inspirationTable.ValueName == inspirationTableName {
      return perterbation.Manager.GetInspirationTable(inspirationTable.Values)
    }
  }

  panic("GetInspirationTable should always return a value!")
}

func StackInspirationExtras(firstInspirationExtras []*InspirationExtra, secondInspirationExtras []*InspirationExtra) []*InspirationExtra {
  newInspirationExtras := make([]*InspirationExtra, len(firstInspirationExtras))
  for i, firstInspirationExtra := range firstInspirationExtras {
    newInspirationExtras[i] = firstInspirationExtra.Clone()
  }

  for _, secondInspirationExtra := range secondInspirationExtras {
    inspirationExtraStacked := false
    for i, newInspirationExtra := range newInspirationExtras {
      if newInspirationExtra.Name == secondInspirationExtra.Name && newInspirationExtra.TypeId == secondInspirationExtra.TypeId {
        inspirationExtraStacked = true
        newInspirationExtras[i] = newInspirationExtra.AddPerterbation(secondInspirationExtra)
        break
      }
    }

    if !inspirationExtraStacked {
      newInspirationExtras = append(newInspirationExtras, secondInspirationExtra.Clone())
    }
  }

  return newInspirationExtras
}

func FetchManyInspirationExtras(manager *ConfigManager, parentId int64, tableName string, valueName string) []*InspirationExtra {
  inspirationExtras := make([]*InspirationExtra, 0)
  inspirationExtraTableName := new(InspirationExtra).TableName("")
  manager.Client.FetchMany(&inspirationExtras, parentId, tableName, inspirationExtraTableName, valueName, "", false)
  for _, inspirationExtra := range inspirationExtras {
    inspirationExtra.FetchChildren(manager)
  }

  return inspirationExtras
}
