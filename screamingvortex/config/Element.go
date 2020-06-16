package config

import "screamingvortex/utilities"

type Element struct {
  Spacing []*Roll
  WeightedTypes []*WeightedValue
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
  return element
}

func (element *Element) AddPerterbation(perterbation *Element) *Element {
  newElement := new(Element)
  newElement.Spacing = append(element.Spacing, perterbation.Spacing...)
  newElement.WeightedTypes = StackWeightedValues(element.WeightedTypes, perterbation.WeightedTypes)
  return newElement
}

func LoadElementConfigFrom(client utilities.ClientInterface, id int64) *Element {
  element := new(Element)
  FetchAllRolls(client, &element.Spacing, id, element.TableName(""), "spacing")
  FetchAllWeightedInspirations(client, &element.WeightedTypes, id, element.TableName(""), "type_inspirations")
  return element
}
