package asset

import "screamingvortex/config"
import "screamingvortex/utilities"

type Element struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  TypeId int64 `sql:"type_id"`
  Type *Detail
  Distance int `sql:"distance"`
}

func (element *Element) TableName(elementType string) string {
  return "plan_asset_element"
}

func (element *Element) GetId() *int64 {
  return &element.Id
}

func (element *Element) GetType() string {
  return element.Type.Inspiration.Name
}

func (element *Element) SetName(name string) {
  element.Name = name
}

func (element *Element) SaveTo(client utilities.ClientInterface) {
  element.SaveParents(client)
  client.Save(element, "")
}

func (element *Element) SaveParents(client utilities.ClientInterface) {
  element.Type.SaveTo(client)
  element.TypeId = element.Type.Id
}

func newElement(perterbation *config.Perterbation, prefix string, index int, distance int, elementType *Detail) (*Element, int) {
  elementConfig := perterbation.ElementConfig

  element := new(Element)
  element.Type = elementType
  SetNameAndGetPrefix(element, prefix, index)
  element.Distance = distance + config.RollAll(elementConfig.Spacing, perterbation.Rand)

  return element, element.Distance
}

func RandomElement(perterbation *config.Perterbation, prefix string, index int, distance int) (*Element, int) {
  elementConfig := perterbation.ElementConfig
  elementType, newPerterbation := RollDetail(elementConfig.WeightedTypes, perterbation)
  return newElement(newPerterbation, prefix, index, distance, elementType)
}

func NewElement(perterbation *config.Perterbation, prefix string, index int, distance int, typeInspirationId int64) (*Element, int) {
  elementType, newPerterbation := NewDetail(typeInspirationId, perterbation)
  return newElement(newPerterbation, prefix, index, distance, elementType)
}
