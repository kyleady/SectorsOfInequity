package grid

type Route struct {
  Id int64 `sql:"id"`
  StartId int64 `sql:"start_id"`
  EndId int64 `sql:"end_id"`
  sourceSystem *System
  targetSystem *System
  X int
  Y int
}

func CreateRoute(sourceSystem *System, targetSystem *System) *Route {
  route := &Route{}
  route.sourceSystem = sourceSystem
  route.targetSystem = targetSystem
  route.X = targetSystem.X
  route.Y = targetSystem.Y

  return route
}

func (route *Route) TableName() string {
  return "config_sectorroute"
}

func (route *Route) GetId() *int64 {
  return &route.Id
}

func (route *Route) CreateReverse() *Route {
    return CreateRoute(route.targetSystem, route.sourceSystem)
}

func (route *Route) TargetSystem() *System {
  return route.targetSystem
}

func (route *Route) SourceSystem() *System {
  return route.sourceSystem
}
