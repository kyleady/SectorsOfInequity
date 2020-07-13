package asset

import "screamingvortex/config"
import "screamingvortex/utilities"

type Territory struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  TypeId int64 `sql:"type_id"`
  Type *Detail
}

func (territory *Territory) TableName(territoryType string) string {
  return "plan_asset_territory"
}

func (territory *Territory) GetId() *int64 {
  return &territory.Id
}

func (territory *Territory) GetType() string {
  return territory.Type.GetName()
}

func (territory *Territory) SetName(name string) {
  territory.Name = name
}

func (territory *Territory) SaveTo(client utilities.ClientInterface) {
  territory.SaveParents(client)
  client.Save(territory, "")
}

func (territory *Territory) SaveParents(client utilities.ClientInterface) {
  territory.Type.SaveTo(client)
  territory.TypeId = territory.Type.Id
}

func (territory *Territory) SaveChildren(client utilities.ClientInterface) {

}

func newTerritory(perterbation *config.Perterbation, prefix string, index int, territoryType *Detail) *Territory {
  territory := new(Territory)
  territory.Type = territoryType
  SetNameAndGetPrefix(territory, prefix, index)
  return territory
}

func RandomTerritory(perterbation *config.Perterbation, prefix string, index int) *Territory {
  territoryConfig := perterbation.TerritoryConfig
  territoryType, newPerterbation := RollDetail(territoryConfig.WeightedTerritoryTypes, perterbation)
  return newTerritory(newPerterbation, prefix, index, territoryType)
}

func NewTerritory(perterbation *config.Perterbation, prefix string, index int, typeInspirationIds []int64) *Territory {
  territoryType, newPerterbation := NewDetail(typeInspirationIds, perterbation)
  return newTerritory(newPerterbation, prefix, index, territoryType)
}
