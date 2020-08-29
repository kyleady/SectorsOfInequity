package config

type GridConfig struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  WeightedRegions []*WeightedValue
  ConnectionTypes []*WeightedValue
  Count []*Roll
  Height []*Roll
  Width []*Roll
  ConnectionRange []*Roll
  PopulationPercent []*Roll
  ConnectionPercent []*Roll
  RangeMultiplierPercent []*Roll
  SmoothingPercent []*Roll
  PopulationDenominator int `sql:"population_denominator"`
  ConnectionDenominator int `sql:"connection_denominator"`
  RangeMultiplierDenominator int `sql:"range_multiplier_denominator"`
  SmoothingDenominator int `sql:"smoothing_denominator"`
}

func (gridConfig *GridConfig) TableName(gridConfigType string) string {
  return "plan_config_grid"
}

func (gridConfig *GridConfig) GetId() *int64 {
  panic("GetId() not implemented. Config should not be editted.")
}

func (gridConfig *GridConfig) AddPerterbation(perterbation *GridConfig) *GridConfig {
  newConfig := new(GridConfig)
  newConfig.Name = gridConfig.Name
  newConfig.WeightedRegions = StackWeightedValues(gridConfig.WeightedRegions, perterbation.WeightedRegions)
  newConfig.ConnectionTypes = StackWeightedValues(gridConfig.ConnectionTypes, perterbation.ConnectionTypes)
  newConfig.Count = append(gridConfig.Count, perterbation.Count...)
  newConfig.Height = append(gridConfig.Height, perterbation.Height...)
  newConfig.Width = append(gridConfig.Width, perterbation.Width...)
  newConfig.ConnectionRange = append(gridConfig.ConnectionRange, perterbation.ConnectionRange...)
  newConfig.PopulationPercent = append(gridConfig.PopulationPercent, perterbation.PopulationPercent...)
  newConfig.ConnectionPercent = append(gridConfig.ConnectionPercent, perterbation.ConnectionPercent...)
  newConfig.RangeMultiplierPercent = append(gridConfig.RangeMultiplierPercent, perterbation.RangeMultiplierPercent...)
  newConfig.SmoothingPercent = append(gridConfig.SmoothingPercent, perterbation.SmoothingPercent...)

  if perterbation.PopulationDenominator > gridConfig.PopulationDenominator {
    newConfig.PopulationDenominator = perterbation.PopulationDenominator
  } else {
    newConfig.PopulationDenominator = gridConfig.PopulationDenominator
  }

  if perterbation.ConnectionDenominator > gridConfig.ConnectionDenominator {
    newConfig.ConnectionDenominator = perterbation.ConnectionDenominator
  } else {
    newConfig.ConnectionDenominator = gridConfig.ConnectionDenominator
  }

  if perterbation.RangeMultiplierDenominator > gridConfig.RangeMultiplierDenominator {
    newConfig.RangeMultiplierDenominator = perterbation.RangeMultiplierDenominator
  } else {
    newConfig.RangeMultiplierDenominator = gridConfig.RangeMultiplierDenominator
  }

  if perterbation.SmoothingDenominator > gridConfig.SmoothingDenominator {
    newConfig.SmoothingDenominator = perterbation.SmoothingDenominator
  } else {
    newConfig.SmoothingDenominator = gridConfig.SmoothingDenominator
  }

  return newConfig
}

func (gridConfig *GridConfig) Clone() *GridConfig {
  newConfig := new(GridConfig)
  newConfig.Name = gridConfig.Name
  newConfig.WeightedRegions = make([]*WeightedValue, len(gridConfig.WeightedRegions))
  copy(newConfig.WeightedRegions, gridConfig.WeightedRegions)
  newConfig.ConnectionTypes = make([]*WeightedValue, len(gridConfig.ConnectionTypes))
  copy(newConfig.ConnectionTypes, gridConfig.ConnectionTypes)
  newConfig.Count = make([]*Roll, len(gridConfig.Count))
  copy(newConfig.Count, gridConfig.Count)
  newConfig.Height = make([]*Roll, len(gridConfig.Height))
  copy(newConfig.Height, gridConfig.Height)
  newConfig.Width = make([]*Roll, len(gridConfig.Width))
  copy(newConfig.Width, gridConfig.Width)
  newConfig.ConnectionRange = make([]*Roll, len(gridConfig.ConnectionRange))
  copy(newConfig.ConnectionRange, gridConfig.ConnectionRange)
  newConfig.PopulationPercent = make([]*Roll, len(gridConfig.PopulationPercent))
  copy(newConfig.PopulationPercent, gridConfig.PopulationPercent)
  newConfig.ConnectionPercent = make([]*Roll, len(gridConfig.ConnectionPercent))
  copy(newConfig.ConnectionPercent, gridConfig.ConnectionPercent)
  newConfig.RangeMultiplierPercent = make([]*Roll, len(gridConfig.RangeMultiplierPercent))
  copy(newConfig.RangeMultiplierPercent, gridConfig.RangeMultiplierPercent)
  newConfig.SmoothingPercent = make([]*Roll, len(gridConfig.SmoothingPercent))
  copy(newConfig.SmoothingPercent, gridConfig.SmoothingPercent)
  newConfig.PopulationDenominator = gridConfig.PopulationDenominator
  newConfig.ConnectionDenominator = gridConfig.ConnectionDenominator
  newConfig.RangeMultiplierDenominator = gridConfig.RangeMultiplierDenominator
  newConfig.SmoothingDenominator = gridConfig.SmoothingDenominator
  return newConfig
}

func StackGridConfigs(firstGridConfigs []*GridConfig, secondGridConfigs []*GridConfig) []*GridConfig {
  newGridConfigs := make([]*GridConfig, len(firstGridConfigs))
  for i, gridConfig := range firstGridConfigs {
    newGridConfigs[i] = gridConfig.Clone()
  }

  for _, perterbationGridConfig := range secondGridConfigs {
    gridConfigStacked := false
    for i, newGridConfig := range newGridConfigs {
      if newGridConfig.Name == perterbationGridConfig.Name {
        gridConfigStacked = true
        newGridConfigs[i] = newGridConfig.AddPerterbation(perterbationGridConfig)
        break
      }
    }

    if !gridConfigStacked {
      newGridConfigs = append(newGridConfigs, perterbationGridConfig)
    }
  }

  return newGridConfigs
}

func FetchGridConfig(manager *ConfigManager, id int64) *GridConfig {
  gridConfig := new(GridConfig)
  manager.Client.Fetch(gridConfig, "", id)
  gridConfig.FetchChildren(manager)
  return gridConfig
}

func (gridConfig *GridConfig) FetchChildren(manager *ConfigManager) {
  gridConfig.WeightedRegions = FetchManyWeightedRegions(manager, gridConfig.Id, gridConfig.TableName(""), "regions")
  gridConfig.ConnectionTypes = FetchManyWeightedTypes(manager, gridConfig.Id, gridConfig.TableName(""), "connection_types")
  gridConfig.Count = FetchManyRolls(manager, gridConfig.Id, gridConfig.TableName(""), "count")
  gridConfig.Height = FetchManyRolls(manager, gridConfig.Id, gridConfig.TableName(""), "height")
  gridConfig.Width = FetchManyRolls(manager, gridConfig.Id, gridConfig.TableName(""), "width")
  gridConfig.ConnectionRange = FetchManyRolls(manager, gridConfig.Id, gridConfig.TableName(""), "connection_range")
  gridConfig.PopulationPercent = FetchManyRolls(manager, gridConfig.Id, gridConfig.TableName(""), "population_percent")
  gridConfig.ConnectionPercent = FetchManyRolls(manager, gridConfig.Id, gridConfig.TableName(""), "connection_percent")
  gridConfig.RangeMultiplierPercent = FetchManyRolls(manager, gridConfig.Id, gridConfig.TableName(""), "range_multiplier_percent")
  gridConfig.SmoothingPercent = FetchManyRolls(manager, gridConfig.Id, gridConfig.TableName(""), "smoothing_percent")
}

func FetchManyGridConfigs(manager *ConfigManager, parentId int64, tableName string, valueName string) []*GridConfig {
  gridConfigs := make([]*GridConfig, 0)
  gridConfigTableName := new(GridConfig).TableName("")
  manager.Client.FetchMany(&gridConfigs, parentId, tableName, gridConfigTableName, valueName, "", false)
  for _, gridConfig := range gridConfigs {
    gridConfig.FetchChildren(manager)
  }

  return gridConfigs
}

func (gridConfig *GridConfig) getFraction(rollableNumerator []*Roll, denominator int, perterbation *Perterbation) float64 {
  numerator := float64(RollAll(rollableNumerator, perterbation))
  return numerator / float64(denominator)
}

func (gridConfig *GridConfig) GetPopulationFraction(perterbation *Perterbation) float64 {
  return gridConfig.getFraction(gridConfig.PopulationPercent, gridConfig.PopulationDenominator, perterbation)
}

func (gridConfig *GridConfig) GetConnectionFraction(perterbation *Perterbation) float64 {
  return gridConfig.getFraction(gridConfig.ConnectionPercent, gridConfig.ConnectionDenominator, perterbation)
}

func (gridConfig *GridConfig) GetRangeMultiplierFraction(perterbation *Perterbation) float64 {
  return gridConfig.getFraction(gridConfig.RangeMultiplierPercent, gridConfig.RangeMultiplierDenominator, perterbation)
}

func (gridConfig *GridConfig) GetSmoothingFraction(perterbation *Perterbation) float64 {
  return gridConfig.getFraction(gridConfig.SmoothingPercent, gridConfig.SmoothingDenominator, perterbation)
}
