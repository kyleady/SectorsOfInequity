package grid

import "testing"
import "github.com/kyleady/SectorsOfInequity/screamingvortex/config"

func TestSectorRandomize(t *testing.T) {
  sectorConfig := new(config.GridConfig)
  sectorConfig.SetToDefault()

  sector := new(Sector)
  sector.Randomize(sectorConfig)

  theSurvivingLabel := -100
  systemCount := 0
  for _, row := range sector.grid {
    for _, system := range row {
      if system.Label() == system.TheUnsetLabel() {
        t.Errorf("Grid System {%d, %d} has not been properly labeled.", system.X, system.Y)
      }

      if system.Label() >= 0 {
        systemCount++
        if theSurvivingLabel >= 0 && theSurvivingLabel != system.Label() {
          t.Errorf("The sector has more than one surviving label. Label: %d, Expected: %d", system.Label(), theSurvivingLabel)
        } else if theSurvivingLabel < 0 {
          theSurvivingLabel = system.Label()
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

  if theSurvivingLabel < 0 {
    t.Errorf("No blobs were found. Possible but improbable. Each and every system would need to be empty space. Label: %d", theSurvivingLabel)
  }

  if systemCount != len(sector.Systems) {
    t.Errorf("The number of surviving systems on the grid does not equal the recorded surviving systems in the sector. %d != %d", systemCount, len(sector.Systems))
  }
}
