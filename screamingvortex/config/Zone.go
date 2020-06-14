package config

import "screamingvortex/utilities"

type Zone struct {
  Id int64 `sql:"id"`
  Type string `sql:"type"`
  Distance int `sql:"distance"`
  ElementRolls []*Roll
}

func (zone *Zone) TableName(zoneType string) string {
  return "plan_config_zone"
}

func (zone *Zone) GetId() *int64 {
  panic("GetId() not implemented. Config should not be editted.")
}

func (zone *Zone) AddPerterbation(perterbation *Zone) *Zone {
  newZone := new(Zone)
  newZone.Type = zone.Type
  newZone.Distance = zone.Distance + perterbation.Distance
  newZone.ElementRolls = append(zone.ElementRolls, perterbation.ElementRolls...)
  return newZone
}

func (zone *Zone) Clone() *Zone {
  newZone := new(Zone)
  newZone.Type = zone.Type
  newZone.Distance = zone.Distance
  newZone.ElementRolls = make([]*Roll, len(zone.ElementRolls))
  copy(newZone.ElementRolls, zone.ElementRolls)
  return newZone
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
      if perterbationZone.Type == newZone.Type {
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

func LoadZoneConfigsFrom(client utilities.ClientInterface, parentId int64) *Zones {
  zones := new(Zones)
  examplePerterbation := new(Perterbation)
  exampleZone := new(Zone)
  client.FetchMany(&zones.Zones, parentId, examplePerterbation.TableName(""), exampleZone.TableName(""), "zones", "", false)
  for _, zone := range zones.Zones {
    FetchAllRolls(client, &zone.ElementRolls, zone.Id, zone.TableName(""), "element_count")
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
