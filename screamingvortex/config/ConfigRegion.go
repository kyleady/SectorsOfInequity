package config

type RegionConfig struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  TypeId int64 `sql:"type_id"`
  Weights []*Roll
  InspirationTables []*InspirationTable
}

func (regionConfig *RegionConfig) TableName(regionConfigType string) string {
  return "plan_config_region"
}

func (regionConfig *RegionConfig) GetId() *int64 {
  return &regionConfig.Id
}

func (regionConfig *RegionConfig) AddPerterbation(perterbation *RegionConfig) *RegionConfig {
  newConfig := new(RegionConfig)
  newConfig.Name = regionConfig.Name
  newConfig.TypeId = regionConfig.TypeId
  newConfig.Weights = append(regionConfig.Weights, perterbation.Weights...)
  newConfig.InspirationTables = StackInspirationTables(regionConfig.InspirationTables, perterbation.InspirationTables)
  return newConfig
}

func (regionConfig *RegionConfig) Clone() *RegionConfig {
  newConfig := new(RegionConfig)
  newConfig.Name = regionConfig.Name
  newConfig.TypeId = regionConfig.TypeId
  newConfig.Weights = make([]*Roll, len(regionConfig.Weights))
  copy(newConfig.Weights, regionConfig.Weights)
  newConfig.InspirationTables = make([]*InspirationTable, len(regionConfig.InspirationTables))
  copy(newConfig.InspirationTables, regionConfig.InspirationTables)
  return newConfig
}

func StackRegionConfigs(firstRegionConfigs []*RegionConfig, secondRegionConfigs []*RegionConfig) []*RegionConfig {
  newRegionConfigs := make([]*RegionConfig, len(firstRegionConfigs))
  for i, regionConfig := range firstRegionConfigs {
    newRegionConfigs[i] = regionConfig.Clone()
  }

  for _, perterbationRegionConfig := range secondRegionConfigs {
    regionConfigStacked := false
    for i, newRegionConfig := range newRegionConfigs {
      if newRegionConfig.TypeId == perterbationRegionConfig.TypeId && newRegionConfig.Name == perterbationRegionConfig.Name {
        regionConfigStacked = true
        newRegionConfigs[i] = newRegionConfig.AddPerterbation(perterbationRegionConfig)
        break
      }
    }

    if !regionConfigStacked {
      newRegionConfigs = append(newRegionConfigs, perterbationRegionConfig)
    }
  }

  return newRegionConfigs
}

func FetchRegionConfig(manager *ConfigManager, id int64) *RegionConfig {
  regionConfig := new(RegionConfig)
  manager.Client.Fetch(regionConfig, "", id)
  regionConfig.FetchChildren(manager)
  return regionConfig
}

func (regionConfig *RegionConfig) FetchChildren(manager *ConfigManager) {
  regionConfig.Weights = FetchManyRolls(manager, regionConfig.Id, regionConfig.TableName(""), "weights")
  regionConfig.InspirationTables = FetchManyInspirationTables(manager, regionConfig.Id, regionConfig.TableName(""), "inspiration_tables")
}

func FetchManyRegionConfigs(manager *ConfigManager, parentId int64, tableName string, valueName string) []*RegionConfig {
  regionConfigs := make([]*RegionConfig, 0)
  regionConfigTableName := new(RegionConfig).TableName("")
  manager.Client.FetchMany(&regionConfigs, parentId, tableName, regionConfigTableName, valueName, "", false)
  for _, regionConfig := range regionConfigs {
    regionConfig.FetchChildren(manager)
  }

  return regionConfigs
}
