package config

type GridConfig struct {
  Height int
  Width int
  ConnectionRange int
  PopulationRate float64
  ConnectionRate float64
  RangeRateMultiplier float64
}

func (config *GridConfig) SetToDefault() {
  config.Height = 20
  config.Width = 20
  config.ConnectionRange = 5
  config.PopulationRate = 0.25
  config.ConnectionRate = 0.4
  config.RangeRateMultiplier = 0.5
}
