package grid

import "testing"

func createRouteTestObjs() (*Route, *System, *System) {
  locationA := Coords{3, 4}
  locationB := Coords{0, 8}

  systemA := new(System)
  systemA.Location = locationA
  systemB := new(System)
  systemB.Location = locationB

  route := new(Route)
  route.InitFromSystems(systemA, systemB)

  return route, systemA, systemB
}

func TestRouteInitializeFromSystems(t *testing.T) {
  _, systemA, systemB := createRouteTestObjs()

  route := new(Route)
  route.InitFromSystems(systemA, systemB)

  if route.SourceCoords != systemA.Location {
    t.Errorf("Route Source Coords did not match Source System Coords. Route: {%d,%d}, System: {%d, %d}",
      route.SourceCoords.X,
      route.SourceCoords.Y,
      systemA.Location.X,
      systemA.Location.Y,
    )
  }

  if route.TargetCoords != systemB.Location {
    t.Errorf("Route Target Coords did not match Target System Coords. Route: {%d,%d}, System: {%d, %d}",
      route.TargetCoords.X,
      route.TargetCoords.Y,
      systemB.Location.X,
      systemB.Location.Y,
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

  if inverseRoute.SourceCoords != route.TargetCoords {
    t.Errorf("Inverse Source Coords did not match original Target Coords. Inverse: {%d,%d}, Original: {%d, %d}",
      inverseRoute.SourceCoords.X,
      inverseRoute.SourceCoords.Y,
      route.TargetCoords.X,
      route.TargetCoords.Y,
    )
  }

  if inverseRoute.TargetCoords != route.SourceCoords {
    t.Errorf("Inverse Target Coords did not match original Source Coords. Inverse: {%d,%d}, Original: {%d, %d}",
      inverseRoute.TargetCoords.X,
      inverseRoute.TargetCoords.Y,
      route.SourceCoords.X,
      route.SourceCoords.Y,
    )
  }

  if inverseRoute.SourceSystem() != route.TargetSystem() {
    t.Errorf("Inverse SourceSystem() did not equal the original TargetSystem().")
  }

  if inverseRoute.TargetSystem() != route.SourceSystem() {
    t.Errorf("Inverse TargetSystem() did not equal the original SourceSystem().")
  }
}
