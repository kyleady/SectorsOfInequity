package grid

type System struct {
  Id int64 `sql:"id"`
  X int `sql:"x"`
  Y int `sql:"y"`
  SectorId int64 `sql:"sector_id"`
  RegionId int64 `sql:"region_id"`
  Routes []*Route
  blobLabel int
}

func (system *System) TableName(systemType string) string {
  return "plan_config_sector_system"
}

func (system *System) GetId() *int64 {
  return &system.Id
}

func (system *System) TheVoidLabel() int {
  return -2
}

func (system *System) TheUnsetLabel() int {
  return -1
}

func (system *System) Label() int {
  return system.blobLabel
}

func (system *System) IsVoidSpace() bool {
  return system.blobLabel == system.TheVoidLabel()
}

func (system *System) LabelIsUnset() bool {
  return system.blobLabel == system.TheUnsetLabel()
}

func (system *System) InitializeAt(i int, j int) {
  system.RegionId = 1

  system.X = i
  system.Y = j
  system.Routes = make([]*Route, 0)
  system.blobLabel = system.TheUnsetLabel()
}

func (system *System) SetToVoidSpace() {
  system.Routes = make([]*Route, 0)
  system.blobLabel = system.TheVoidLabel()
}

func (system *System) ConnectTo(targetSystem *System) {
  if targetSystem == nil || system == targetSystem {
    return
  }

  if system.IsVoidSpace() || targetSystem.IsVoidSpace() {
    return
  }

  for _, route := range system.Routes {
    if route.TargetSystem() == targetSystem {
      return
    }
  }

  newRoute := CreateRoute(system, targetSystem)
  system.Routes = append(system.Routes, newRoute)
  targetSystem.Routes = append(targetSystem.Routes, newRoute.CreateReverse())
}

func (system *System) LabelBlob(label int) int {
  if system.blobLabel == system.TheVoidLabel() || system.blobLabel >= 0 {
    return 0
  }

  system.blobLabel = label
  blobSize := 1
  for _, route := range system.Routes {
    blobSize += route.TargetSystem().LabelBlob(label)
  }

  return blobSize
}

func (system *System) VoidNonMatchingLabel(label int) {
  if system.blobLabel != label {
    system.SetToVoidSpace()
  }
}
