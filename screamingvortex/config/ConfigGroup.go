package config

type GroupConfig struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  Types []*WeightedValue
  Count []*Roll
  Extras []*InspirationExtra
  PerterbationIds []int64
}

func (groupConfig *GroupConfig) TableName(groupConfigType string) string {
  return "plan_config_group"
}

func (groupConfig *GroupConfig) GetId() *int64 {
  return &groupConfig.Id
}

func (groupConfig *GroupConfig) AddPerterbation(perterbation *GroupConfig) *GroupConfig {
  newConfig := new(GroupConfig)
  newConfig.Name = groupConfig.Name
  newConfig.Types = StackWeightedValues(groupConfig.Types, perterbation.Types)
  newConfig.Count = append(groupConfig.Count, perterbation.Count...)
  newConfig.Extras = StackInspirationExtras(groupConfig.Extras, perterbation.Extras)
  newConfig.PerterbationIds = append(groupConfig.PerterbationIds, perterbation.PerterbationIds...)
  return newConfig
}

func (groupConfig *GroupConfig) Clone() *GroupConfig {
  newConfig := new(GroupConfig)
  newConfig.Name = groupConfig.Name
  newConfig.Types = make([]*WeightedValue, len(groupConfig.Types))
  copy(newConfig.Types, groupConfig.Types)
  newConfig.Count = make([]*Roll, len(groupConfig.Count))
  copy(newConfig.Count, groupConfig.Count)
  newConfig.Extras = make([]*InspirationExtra, len(groupConfig.Extras))
  copy(newConfig.Extras, groupConfig.Extras)
  return newConfig
}

func StackGroupConfigs(firstGroupConfigs []*GroupConfig, secondGroupConfigs []*GroupConfig) []*GroupConfig {
  newGroupConfigs := make([]*GroupConfig, len(firstGroupConfigs))
  for i, groupConfig := range firstGroupConfigs {
    newGroupConfigs[i] = groupConfig.Clone()
  }

  for _, perterbationGroupConfig := range secondGroupConfigs {
    groupConfigStacked := false
    for i, newGroupConfig := range newGroupConfigs {
      if newGroupConfig.Name == perterbationGroupConfig.Name  {
        groupConfigStacked = true
        newGroupConfigs[i] = newGroupConfig.AddPerterbation(perterbationGroupConfig)
        break
      }
    }

    if !groupConfigStacked {
      newGroupConfigs = append(newGroupConfigs, perterbationGroupConfig)
    }
  }

  return newGroupConfigs
}

func FetchGroupConfig(manager *ConfigManager, id int64) *GroupConfig {
  groupConfig := new(GroupConfig)
  manager.Client.Fetch(groupConfig, "", id)
  groupConfig.FetchChildren(manager)
  return groupConfig
}

func (groupConfig *GroupConfig) FetchChildren(manager *ConfigManager) {
  groupConfig.Types = FetchManyWeightedTypes(manager, groupConfig.Id, groupConfig.TableName(""), "types")
  groupConfig.Count = FetchManyRolls(manager, groupConfig.Id, groupConfig.TableName(""), "count")
  groupConfig.Extras = FetchManyInspirationExtras(manager, groupConfig.Id, groupConfig.TableName(""), "extras")
  groupConfig.PerterbationIds = FetchManyPerterbationIds(manager, groupConfig.Id, groupConfig.TableName(""), "perterbations")
}

func (groupConfig *GroupConfig) GetInspirationExtra(inspirationExtraName string, inspirationExtraTypeId int64) *InspirationExtra {
  for _, inspirationExtra := range groupConfig.Extras {
    if inspirationExtra.Name == inspirationExtraName && inspirationExtra.TypeId == inspirationExtraTypeId {
      return inspirationExtra
    }
  }

  panic("GetInspirationExtra should always return a value!")
}

func FetchManyGroupConfigs(manager *ConfigManager, parentId int64, tableName string, valueName string) []*GroupConfig {
  groupConfigs := make([]*GroupConfig, 0)
  groupConfigTableName := new(GroupConfig).TableName("")
  manager.Client.FetchMany(&groupConfigs, parentId, tableName, groupConfigTableName, valueName, "", false)
  for _, groupConfig := range groupConfigs {
    groupConfig.FetchChildren(manager)
  }

  return groupConfigs
}
