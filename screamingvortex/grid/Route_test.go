package grid

import "testing"

func createRouteTestObjs() (*Route, *System, *System) {
  xA := 3
  yA := 4
  xB := 0
  yB := 8

  systemA := new(System)
  systemA.InitializeAt(xA, yA)
  systemB := new(System)
  systemB.InitializeAt(xB, yB)

  route := new(Route)
  route.InitFromSystems(systemA, systemB)

  return route, systemA, systemB
}

func TestRouteInitializeFromSystems(t *testing.T) {
  _, systemA, systemB := createRouteTestObjs()

  route := new(Route)
  route.InitFromSystems(systemA, systemB)

  if route.X != systemB.X || route.Y != systemB.Y {
    t.Errorf("Route Target Coords did not match Target System Coords. Route: {%d,%d}, System: {%d, %d}",
      route.X,
      route.Y,
      systemB.X,
      systemB.Y,
    )
  }
}

func TestRouteSourceSystem(t *testing.T) {
  route, systemA, _ := createRouteTestObjs()

  sourceSystem := route.SourceSystem()

  if sourceSystem != systemA {
    t.Errorf("Route SourceSystem() did not return the private Source System.")
  }
}

func TestRouteTargetSystem(t *testing.T) {
  route, _, systemB := createRouteTestObjs()

  targetSystem := route.TargetSystem()

  if targetSystem != systemB {
    t.Errorf("Route TargetSystem() did not return the private Target System.")
  }
}

func TestRouteCreateReverse(t *testing.T) {
  route, _, _ := createRouteTestObjs()

  inverseRoute := route.CreateReverse()

  if inverseRoute.SourceSystem().X != route.X || inverseRoute.SourceSystem().Y != route.Y {
    t.Errorf("Inverse Source Coords did not match original Target Coords. Inverse: {%d,%d}, Original: {%d, %d}",
      inverseRoute.SourceSystem().X,
      inverseRoute.SourceSystem().Y,
      route.X,
      route.Y,
    )
  }

  if inverseRoute.SourceSystem() != route.TargetSystem() {
    t.Errorf("Inverse SourceSystem() did not equal the original TargetSystem().")
  }

  if inverseRoute.TargetSystem() != route.SourceSystem() {
    t.Errorf("Inverse TargetSystem() did not equal the original SourceSystem().")
  }
}
