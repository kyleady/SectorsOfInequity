package config

import "screamingvortex/utilities"

type System struct {
  WeightedInspirations []*WeightedValue
  ExtraInspirationIds []int64
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
  system.ExtraInspirationIds = make([]int64, 0)
  system.SystemFeaturesRolls = make([]*Roll, 0)
  system.SystemStarClustersRolls = make([]*Roll, 0)
  return system
}

func (system *System) AddPerterbation(perterbation *System) *System {
  newSystem := new(System)
  newSystem.SystemFeaturesRolls = append(system.SystemFeaturesRolls, perterbation.SystemFeaturesRolls...)
  newSystem.SystemStarClustersRolls = append(system.SystemStarClustersRolls, perterbation.SystemStarClustersRolls...)
  newSystem.WeightedInspirations = StackWeightedValues(system.WeightedInspirations, perterbation.WeightedInspirations)
  newSystem.ExtraInspirationIds = append(system.ExtraInspirationIds, perterbation.ExtraInspirationIds...)
  return newSystem
}

func LoadSystemConfigFrom(client utilities.ClientInterface, id int64) *System {
  system := new(System)
  FetchAllWeightedInspirations(client, &system.WeightedInspirations, id, system.TableName(""), "system_feature_inspirations")
  FetchManyInspirationIds(client, &system.ExtraInspirationIds, id, system.TableName(""), "system_feature_extra")
  FetchAllRolls(client, &system.SystemFeaturesRolls, id, system.TableName(""), "system_feature_count")
  FetchAllRolls(client, &system.SystemStarClustersRolls, id, system.TableName(""), "star_cluster_count")
  return system
}
