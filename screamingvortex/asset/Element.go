package asset

import "screamingvortex/config"
import "screamingvortex/utilities"

type Element struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  TypeId int64 `sql:"type_id"`
  Type *Detail
  Distance int `sql:"distance"`
  Satellites []*Element
  Territories []*Territory
}

func (element *Element) TableName(elementType string) string {
  return "plan_asset_element"
}

func (element *Element) GetId() *int64 {
  return &element.Id
}

func (element *Element) GetType() string {
  return element.Type.GetName()
}

func (element *Element) SetName(name string) {
  element.Name = name
}

func (element *Element) SaveTo(client utilities.ClientInterface) {
  element.SaveParents(client)
  client.Save(element, "")
  element.SaveChildren(client)
}

func (element *Element) SaveParents(client utilities.ClientInterface) {
  element.Type.SaveTo(client)
  element.TypeId = element.Type.Id
}

func (element *Element) SaveChildren(client utilities.ClientInterface) {
  for _, satellite := range element.Satellites {
    satellite.SaveParents(client)
  }

  for _, territory := range element.Territories {
    territory.SaveParents(client)
  }

  client.SaveAll(&element.Satellites, "")
  client.SaveAll(&element.Territories, "")

  client.SaveMany2ManyLinks(element, &element.Satellites, "", "", "satellites", false)
  client.SaveMany2ManyLinks(element, &element.Territories, "", "", "territories", false)

  for _, satellite := range element.Satellites {
    satellite.SaveChildren(client)
  }

  for _, territory := range element.Territories {
    territory.SaveChildren(client)
  }
}

func newElement(perterbation *config.Perterbation, prefix string, index int, distance int, elementType *Detail, isSatellite bool) (*Element, int) {
  elementConfig := perterbation.ElementConfig
  satelliteConfig := perterbation.SatelliteConfig
  if isSatellite {
    elementConfig = satelliteConfig
    satelliteConfig = nil
  }

  element := new(Element)
  element.Type = elementType
  newPrefix := SetNameAndGetPrefix(element, prefix, index)
  element.Distance = distance + config.RollAll(elementConfig.Spacing, perterbation)

  if !isSatellite {
    assetInspirationGroups := RollAssetInspirations(elementConfig.SatelliteCount, elementConfig.SatelliteExtra, satelliteConfig.WeightedTypes, perterbation)
    satelliteDistance := 0
    for i, assetInspirations := range assetInspirationGroups {
      satellite := new(Element)
      newSatelliteDistance := 0
      if assetInspirations != nil {
        satellite, newSatelliteDistance = NewSatellite(perterbation, newPrefix, i+1, satelliteDistance, assetInspirations)
      } else {
        satellite, newSatelliteDistance = RandomSatellite(perterbation, newPrefix, i+1, satelliteDistance)
      }

      satelliteDistance = newSatelliteDistance
      element.Satellites = append(element.Satellites, satellite)
    }
  }

  perterbation = perterbation.CombineElementConfigs(isSatellite)
  territoryConfig := perterbation.TerritoryConfig

  territoryInspirationGroups := RollAssetInspirations(elementConfig.TerritoryCount, elementConfig.TerritoryExtra, territoryConfig.WeightedTerritoryTypes, perterbation)
  for i, territoryInspirationGroup := range territoryInspirationGroups {
    territory := new(Territory)
    if territoryInspirationGroup != nil {
      territory = NewTerritory(perterbation, newPrefix, i+1, territoryInspirationGroup)
    } else {
      territory = RandomTerritory(perterbation, newPrefix, i+1)
    }

    element.Territories = append(element.Territories, territory)
  }

  return element, element.Distance
}

func RandomElement(perterbation *config.Perterbation, prefix string, index int, distance int) (*Element, int) {
  elementConfig := perterbation.ElementConfig
  elementType, newPerterbation := RollDetail(elementConfig.WeightedTypes, perterbation)
  return newElement(newPerterbation, prefix, index, distance, elementType, false)
}

func NewElement(perterbation *config.Perterbation, prefix string, index int, distance int, typeInspirationIds []int64) (*Element, int) {
  elementType, newPerterbation := NewDetail(typeInspirationIds, perterbation)
  return newElement(newPerterbation, prefix, index, distance, elementType, false)
}

func RandomSatellite(perterbation *config.Perterbation, prefix string, index int, distance int) (*Element, int) {
  satelliteConfig := perterbation.SatelliteConfig
  elementType, newPerterbation := RollSatelliteDetail(satelliteConfig.WeightedTypes, perterbation)
  return newElement(newPerterbation, prefix, index, distance, elementType, true)
}

func NewSatellite(perterbation *config.Perterbation, prefix string, index int, distance int, typeInspirationIds []int64) (*Element, int) {
  elementType, newPerterbation := NewSatelliteDetail(typeInspirationIds, perterbation)
  return newElement(newPerterbation, prefix, index, distance, elementType, true)
}
