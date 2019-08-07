package config

import "time"

type GridConfig struct {
  Seed int64
  Height int
  Width int
  ConnectionRange int
  PopulationRate float64
  ConnectionRate float64
  RangeRateMultiplier float64
}

func (config *GridConfig) SetToDefault() {
  config.Seed = time.Now().UnixNano()
  config.Height = 20
  config.Width = 20
  config.ConnectionRange = 5
  config.PopulationRate = 0.5
  config.ConnectionRate = 0.4
  config.RangeRateMultiplier = 0.5
}
