package config

import "screamingvortex/utilities"

type Route struct {
  WeightedDays []*WeightedValue
  WeightedTypes []*WeightedValue
}

func (route *Route) TableName(routeType string) string {
  return "plan_config_route"
}

func (route *Route) GetId() *int64 {
  panic("GetId() not implemented. Config should not be editted.")
}

func CreateEmptyRouteConfig() *Route {
  route := new(Route)
  route.WeightedDays = make([]*WeightedValue, 0)
  route.WeightedTypes = make([]*WeightedValue, 0)
  return route
}

func (route *Route) AddPerterbation(perterbation *Route) *Route {
  newRoute := new(Route)
  newRoute.WeightedDays = StackWeightedValues(route.WeightedDays, perterbation.WeightedDays)
  newRoute.WeightedTypes = StackWeightedValues(route.WeightedTypes, perterbation.WeightedTypes)
  return newRoute
}

func LoadRouteConfigFrom(client utilities.ClientInterface, id int64) *Route {
  route := new(Route)
  FetchAllWeightedInspirations(client, &route.WeightedDays, id, route.TableName(""), "days_inspirations")
  FetchAllWeightedInspirations(client, &route.WeightedTypes, id, route.TableName(""), "type_inspirations")
  return route
}
