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
  return gridConfig
}

func ExampleGridConfig() *GridConfig{
  return &GridConfig{
    314159265,      //Id int64 `sql:"id"`
    "Test Config",  //Name string `sql:"name"`
    20,             //Height int `sql:"height"`
    21,             //Width int `sql:"width"`
    5,              //ConnectionRange int `sql:"connectionRange"`
    0.24,           //PopulationRate float64 `sql:"populationRate"`
    0.25,           //ConnectionRate float64 `sql:"connectionRate"`
    0.26,           //RangeRateMultiplier float64 `sql:"rangeRateMultiplier"`
  }
}
