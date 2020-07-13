package config

import "database/sql"

type Element struct {
  Spacing []*Roll
  WeightedTypes []*WeightedValue
  SatelliteCount []*Roll
  SatelliteExtra []*WeightedValue
  TerritoryId sql.NullInt64 `sql:"territory_id"`
  TerritoryCount []*Roll
  TerritoryExtra []*WeightedValue

  TerritoryConfig *Territory
}

func (element *Element) TableName(routeType string) string {
  return "plan_config_element"
}

func (element *Element) GetId() *int64 {
  panic("GetId() not implemented. Config should not be editted.")
}

func CreateEmptyElementConfig() *Element {
  element := new(Element)
  element.Spacing = make([]*Roll, 0)
  element.WeightedTypes = make([]*WeightedValue, 0)
  element.SatelliteCount = make([]*Roll, 0)
  element.SatelliteExtra = make([]*WeightedValue, 0)
  element.TerritoryCount = make([]*Roll, 0)
  element.TerritoryExtra = make([]*WeightedValue, 0)
  element.TerritoryConfig = CreateEmptyTerritoryConfig()
  return element
}

func (element *Element) AddPerterbation(perterbation *Element) *Element {
  newElement := new(Element)
  newElement.Spacing = append(element.Spacing, perterbation.Spacing...)
  newElement.WeightedTypes = StackWeightedInspirations(element.WeightedTypes, perterbation.WeightedTypes)
  newElement.SatelliteCount = append(element.SatelliteCount, perterbation.SatelliteCount...)
  newElement.SatelliteExtra = StackWeightedInspirations(element.SatelliteExtra, perterbation.SatelliteExtra)
  newElement.TerritoryCount = append(element.TerritoryCount, perterbation.TerritoryCount...)
  newElement.TerritoryExtra = StackWeightedInspirations(element.TerritoryExtra, perterbation.TerritoryExtra)
  newElement.TerritoryConfig = element.TerritoryConfig.AddPerterbation(perterbation.TerritoryConfig)
  return newElement
}

func FetchElementConfig(manager *ConfigManager, id int64) *Element {
  element := new(Element)
  manager.Client.Fetch(element, "", id)
  element.Spacing = FetchManyRolls(manager, id, element.TableName(""), "spacing")
  element.WeightedTypes = FetchManyWeightedInspirations(manager, id, element.TableName(""), "type_inspirations")
  element.SatelliteCount = FetchManyRolls(manager, id, element.TableName(""), "satellite_count")
  element.SatelliteExtra = FetchManyWeightedInspirations(manager, id, element.TableName(""), "satellite_extra")
  element.TerritoryCount = FetchManyRolls(manager, id, element.TableName(""), "territory_count")
  element.TerritoryExtra = FetchManyWeightedInspirations(manager, id, element.TableName(""), "territory_extra")
  if element.TerritoryId.Valid {
    element.TerritoryConfig = FetchTerritoryConfig(manager, element.TerritoryId.Int64)
  } else {
    element.TerritoryConfig = CreateEmptyTerritoryConfig()
  }

  return element
}
