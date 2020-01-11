package config

import "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"

type ConfigManager struct {
  cachedPerterbations map[string]map[int64]*Perterbation
  cachedInspirations map[string]map[int64]*Inspiration
  client *utilities.Client
}

func CreateEmptyManager(client *utilities.Client) *ConfigManager {
  manager := new(ConfigManager)
  manager.client = client
  manager.cachedPerterbations = make(map[string]map[int64]*Perterbation)
  manager.cachedInspirations = make(map[string]map[int64]*Inspiration)
  return manager
}

func (manager *ConfigManager) GetRegion(regionId int64) *Perterbation {
  return manager.GetPerterbation("sector", regionId)
}

func (manager *ConfigManager) GetPerterbation(perterbationType string, perterbationId int64) *Perterbation {
  if _, ok := manager.cachedPerterbations[perterbationType]; !ok {
    manager.cachedPerterbations[perterbationType] = make(map[int64]*Perterbation)
  }

  if _, ok := manager.cachedPerterbations[perterbationType][perterbationId]; !ok {
    manager.cachedPerterbations[perterbationType][perterbationId] = LoadPerterbationFrom(
      manager.client, perterbationType, perterbationId)
  }

  return manager.cachedPerterbations[perterbationType][perterbationId]
}

func (manager *ConfigManager) GetInspiration(inspirationType string, inspirationId int64) *Inspiration {
  if _, ok := manager.cachedInspirations[inspirationType]; !ok {
    manager.cachedInspirations[inspirationType] = make(map[int64]*Inspiration)
  }

  if _, ok := manager.cachedInspirations[inspirationType][inspirationId]; !ok {
    manager.cachedInspirations[inspirationType][inspirationId] = LoadInspirationFrom(
      manager.client, inspirationType, inspirationId)
  }

  return manager.cachedInspirations[inspirationType][inspirationId]
}
