package config

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
  newRoute.WeightedDays = StackWeightedInspirations(route.WeightedDays, perterbation.WeightedDays)
  newRoute.WeightedStability = StackWeightedInspirations(route.WeightedStability, perterbation.WeightedStability)
  return newRoute
}

func FetchRouteConfig(manager *ConfigManager, id int64) *Route {
  route := new(Route)
  route.WeightedDays = FetchManyWeightedInspirations(manager, id, route.TableName(""), "days_inspirations")
  route.WeightedStability = FetchManyWeightedInspirations(manager, id, route.TableName(""), "stability_inspirations")
  return route
}
