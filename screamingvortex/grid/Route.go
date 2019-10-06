package grid

type Route struct {
  sourceSystem *System
  targetSystem *System
  X int
  Y int
}

func (route *Route) InitFromSystems(sourceSystem *System, targetSystem *System) *Route {
  route.sourceSystem = sourceSystem
  route.targetSystem = targetSystem
  route.X = targetSystem.X
  route.Y = targetSystem.Y

  return route
}

func (route *Route) CreateReverse() Route {
    newRoute := new(Route)
    return *newRoute.InitFromSystems(route.targetSystem, route.sourceSystem)
}

func (route *Route) TargetSystem() *System {
  return route.targetSystem
}

func (route *Route) SourceSystem() *System {
  return route.sourceSystem
}
