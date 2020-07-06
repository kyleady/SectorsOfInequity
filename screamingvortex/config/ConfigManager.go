package config

import "screamingvortex/utilities"
import "fmt"

type ConfigManager struct {
  cachedPerterbations map[int64]*Perterbation
  cachedInspirations map[int64]*Inspiration
  activeFlags map[string]bool
  Client *utilities.Client
}

func CreateEmptyManager(client *utilities.Client) *ConfigManager {
  manager := new(ConfigManager)
  manager.Client = client
  manager.activeFlags = make(map[string]bool)
  manager.cachedPerterbations = make(map[int64]*Perterbation)
  manager.cachedInspirations = make(map[int64]*Inspiration)
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
    inspiration := new(Inspiration)
    manager.Client.Fetch(inspiration, "", inspirationId)
    LoadInspiration(manager, inspiration)
    manager.cachedInspirations[inspirationId] = inspiration
  }

  return manager.cachedInspirations[inspirationId]
}

func (manager *ConfigManager) AddFlags(flags []string) {
  for _, flag := range flags {
    manager.activeFlags[flag] = true
  }
}

func (manager *ConfigManager) RemoveFlags(flags []string) {
  for _, flag := range flags {
    manager.activeFlags[flag] = false
  }
}

func (manager *ConfigManager) HasFlags(flags []string) bool {
  for _, flag := range flags {
    if hasFlag, _ := manager.activeFlags[flag]; !hasFlag {
      return false
    }
  }

  return true
}
