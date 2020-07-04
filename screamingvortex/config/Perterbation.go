package config

import "database/sql"
import "math/rand"

import "screamingvortex/utilities"

type Perterbation struct {
  Id int64 `sql:"id"`
  SystemId sql.NullInt64 `sql:"system_id"`
  StarClusterId sql.NullInt64 `sql:"star_cluster_id"`
  RouteId sql.NullInt64 `sql:"route_id"`
  ElementId sql.NullInt64 `sql:"element_id"`
  SatelliteId sql.NullInt64 `sql:"satellite_id"`

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

func (basePerterbation *Perterbation) AddInspiration(inspirationId int64) (*Inspiration, *Perterbation) {
  inspiration := basePerterbation.Manager.GetInspiration(inspirationId)
  var newPerterbation *Perterbation
  if inspiration.PerterbationId.Valid {
    newPerterbation = basePerterbation.AddPerterbation(inspiration.PerterbationId.Int64)
  } else {
    newPerterbation = basePerterbation.Copy()
  }

  return inspiration, newPerterbation
}

func (basePerterbation *Perterbation) Copy() *Perterbation {
  return basePerterbation.addPerterbation(CreateEmptyPerterbation(nil, nil))
}

func (basePerterbation *Perterbation) AddPerterbation(perterbationId int64) *Perterbation {
  modifyingPerterbation := basePerterbation.Manager.GetPerterbation(perterbationId)
  return basePerterbation.addPerterbation(modifyingPerterbation)
}

func (basePerterbation *Perterbation) addPerterbation(modifyingPerterbation *Perterbation) *Perterbation {
  newPerterbation := new(Perterbation)
  newPerterbation.Rand = basePerterbation.Rand
  newPerterbation.Manager = basePerterbation.Manager

  newPerterbation.SystemConfig = basePerterbation.SystemConfig.AddPerterbation(modifyingPerterbation.SystemConfig)
  newPerterbation.StarClusterConfig = basePerterbation.StarClusterConfig.AddPerterbation(modifyingPerterbation.StarClusterConfig)
  newPerterbation.RouteConfig = basePerterbation.RouteConfig.AddPerterbation(modifyingPerterbation.RouteConfig)
  newPerterbation.ZoneConfigs = basePerterbation.ZoneConfigs.AddPerterbation(modifyingPerterbation.ZoneConfigs)
  newPerterbation.ElementConfig = basePerterbation.ElementConfig.AddPerterbation(modifyingPerterbation.ElementConfig)
  newPerterbation.SatelliteConfig = basePerterbation.SatelliteConfig.AddPerterbation(modifyingPerterbation.SatelliteConfig)

  return newPerterbation
}
