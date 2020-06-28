package config

type System struct {
  WeightedInspirations []*WeightedValue
  ExtraInspirations []*WeightedValue
  SystemFeaturesRolls []*Roll
  SystemStarClustersRolls []*Roll
}

func (system *System) TableName(systemType string) string {
  return "plan_config_system"
}

func (system *System) GetId() *int64 {
  panic("GetId() not implemented. Config should not be editted.")
}

func CreateEmptySystemConfig() *System {
  system := new(System)
  system.WeightedInspirations = make([]*WeightedValue, 0)
  system.ExtraInspirations = make([]*WeightedValue, 0)
  system.SystemFeaturesRolls = make([]*Roll, 0)
  system.SystemStarClustersRolls = make([]*Roll, 0)
  return system
}

func (system *System) AddPerterbation(perterbation *System) *System {
  newSystem := new(System)
  newSystem.SystemFeaturesRolls = append(system.SystemFeaturesRolls, perterbation.SystemFeaturesRolls...)
  newSystem.SystemStarClustersRolls = append(system.SystemStarClustersRolls, perterbation.SystemStarClustersRolls...)
  newSystem.WeightedInspirations = StackWeightedValues(system.WeightedInspirations, perterbation.WeightedInspirations)
  newSystem.ExtraInspirations = StackWeightedValues(system.ExtraInspirations, perterbation.ExtraInspirations)
  return newSystem
}

func FetchSystemConfig(manager *ConfigManager, id int64) *System {
  system := new(System)
  system.WeightedInspirations = FetchManyWeightedInspirations(manager, id, system.TableName(""), "system_feature_inspirations")
  system.ExtraInspirations = FetchManyWeightedInspirations(manager, id, system.TableName(""), "system_feature_extra")
  system.SystemFeaturesRolls = FetchManyRolls(manager, id, system.TableName(""), "system_feature_count")
  system.SystemStarClustersRolls = FetchManyRolls(manager, id, system.TableName(""), "star_cluster_count")
  return system
}
