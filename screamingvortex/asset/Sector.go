package asset

import "math/rand"
import "time"

import "github.com/kyleady/SectorsOfInequity/screamingvortex/config"
import "github.com/kyleady/SectorsOfInequity/screamingvortex/grid"
import "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"

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
  for _, system := range sector.Systems {
    client.Save(&utilities.SectorToSystemLink{ParentId: sector.Id, ChildId: system.Id}, "")
    system.SaveChildren(client)
  }
}

func RandomSector(sectorGrid *grid.Sector, client *utilities.Client, job *utilities.Job) *Sector {
  sector := new(Sector)
  t := time.Now()
  sector.Name = sectorGrid.Name + t.Format("_20060102150405")


  randSource := rand.NewSource(t.UnixNano())
  rRand := rand.New(randSource)
  emptyPerterbation := config.CreateEmptyPerterbation(client, rRand)
  for i, systemGrid := range sectorGrid.Systems {
    systemPerterbation := emptyPerterbation.AddPerterbation(systemGrid.RegionId)
    system := RandomSystem(systemPerterbation, "", i)

    sector.Systems = append(sector.Systems, system)
    job.Step(client)
  }

  return sector
}
