package config

import "testing"
import "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"

func TestLoadFrom(t *testing.T) {
  client := &utilities.ClientMock{}
  client.Open()
  defer client.Close()
  client.AddTable_((&GridConfig{}).TableName())

  gridConfig := ExampleGridConfig()
  blankConfig := &GridConfig{}
  if *gridConfig == *blankConfig {
    t.Errorf("Example gridConfig == blank gridConfig:\n%+v\n", gridConfig)
  }

  client.Save(gridConfig)
  loadedConfig := LoadFrom(client, gridConfig.Id)
  if *gridConfig != *loadedConfig {
    t.Errorf("Loaded gridConfig does not equal example gridConfig after LoadFrom():\n%+v\n", loadedConfig)
  }
}
