package config

type Element struct {
  Spacing []*Roll
  WeightedTypes []*WeightedValue
  SatelliteCount []*Roll
  SatelliteExtra []*WeightedValue
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
  return element
}

func (element *Element) AddPerterbation(perterbation *Element) *Element {
  newElement := new(Element)
  newElement.Spacing = append(element.Spacing, perterbation.Spacing...)
  newElement.WeightedTypes = StackWeightedInspirations(element.WeightedTypes, perterbation.WeightedTypes)
  newElement.SatelliteCount = append(element.SatelliteCount, perterbation.SatelliteCount...)
  newElement.SatelliteExtra = StackWeightedInspirations(element.SatelliteExtra, perterbation.SatelliteExtra)
  return newElement
}

func FetchElementConfig(manager *ConfigManager, id int64) *Element {
  element := new(Element)
  element.Spacing = FetchManyRolls(manager, id, element.TableName(""), "spacing")
  element.WeightedTypes = FetchManyWeightedInspirations(manager, id, element.TableName(""), "type_inspirations")
  element.SatelliteCount = FetchManyRolls(manager, id, element.TableName(""), "satellite_count")
  element.SatelliteExtra = FetchManyWeightedInspirations(manager, id, element.TableName(""), "satellite_extra")
  return element
}
