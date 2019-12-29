package grid

import "testing"
import "github.com/kyleady/SectorsOfInequity/screamingvortex/config"
import "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"
import "encoding/json"

func TestSectorRandomize(t *testing.T) {
  sectorConfig := config.TestGrid()

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
        t.Errorf("Grid System is not connected to other systems. This is possible but improbable if there is only one system in the blob. Label: %d", theSurvivingLabel)
      }

      if len(system.Routes) > 0 && system.IsVoidSpace() {
        t.Errorf("Void Space is still connected to other systems.")
      }

      for _, route := range system.Routes {
        if route.TargetSystem().Label() != theSurvivingLabel {
          t.Errorf("Route is connected to system outside the main label. TargetSystem: %v", route.TargetSystem())
        }
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

func TestLoadSectorFrom(t *testing.T) {
  client := &utilities.ClientMock{}
  client.Open()
  defer client.Close()
  client.AddTable_((&Sector{}).TableName(""))
  client.AddTable_((&System{}).TableName(""))
  client.AddTable_((&Route{}).TableName(""))

  gridConfig := config.TestGrid()
  sectorConfig := &Sector{}
  sectorConfig.Randomize(gridConfig)
  blankConfig := &Sector{}
  originalJson, _ := json.MarshalIndent(sectorConfig, "", "    ")
  blankJson, _ := json.MarshalIndent(blankConfig, "", "    ")
  if string(originalJson) == string(blankJson) {
    t.Errorf("Example sectorConfig == blank sectorConfig:\n%+v\n", sectorConfig)
  }

  sectorConfig.SaveTo(client)
  loadedConfig := LoadSectorFrom(client, sectorConfig.Id)
  sectorJson, _ := json.MarshalIndent(sectorConfig, "", "    ")
  loadedJson, _ := json.MarshalIndent(loadedConfig, "", "    ")
  if string(sectorJson) != string(loadedJson) {
    t.Errorf("Loaded sectorConfig does not equal example sectorConfig after LoadSectorFrom()\nRandomized:\n%s\nSaved:\n%s\nLoaded:\n%s", originalJson, sectorJson, loadedJson)
  }
}
