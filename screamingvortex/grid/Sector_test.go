package grid

import "testing"
import "github.com/kyleady/SectorsOfInequity/screamingvortex/config"

func TestSectorRandomize(t *testing.T) {
  sectorConfig := new(config.GridConfig)
  sectorConfig.SetToDefault()

  sector := new(Sector)
  sector.Randomize(sectorConfig)

  theOneLabel := -100
  for _, row := range sector.Grid {
    for _, system := range row {
      if system.Label() == system.TheUnsetLabel() {
        t.Errorf("Grid System {%d, %d} has not been properly labeled.", system.Location.X, system.Location.Y)
      }

      if system.Label() >= 0 {
        if theOneLabel >= 0 && theOneLabel != system.Label() {
          t.Errorf("The sector has more than one surviving label. Label: %d, Expected: %d", system.Location.X, system.Location.Y, system.Label(), theOneLabel)
        } else if theOneLabel < 0 {
          theOneLabel = system.Label()
        }
      }

      if len(system.Routes) == 0 && !system.IsVoidSpace() {
        t.Errorf("Grid System is not connected to other systems. This is possible but improbable if there is only one system in the blob.")
      }

      if len(system.Routes) > 0 && system.IsVoidSpace() {
        t.Errorf("Void Space is still connected to other systems.")
      }
    }
  }

  if theOneLabel < 0 {
    t.Errorf("No blobs were found. Possible but improbable. Each and every system would need to be empty space. Label: %d", theOneLabel)
  }
}
