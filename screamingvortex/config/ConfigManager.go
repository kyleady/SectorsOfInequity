package config

import "screamingvortex/utilities"

type ConfigManager struct {
  cachedPerterbations map[int64]*Perterbation
  cachedInspirations map[int64]*Inspiration
  client *utilities.Client
}

func CreateEmptyManager(client *utilities.Client) *ConfigManager {
  manager := new(ConfigManager)
  manager.client = client
  manager.cachedPerterbations = make(map[int64]*Perterbation)
  manager.cachedInspirations = make(map[int64]*Inspiration)
  return manager
}

func (manager *ConfigManager) GetPerterbation(perterbationId int64) *Perterbation {
  if _, ok := manager.cachedPerterbations[perterbationId]; !ok {
    manager.cachedPerterbations[perterbationId] = LoadPerterbationFrom(
      manager.client, perterbationId)
  }

  return manager.cachedPerterbations[perterbationId]
}

func (manager *ConfigManager) GetInspiration(inspirationId int64) *Inspiration {
  if _, ok := manager.cachedInspirations[inspirationId]; !ok {
    manager.cachedInspirations[inspirationId] = LoadInspirationFrom(
      manager.client, inspirationId)
  }

  return manager.cachedInspirations[inspirationId]
}
