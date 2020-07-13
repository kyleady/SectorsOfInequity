package config

type Territory struct {
  WeightedTerritoryTypes []*WeightedValue
}

func (territory *Territory) TableName(territoryType string) string {
  return "plan_config_territory"
}

func (territory *Territory) GetId() *int64 {
  panic("GetId() not implemented. Config should not be editted.")
}

func CreateEmptyTerritoryConfig() *Territory {
  territory := new(Territory)
  territory.WeightedTerritoryTypes = make([]*WeightedValue, 0)
  return territory
}

func (territory *Territory) AddPerterbation(perterbation *Territory) *Territory {
  newTerritory := new(Territory)
  newTerritory.WeightedTerritoryTypes = StackWeightedInspirations(territory.WeightedTerritoryTypes, perterbation.WeightedTerritoryTypes)
  return newTerritory
}

func FetchTerritoryConfig(manager *ConfigManager, id int64) *Territory {
  territory := new(Territory)
  territory.WeightedTerritoryTypes = FetchManyWeightedInspirations(manager, id, territory.TableName(""), "territory_inspirations")
  return territory
}
