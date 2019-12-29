package config

import "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"

type Grid struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  Height int `sql:"height"`
  Width int `sql:"width"`
  ConnectionRange int `sql:"connectionRange"`
  PopulationRate float64 `sql:"populationRate"`
  ConnectionRate float64 `sql:"connectionRate"`
  RangeRateMultiplier float64 `sql:"rangeRateMultiplier"`
  SmoothingFactor float64 `sql:"smoothingFactor"`
  WeightedRegions []*utilities.WeightedValue
}

func TestGrid() *Grid {
  weightedRegions := []*utilities.WeightedValue{
    &utilities.WeightedValue{1, 3, 2320},
    &utilities.WeightedValue{2, 2, 320},
    &utilities.WeightedValue{3, 4, 3499},
  }
  return &Grid{
    1234,             //Id int64 `sql:"id"`
    "test config",    //Name string `sql:"name"`
    20,               //Height int `sql:"height"`
    20,               //Width int `sql:"width"`
    3,                //ConnectionRange int `sql:"connectionRange"`
    0.75,             //PopulationRate float64 `sql:"populationRate"`
    0.53,             //RangeRateMultiplier float64 `sql:"rangeRateMultiplier"`
    0.51,             //ConnectionRate float64 `sql:"connectionRate"`
    2.0,              //SmoothingFactor float64 `sql:"smoothingFactor"`
    weightedRegions,  //WeightedRegions []*utilities.WeightedValue
  }
}

func (config *Grid) TableName(gridType string) string {
  return "plan_config_grid"
}

func (config *Grid) GetId() *int64 {
  return &config.Id
}

func LoadGridFrom(client utilities.ClientInterface, id int64) *Grid {
  gridConfig := new(Grid)
  client.Fetch(gridConfig, "", id)
  client.FetchAll(&gridConfig.WeightedRegions, "", "parent_id = ?", id)
  return gridConfig
}
