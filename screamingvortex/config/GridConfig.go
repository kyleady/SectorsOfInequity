package config

import "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"

type GridConfig struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  Height int `sql:"height"`
  Width int `sql:"width"`
  ConnectionRange int `sql:"connectionRange"`
  PopulationRate float64 `sql:"populationRate"`
  ConnectionRate float64 `sql:"connectionRate"`
  RangeRateMultiplier float64 `sql:"rangeRateMultiplier"`
  WeightedRegions []*WeightedRegion
}

func (config *GridConfig) TableName() string {
  return "config_grid"
}

func (config *GridConfig) GetId() *int64 {
  return &config.Id
}

func LoadFrom(client utilities.ClientInterface, id int64) *GridConfig {
  gridConfig := new(GridConfig)
  client.Fetch(gridConfig, id)
  client.FetchAll(&gridConfig.WeightedRegions, "id IN (SELECT weightedregion_id FROM config_grid_weightedRegions WHERE grid_id = 1)", )
  return gridConfig
}

type WeightedRegion struct {
  Id int64 `sql:"id"`
  Weight int `sql:"weight"`
  RegionId int64 `sql:"region_id"`
}

func (weightedRegion *WeightedRegion) TableName() string {
  return "config_weightedregion"
}

func (weightedRegion *WeightedRegion) GetId() *int64 {
  return &weightedRegion.Id
}

func (weightedRegion *WeightedRegion) GetWeight() int {
  return weightedRegion.Weight
}

func (weightedRegion *WeightedRegion) GetValue() interface{} {
  return weightedRegion.RegionId
}
