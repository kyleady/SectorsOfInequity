package grid

import "testing"
import "math/rand"
import "math"
import "strconv"

func createTestObjs(depth int, branches int) ([][]*System) {
  systems := make([][]*System, depth)

  var baseSystem *System
  width := 1
  for i := range systems {
    systems[i] = make([]*System, width)
    base_j := 0
    relative_j := 0
    for j := range systems[i] {
      systems[i][j] = new(System)
      systems[i][j].InitializeAt(i, j)
      if baseSystem != nil {
        systems[i][j].ConnectTo(baseSystem)
      }

      relative_j += 1
      if i == 0 || relative_j >= branches {
        relative_j = 0
        base_j += 1
        if i > 0 && base_j < len(systems[i-1]) {
          baseSystem = systems[i-1][base_j]
        } else {
          baseSystem = systems[i][0]
        }
      }
    }

    width *= branches
  }

  return systems
}

func showBlob(blob [][]*System) string {
  output := ""

  for _, row := range blob {
    for _, system := range row {
      output += strconv.Itoa(len(system.Routes)) + ","
    }

    output += "\n\n"
  }

  return output
}

func TestSystemInitializeAt(t *testing.T) {
  system := new(System)
  location := Coords{rand.Intn(100), rand.Intn(100)}

  system.InitializeAt(location.X, location.Y)


  if system.Label() != system.TheUnsetLabel() {
    t.Errorf("System's label was not initialized to the UnsetLabel. Label: %d, UnsetLabel: %d", system.Label(), system.TheUnsetLabel())
  }

  if system.Location != location {
    t.Errorf("System's location does not match the initilization coords. System: {%d,%d}, Location: {%d, %d}", system.Location.X, system.Location.Y, location.X, location.Y)
  }

  if len(system.Routes) != 0 {
    t.Errorf("System has more than zero routes. Routes: %d", len(system.Routes))
  }
}

func TestSystemSetToVoidSpace(t *testing.T) {
  system := new(System)
  system.SetToVoidSpace()

  if !system.IsVoidSpace() {
    t.Errorf("System was set to void space, but is not recognized as void space.")
  }
}

func TestSystemConnectTo(t *testing.T) {
  systemA := new(System)
  systemB := new(System)
  locationA := Coords{rand.Intn(100), rand.Intn(100)}
  locationB := Coords{rand.Intn(100), rand.Intn(100)}

  systemA.InitializeAt(locationA.X, locationA.Y)
  systemA.ConnectTo(systemA)
  if len(systemA.Routes) != 0 {
    t.Errorf("A system should not be allowed to connect to itself. Routes: %d", len(systemA.Routes))
  }

  systemA.InitializeAt(locationA.X, locationA.Y)
  systemA.ConnectTo(nil)
  if len(systemA.Routes) != 0 {
    t.Errorf("A system should not be allowed to connect to itself. Routes: %d", len(systemA.Routes))
  }

  systemA.SetToVoidSpace()
  systemB.InitializeAt(locationB.X, locationB.Y)
  systemA.ConnectTo(systemB)
  if len(systemA.Routes) != 0 || len(systemB.Routes) != 0 {
    t.Errorf("Void space should not be allowed to connect. Void Routes: %d, SystemB Routes: %d", len(systemA.Routes), len(systemB.Routes))
  }

  systemA.InitializeAt(locationA.X, locationA.Y)
  systemB.SetToVoidSpace()
  systemA.ConnectTo(systemB)
  if len(systemA.Routes) != 0 || len(systemB.Routes) != 0 {
    t.Errorf("A system should not be allowed to connect to void space. SystemA Routes: %d, Void Routes: %d", len(systemA.Routes), len(systemB.Routes))
  }

  systemA.InitializeAt(locationA.X, locationA.Y)
  systemB.InitializeAt(locationB.X, locationB.Y)
  systemA.ConnectTo(systemB)
  if len(systemA.Routes) != 1 || len(systemB.Routes) != 1 {
    t.Errorf("Connecting two systems should add a route to both systems. SystemA Routes: %d, SystemB Routes: %d", len(systemA.Routes), len(systemB.Routes))
  }

  if systemA.Routes[0].TargetCoords != locationB || systemB.Routes[0].TargetCoords != locationA {
    t.Errorf("The routes between two systems should point to each other. SystemA TargetCoords: {%d,%d}, SystemB TargetCoords: {%d,%d}, SystemA Location: {%d,%d}, SystemB Location{%d,%d}", systemA.Routes[0].TargetCoords.X, systemA.Routes[0].TargetCoords.Y, systemB.Routes[0].TargetCoords.X, systemB.Routes[0].TargetCoords.Y, systemA.Location.X, systemA.Location.Y, systemB.Location.X, systemB.Location.Y)
  }

  systemA.InitializeAt(locationA.X, locationA.Y)
  systemB.InitializeAt(locationB.X, locationB.Y)
  systemA.ConnectTo(systemB)
  systemB.ConnectTo(systemA)
  systemA.ConnectTo(systemB)
  if len(systemA.Routes) != 1 || len(systemB.Routes) != 1 {
    t.Errorf("Connecting two systems repeatedly should not add additional routes. SystemA Routes: %d, SystemB Routes: %d", len(systemA.Routes), len(systemB.Routes))
  }
}

func TestSystemLabelBlob(t *testing.T) {
  numberOfBlobs := 10
  maxDepth := 5
  maxBranches := 5

  blobs := make([][][]*System, numberOfBlobs)
  labels := make([]int, numberOfBlobs)
  depths := make([]int, numberOfBlobs)
  branches := make([]int, numberOfBlobs)
  expectedSizes := make([]int, numberOfBlobs)
  sizesByLength := make([]int, numberOfBlobs)
  randomSystem := make([]*System, numberOfBlobs)

  for i := range blobs {
    labels[i] = i
    depths[i] = rand.Intn(maxDepth-2) + 2
    branches[i] = rand.Intn(maxBranches-2) + 2
    blobs[i] = createTestObjs(depths[i], branches[i])
    expectedSizes[i] = int( ( math.Pow(float64(branches[i]), float64(depths[i])) - 1 ) / ( float64(branches[i]) - 1 ) )
    sizesByLength[i] = 0
    for j := range blobs[i] {
      sizesByLength[i] += len(blobs[i][j])
    }
  }

  for i := range blobs {
    randomSystem[i] = blobs[i][rand.Intn(depths[i]-1)+1][rand.Intn(branches[i])]
    sizeOfBlob := randomSystem[i].LabelBlob(labels[i])
    if sizeOfBlob != expectedSizes[i] {
      t.Errorf("The blob size was incorrectly calculated. Output: %d, Expected: %d, ByLength: %d\n" + showBlob(blobs[i]), sizeOfBlob, expectedSizes[i], sizesByLength[i])
    }

    errorFound := false
    for _, row := range blobs[i] {
      for _, system := range row {
        label := system.Label()
        if label != labels[i] {
          t.Errorf("The system does not have the proper label. System: %d, Expected: %d", label, labels[i])
          errorFound = true
          break
        }
      }

      if errorFound {
        break
      }
    }
  }
}

func TestSystemVoidNonMatchingLabel(t *testing.T) {
  maxDepth := 5
  maxBranches := 5
  depth := rand.Intn(maxDepth - 3) + 3
  branches := rand.Intn(maxBranches - 3) + 3
  blob := createTestObjs(depth, branches)
  label := rand.Intn(10)
  blob[rand.Intn(depth-1)+1][rand.Intn(branches-1)+1].LabelBlob(label)
  systemThatShouldVoid := blob[rand.Intn(depth-1)+1][rand.Intn(branches-1)+1]
  systemThatShouldNotVoid := blob[rand.Intn(depth-1)+1][rand.Intn(branches-1)+1]
  for systemThatShouldVoid == systemThatShouldNotVoid {
    systemThatShouldNotVoid = blob[rand.Intn(depth-1)+1][rand.Intn(branches-1)+1]
  }

  routesBeforeRemoval := len(systemThatShouldVoid.Routes)

  systemThatShouldVoid.VoidNonMatchingLabel(label + 1)
  systemThatShouldNotVoid.VoidNonMatchingLabel(label)

  if !systemThatShouldVoid.IsVoidSpace() {
    t.Errorf("The system's label is not void. System: %d, Expected: %d", systemThatShouldVoid.Label(), systemThatShouldVoid.TheVoidLabel())
  }

  if len(systemThatShouldVoid.Routes) != 0 {
    t.Errorf("The system's routes were not removed. System: %d, Expected: 0", len(systemThatShouldVoid.Routes))
  }

  if systemThatShouldNotVoid.IsVoidSpace() {
    t.Errorf("The system's label is void. System: %d, Expected: %d", systemThatShouldNotVoid.Label(), label)
  }

  if len(systemThatShouldNotVoid.Routes) == 0 {
    t.Errorf("The system's routes were removed. System: 0, Expected: %d", routesBeforeRemoval)
  }
}
