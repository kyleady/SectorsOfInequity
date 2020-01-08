package config

import "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"

type System struct {
  WeightedInspirations []*WeightedValue
  SystemFeaturesRolls []*Roll
  SystemStarClustersRolls []*Roll
}

func (system *System) TableName(systemType string) string {
  return "plan_config_system"
}

func (system *System) GetId() *int64 {
  panic("GetId() not implemented. Config should not be editted.")
}

func (system *System) AddPerterbation(perterbation *System) *System {
  newSystem := new(System)
  return newSystem
}

func LoadSystemConfigFrom(client utilities.ClientInterface, id int64) *System {
  systemConfig := new(System)
  FetchAllWeightedValues(client, &systemConfig.WeightedInspirations, WeightedSystemInspirationTag(), id)
  FetchAllRolls(client, &systemConfig.SystemFeaturesRolls, RollSystemFeaturesTag(), id)
  FetchAllRolls(client, &systemConfig.SystemStarClustersRolls, RollSystemStarClustersTag(), id)
  return systemConfig
}
