package config

import "screamingvortex/utilities"

type ConfigManager struct {
  cachedPerterbations map[int64]*Perterbation
  cachedInspirations map[int64]*Inspiration
  cachedConfigTypes map[int64]*ConfigType
  Client utilities.ClientInterface
}

func CreateEmptyManager(client utilities.ClientInterface) *ConfigManager {
  manager := new(ConfigManager)
  manager.Client = client
  manager.cachedPerterbations = make(map[int64]*Perterbation)
  manager.cachedInspirations = make(map[int64]*Inspiration)
  manager.cachedConfigTypes = make(map[int64]*ConfigType)
  return manager
}

func (manager *ConfigManager) GetPerterbation(perterbationId int64) *Perterbation {
  if _, ok := manager.cachedPerterbations[perterbationId]; !ok {
    perterbation := new(Perterbation)
    manager.Client.Fetch(perterbation, "", perterbationId)
    LoadPerterbation(manager, perterbation)
    manager.cachedPerterbations[perterbationId] = perterbation
  }

  return manager.cachedPerterbations[perterbationId]
}

func (manager *ConfigManager) GetInspiration(inspirationId int64) *Inspiration {
  if _, ok := manager.cachedInspirations[inspirationId]; !ok {
    manager.cachedInspirations[inspirationId] = FetchInspiration(manager, inspirationId)
  }

  return manager.cachedInspirations[inspirationId]
}

func(manager *ConfigManager) GetConfigType(configTypeId int64) *ConfigType {
  if _, ok := manager.cachedConfigTypes[configTypeId]; !ok {
    manager.cachedConfigTypes[configTypeId] = FetchConfigType(manager, configTypeId)
  }

  return manager.cachedConfigTypes[configTypeId]
}
