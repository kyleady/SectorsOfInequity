package config

type RegionConfig struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  Types []*WeightedValue
  PerterbationIds []int64
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
  newConfig.Types = StackWeightedValues(regionConfig.Types, perterbation.Types)
  newConfig.PerterbationIds = append(regionConfig.PerterbationIds, perterbation.PerterbationIds...)
  return newConfig
}

func (regionConfig *RegionConfig) Clone() *RegionConfig {
  newConfig := new(RegionConfig)
  newConfig.Name = regionConfig.Name
  newConfig.Types = make([]*WeightedValue, len(regionConfig.Types))
  copy(newConfig.Types, regionConfig.Types)
  newConfig.PerterbationIds = make([]int64, len(regionConfig.PerterbationIds))
  copy(newConfig.PerterbationIds, regionConfig.PerterbationIds)
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
      if newRegionConfig.Name == perterbationRegionConfig.Name {
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
  regionConfig.Types = FetchManyWeightedTypes(manager, regionConfig.Id, regionConfig.TableName(""), "types")
  regionConfig.PerterbationIds = FetchManyPerterbationIds(manager, regionConfig.Id, regionConfig.TableName(""), "perterbations")
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
