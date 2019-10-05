package config

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

func (config *GridConfig) SetToDefault() {
  config.Height = 20
  config.Width = 20
  config.ConnectionRange = 5
  config.PopulationRate = 0.5
  config.ConnectionRate = 0.4
  config.RangeRateMultiplier = 0.5
}

func (config *GridConfig) TableName() string {
  return "config_grid"
}

func (config *GridConfig) GetId() *int64 {
  return &config.Id
}
