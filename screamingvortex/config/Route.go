package config

import "screamingvortex/utilities"

type Route struct {
  WeightedDays []*WeightedValue
  WeightedStability []*WeightedValue
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
  route.WeightedStability = make([]*WeightedValue, 0)
  return route
}

func (route *Route) AddPerterbation(perterbation *Route) *Route {
  newRoute := new(Route)
  newRoute.WeightedDays = StackWeightedValues(route.WeightedDays, perterbation.WeightedDays)
  newRoute.WeightedStability = StackWeightedValues(route.WeightedStability, perterbation.WeightedStability)
  return newRoute
}

func LoadRouteConfigFrom(client utilities.ClientInterface, id int64) *Route {
  route := new(Route)
  FetchAllWeightedInspirations(client, &route.WeightedDays, id, route.TableName(""), "days_inspirations")
  FetchAllWeightedInspirations(client, &route.WeightedStability, id, route.TableName(""), "stability_inspirations")
  return route
}
