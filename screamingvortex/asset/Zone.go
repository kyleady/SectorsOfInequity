package asset

import "screamingvortex/config"
import "screamingvortex/utilities"

type Zone struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  Distance int `sql:"distance"`
  Elements []*Element
  Zone string
}

func (zone *Zone) TableName(zoneType string) string {
  return "plan_asset_zone"
}

func (zone *Zone) GetId() *int64 {
  return &zone.Id
}

func (zone *Zone) GetType() string {
  return zone.Zone
}

func (zone *Zone) SetName(name string) {
  zone.Name = name
}

func (zone *Zone) SaveTo(client utilities.ClientInterface) {
  client.Save(zone, "")
  zone.SaveChildren(client)
}

func (zone *Zone) SaveChildren(client utilities.ClientInterface) {
  for _, element := range zone.Elements {
    element.SaveParents(client)
  }
  client.SaveAll(&zone.Elements, "")
  client.SaveMany2ManyLinks(zone, &zone.Elements, "", "", "elements", false)
  //for _, element := range zone.Elements {
    //element.SaveChildren(client)
  //}
}

func RandomZones(perterbation *config.Perterbation, prefix string) []*Zone {
  zoneConfigs := perterbation.ZoneConfigs

  baseConfigs := []*config.Zone{}
  for _, zoneConfig := range zoneConfigs.Zones {
    if !zoneConfig.Zone.Valid {
      baseConfigs = append(baseConfigs, zoneConfig)
    }
  }

  zones := []*Zone{}
  zoneCount := 1
  for _, zoneConfig := range zoneConfigs.Zones {
    if zoneConfig.Zone.Valid {
      zoneAndBaseConfig := zoneConfig
      for _, baseConfig := range baseConfigs {
        zoneAndBaseConfig = zoneAndBaseConfig.AddPerterbation(baseConfig)
      }

      zone := new(Zone)
      zone.Zone = zoneAndBaseConfig.Zone.String
      zone.Distance =  config.RollAll(zoneAndBaseConfig.Distance, perterbation.Rand)
      newPrefix := SetNameAndGetPrefix(zone, prefix, zoneCount)
      zoneCount++
      zonePerterbation := perterbation
      for _, perterbationId := range zoneAndBaseConfig.PerterbationIds {
        zonePerterbation = zonePerterbation.AddPerterbation(perterbationId)
      }

      numberOfRandomElements := config.RollAll(zoneAndBaseConfig.ElementRolls, zonePerterbation.Rand)
      shuffledExtraIds := config.ExtraInspirationsToShuffledExtraIds(zoneAndBaseConfig.ExtraElementTypes, zonePerterbation.ElementConfig.WeightedTypes, zonePerterbation.Rand)
      numberOfExtraElements := len(shuffledExtraIds)
      numberOfElements := numberOfRandomElements + numberOfExtraElements
      numberOfRandomElementsCreated := 0
      numberOfExtraElementsCreated := 0
      distance := 0
      for i := 1; i <= numberOfElements; i++ {
        element := new(Element)
        newDistance := 0
        if numberOfExtraElements - numberOfExtraElementsCreated > 0 && zonePerterbation.Rand.Intn(numberOfElements - numberOfExtraElementsCreated - numberOfRandomElementsCreated) < numberOfExtraElements - numberOfExtraElementsCreated {
          extraInspirationIds := shuffledExtraIds[numberOfExtraElementsCreated]
          element, newDistance = NewElement(zonePerterbation, newPrefix, i, distance, extraInspirationIds)
          numberOfExtraElementsCreated++
        } else {
          element, newDistance = RandomElement(zonePerterbation, newPrefix, i, distance)
          numberOfRandomElementsCreated++
        }

        distance = newDistance
        zone.Elements = append(zone.Elements, element)
      }

      zones = append(zones, zone)
    }
  }

  return zones
}
