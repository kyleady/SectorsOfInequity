package asset

import "math/rand"
import "time"

import "screamingvortex/config"
import "screamingvortex/grid"
import "screamingvortex/utilities"

type Sector struct {
  Id int64 `sql:"id"`
  Name string `sql:"name"`
  Systems []*System
}

func (sector *Sector) TableName(sectorType string) string {
  return "plan_asset_sector"
}

func (sector *Sector) GetId() *int64 {
  return &sector.Id
}

func (sector *Sector) GetType() string {
  return "Sector"
}

func (sector *Sector) SaveTo(client utilities.ClientInterface) {
  client.Save(sector, "")
  sector.SaveChildren(client)
}

func (sector *Sector) SaveChildren(client utilities.ClientInterface) {
  client.SaveAll(&sector.Systems, "")
  client.SaveMany2ManyLinks(sector, &sector.Systems, "", "", "systems", false)
  for _, system := range sector.Systems {
    system.SaveChildren(client)
  }

  for _, system := range sector.Systems {
    for _, route := range system.Routes {
      for _, targetSystem := range sector.Systems {
        if route.TargetId == targetSystem.GridId {
          client.Save(&utilities.RouteToTargetSystemLink{ParentId: route.Id, ChildId: targetSystem.Id}, "")
          break
        }
      }
    }
  }
}

func RandomSector(gridSector *grid.Sector, client *utilities.Client, job *utilities.Job) *Sector {
  sector := new(Sector)
  t := time.Now()
  sector.Name = gridSector.Name + t.Format("_20060102150405")


  randSource := rand.NewSource(t.UnixNano())
  rRand := rand.New(randSource)
  emptyPerterbation := config.CreateEmptyPerterbation(client, rRand)
  for i, gridSystem := range gridSector.Systems {
    systemPerterbation := emptyPerterbation.AddPerterbation(gridSystem.RegionId)
    system := RandomSystem(systemPerterbation, "", i+1, gridSystem)

    sector.Systems = append(sector.Systems, system)
    job.Step(client)
  }

  return sector
}
