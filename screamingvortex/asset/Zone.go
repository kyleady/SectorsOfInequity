package asset

import "screamingvortex/config"
import "screamingvortex/utilities"

type Zone struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  Distance int `sql:"distance"`
  Type string
}

func (zone *Zone) TableName(zoneType string) string {
  return "plan_asset_zone"
}

func (zone *Zone) GetId() *int64 {
  return &zone.Id
}

func (zone *Zone) GetType() string {
  return zone.Type
}

func (zone *Zone) SetName(name string) {
  zone.Name = name
}

func (zone *Zone) SaveTo(client utilities.ClientInterface) {
  client.Save(zone, "")
}

func RandomZones(perterbation *config.Perterbation, prefix string) []*Zone {
  zoneConfigs := perterbation.ZoneConfigs
  zones := []*Zone{}
  for _, zoneConfig := range zoneConfigs.Zones {
    zone := new(Zone)
    zone.Type = zoneConfig.Type
    zone.Distance = zoneConfig.Distance
    //newPrefix :=
    SetNameAndGetPrefix(zone, prefix, 1)

    zones = append(zones, zone)
  }

  return zones
}
