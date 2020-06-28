package config

import "database/sql"

type Zone struct {
  Id int64 `sql:"id"`
  Zone sql.NullString `sql:"zone"`
  Distance []*Roll
  PerterbationId sql.NullInt64 `sql:"perterbation_id"`
  PerterbationIds []int64
  ElementRolls []*Roll
  ExtraElementTypes []*WeightedValue
}

func (zone *Zone) TableName(zoneType string) string {
  return "plan_config_zone"
}

func (zone *Zone) GetId() *int64 {
  panic("GetId() not implemented. Config should not be editted.")
}

func (zone *Zone) AddPerterbation(perterbation *Zone) *Zone {
  newZone := new(Zone)
  newZone.Zone = zone.Zone
  newZone.Distance = append(zone.Distance, perterbation.Distance...)
  newZone.PerterbationId = sql.NullInt64{Valid: false, Int64: 0}
  newZone.PerterbationIds = append(zone.PerterbationIds, perterbation.PerterbationIds...)
  newZone.ElementRolls = append(zone.ElementRolls, perterbation.ElementRolls...)
  newZone.ExtraElementTypes = StackWeightedInspirations(zone.ExtraElementTypes, perterbation.ExtraElementTypes)
  return newZone
}

func (zone *Zone) Clone() *Zone {
  newZone := new(Zone)
  newZone.Zone = zone.Zone
  newZone.Distance = make([]*Roll, len(zone.Distance))
  copy(newZone.Distance, zone.Distance)
  newZone.ElementRolls = make([]*Roll, len(zone.ElementRolls))
  copy(newZone.ElementRolls, zone.ElementRolls)
  newZone.PerterbationIds = make([]int64, len(zone.PerterbationIds))
  copy(newZone.PerterbationIds, zone.PerterbationIds)
  newZone.ExtraElementTypes = make([]*WeightedValue, len(zone.ExtraElementTypes))
  copy(newZone.ExtraElementTypes, zone.ExtraElementTypes)
  return newZone
}

func (zone *Zone) SameZoneType(otherZone *Zone) bool {
  if zone.Zone.Valid == false || otherZone.Zone.Valid == false {
    return zone.Zone.Valid == otherZone.Zone.Valid
  }

  return zone.Zone.String == otherZone.Zone.String
}

//

type Zones struct {
  Zones []*Zone
}

func (zones *Zones) AddPerterbation(perterbation *Zones) *Zones {
  newZones := zones.Clone()
  for _, perterbationZone := range perterbation.Zones {
    zoneStacked := false
    for i, newZone := range newZones.Zones {
      if perterbationZone.SameZoneType(newZone) {
        newZones.Zones[i] = newZone.AddPerterbation(perterbationZone)
        zoneStacked = true
        break
      }
    }

    if !zoneStacked {
      newZones.Zones = append(newZones.Zones, perterbationZone.Clone())
    }
  }

  return newZones
}

func FetchZoneConfigs(manager *ConfigManager, parentId int64) *Zones {
  zones := new(Zones)
  examplePerterbation := new(Perterbation)
  exampleZone := new(Zone)
  manager.Client.FetchMany(&zones.Zones, parentId, examplePerterbation.TableName(""), exampleZone.TableName(""), "zones", "", false)
  for _, zone := range zones.Zones {
    zone.Distance = FetchManyRolls(manager, zone.Id, zone.TableName(""), "distance")
    zone.ElementRolls = FetchManyRolls(manager, zone.Id, zone.TableName(""), "element_count")
    zone.ExtraElementTypes = FetchManyWeightedInspirations(manager, zone.Id, zone.TableName(""), "element_extra")
    if zone.PerterbationId.Valid {
      zone.PerterbationIds = append(zone.PerterbationIds, zone.PerterbationId.Int64)
    } else {
      zone.PerterbationIds = make([]int64, 0)
    }
  }

  if zones.Zones == nil {
    zones.Zones = make([]*Zone, 0)
  }

  return zones
}

func (zones *Zones) Clone() *Zones {
  newZones := new(Zones)
  for _, zone := range zones.Zones {
    newZones.Zones = append(newZones.Zones, zone.Clone())
  }

  return newZones
}
