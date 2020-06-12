package asset

import "screamingvortex/config"
import "screamingvortex/utilities"

type Route struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  DaysId int64 `sql:"days_id"`
  TypeId int64 `sql:"type_id"`
  Days *Detail
  Type *Detail
  TargetId int64
}

func (route *Route) TableName(routeType string) string {
  return "plan_asset_route"
}

func (route *Route) GetId() *int64 {
  return &route.Id
}

func (route *Route) GetType() string {
  return "Route"
}

func (route *Route) SetName(name string) {
  route.Name = name
}

func (route *Route) SaveTo(client utilities.ClientInterface) {
  route.SaveParents(client)
  client.Save(route, "")
}

func (route *Route) SaveParents(client utilities.ClientInterface) {
  route.Days.SaveTo(client)
  route.DaysId = route.Days.Id
  route.Type.SaveTo(client)
  route.TypeId = route.Type.Id
}

func RandomRoute(perterbation *config.Perterbation, prefix string, index int) *Route {
  routeConfig := perterbation.RouteConfig

  route := new(Route)
  SetNameAndGetPrefix(route, prefix, index)
  newPerterbation := new(config.Perterbation)
  route.Type, newPerterbation = RollDetail(routeConfig.WeightedTypes, perterbation)
  route.Days, _ = RollDetail(routeConfig.WeightedDays, newPerterbation)

  return route
}