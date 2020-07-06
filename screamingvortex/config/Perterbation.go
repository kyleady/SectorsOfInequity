package config

import "database/sql"
import "math/rand"
import "strings"

import "screamingvortex/utilities"

type Perterbation struct {
  Id int64 `sql:"id"`

  FlagsString sql.NullString `sql:"flags"`
  MutedFlagsString sql.NullString `sql:"muted_flags"`
  RequiredFlagsString sql.NullString `sql:"required_flags"`

  SystemId sql.NullInt64 `sql:"system_id"`
  StarClusterId sql.NullInt64 `sql:"star_cluster_id"`
  RouteId sql.NullInt64 `sql:"route_id"`
  ElementId sql.NullInt64 `sql:"element_id"`
  SatelliteId sql.NullInt64 `sql:"satellite_id"`

  Flags []string
  MutedFlags []string
  RequiredFlags []string

  SystemConfig *System
  StarClusterConfig *StarCluster
  RouteConfig *Route
  ZoneConfigs *Zones
  ElementConfig *Element
  SatelliteConfig *Element

  Manager *ConfigManager
  Rand *rand.Rand
}

func CreateEmptyPerterbation(client *utilities.Client, rRand *rand.Rand) *Perterbation {
  perterbation := new(Perterbation)
  if client != nil {
    perterbation.Manager = CreateEmptyManager(client)
  }
  perterbation.Rand = rRand
  perterbation.SystemConfig = CreateEmptySystemConfig()
  perterbation.StarClusterConfig = CreateEmptyStarClusterConfig()
  perterbation.RouteConfig = CreateEmptyRouteConfig()
  perterbation.ZoneConfigs = new(Zones)
  perterbation.ElementConfig = CreateEmptyElementConfig()
  perterbation.SatelliteConfig = CreateEmptyElementConfig()
  return perterbation
}

func LoadPerterbation(manager *ConfigManager, perterbation *Perterbation) {
  if perterbation.FlagsString.String != "" {
    perterbation.Flags = strings.Split(perterbation.FlagsString.String, ",")
  } else {
    perterbation.Flags = make([]string, 0)
  }

  if perterbation.MutedFlagsString.String != "" {
    perterbation.MutedFlags = strings.Split(perterbation.MutedFlagsString.String, ",")
  } else {
    perterbation.MutedFlags = make([]string, 0)
  }

  if perterbation.RequiredFlagsString.String != "" {
    perterbation.RequiredFlags = strings.Split(perterbation.RequiredFlagsString.String, ",")
  } else {
    perterbation.RequiredFlags = make([]string, 0)
  }

  if perterbation.SystemId.Valid {
    perterbation.SystemConfig = FetchSystemConfig(manager, perterbation.SystemId.Int64)
  } else {
    perterbation.SystemConfig = CreateEmptySystemConfig()
  }

  if perterbation.StarClusterId.Valid {
    perterbation.StarClusterConfig = FetchStarClusterConfig(manager, perterbation.StarClusterId.Int64)
  } else {
    perterbation.StarClusterConfig = CreateEmptyStarClusterConfig()
  }

  if perterbation.RouteId.Valid {
    perterbation.RouteConfig = FetchRouteConfig(manager, perterbation.RouteId.Int64)
  } else {
    perterbation.RouteConfig = CreateEmptyRouteConfig()
  }

  perterbation.ZoneConfigs = FetchZoneConfigs(manager, perterbation.Id)

  if perterbation.ElementId.Valid {
    perterbation.ElementConfig = FetchElementConfig(manager, perterbation.ElementId.Int64)
  } else {
    perterbation.ElementConfig = CreateEmptyElementConfig()
  }

  if perterbation.SatelliteId.Valid {
    perterbation.SatelliteConfig = FetchElementConfig(manager, perterbation.SatelliteId.Int64)
  } else {
    perterbation.SatelliteConfig = CreateEmptyElementConfig()
  }
}

func (perterbation *Perterbation) TableName(perterbationType string) string {
  return "plan_perterbation"
}

func (perterbation *Perterbation) GetId() *int64 {
  panic("GetId() not implemented. Config should not be editted.")
}

func (basePerterbation *Perterbation) AddSatellitedInspiration(inspirationId int64) (*Inspiration, *Perterbation) {
  return basePerterbation.addInspiration(inspirationId, true)
}

func (basePerterbation *Perterbation) AddInspiration(inspirationId int64) (*Inspiration, *Perterbation) {
  return basePerterbation.addInspiration(inspirationId, false)
}

func (basePerterbation *Perterbation) addInspiration(inspirationId int64, isSatellite bool) (*Inspiration, *Perterbation) {
  inspiration := basePerterbation.Manager.GetInspiration(inspirationId)
  newPerterbation := basePerterbation.Copy()
  for _, perterbationId := range inspiration.PerterbationIds {
    if isSatellite {
      newPerterbation = newPerterbation.AddSatellitePerterbation(perterbationId)
    } else {
      newPerterbation = newPerterbation.AddPerterbation(perterbationId)
    }
  }

  return inspiration, newPerterbation
}

func (basePerterbation *Perterbation) Copy() *Perterbation {
  return basePerterbation.addPerterbation(CreateEmptyPerterbation(nil, nil), false)
}

func (basePerterbation *Perterbation) AddPerterbation(perterbationId int64) *Perterbation {
  modifyingPerterbation := basePerterbation.Manager.GetPerterbation(perterbationId)
  return basePerterbation.addPerterbation(modifyingPerterbation, false)
}

func (basePerterbation *Perterbation) AddSatellitePerterbation(perterbationId int64) *Perterbation {
  modifyingPerterbation := basePerterbation.Manager.GetPerterbation(perterbationId)
  return basePerterbation.addPerterbation(modifyingPerterbation, true)
}

func (basePerterbation *Perterbation) addPerterbation(modifyingPerterbation *Perterbation, isSatellite bool) *Perterbation {
  if !basePerterbation.Manager.HasFlags(modifyingPerterbation.RequiredFlags) {
    return basePerterbation.Copy()
  }

  newPerterbation := new(Perterbation)
  newPerterbation.Rand = basePerterbation.Rand
  newPerterbation.Manager = basePerterbation.Manager

  newPerterbation.Manager.AddFlags(modifyingPerterbation.Flags)
  newPerterbation.Manager.RemoveFlags(modifyingPerterbation.MutedFlags)

  newPerterbation.SystemConfig = basePerterbation.SystemConfig.AddPerterbation(modifyingPerterbation.SystemConfig)
  newPerterbation.StarClusterConfig = basePerterbation.StarClusterConfig.AddPerterbation(modifyingPerterbation.StarClusterConfig)
  newPerterbation.RouteConfig = basePerterbation.RouteConfig.AddPerterbation(modifyingPerterbation.RouteConfig)
  newPerterbation.ZoneConfigs = basePerterbation.ZoneConfigs.AddPerterbation(modifyingPerterbation.ZoneConfigs)
  if isSatellite {
    newPerterbation.SatelliteConfig = basePerterbation.SatelliteConfig.AddPerterbation(modifyingPerterbation.ElementConfig)
  } else {
    newPerterbation.ElementConfig = basePerterbation.ElementConfig.AddPerterbation(modifyingPerterbation.ElementConfig)
  }
  newPerterbation.SatelliteConfig = basePerterbation.SatelliteConfig.AddPerterbation(modifyingPerterbation.SatelliteConfig)

  return newPerterbation
}

func FetchManyPerterbationIds(manager *ConfigManager, parentId int64, tableName string, valueName string) []int64 {
  ids := make([]int64, 0)
  examplePerterbation := new(Perterbation)
  manager.Client.FetchManyToManyChildIds(&ids, parentId, tableName, examplePerterbation.TableName(""), valueName, "", false)
  return ids
}
