package asset

import "math/rand"
import "strconv"

import "github.com/kyleady/SectorsOfInequity/screamingvortex/config"
import "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"

type Detail struct {
  Id int64 `sql:"id"`
  InspirationId int64 `sql:"inspiration_id"`
  ParentId int64 `sql:"asset_id"`
  RollsAsString string `sql:"rolls"`
  Type string
}

func (detail *Detail) TableName(detailType string) string {
  switch detailType {
  case config.InspirationSystemFeatureTag():
    return "plan_detail_system_feature"
  default:
    panic("Unexpected detailType.")
  }
}

func (detail *Detail) GetId() *int64 {
  return &detail.Id
}

func (detail *Detail) SaveTo(client utilities.ClientInterface) {
  client.Save(detail, detail.Type)
}


func RandomDetail(inspiration *config.Inspiration, rRand *rand.Rand) *Detail {
  detail := new(Detail)
  detail.InspirationId = inspiration.Id
  detail.RollsAsString = ""
  detail.Type = inspiration.Type

  for index, roll := range inspiration.InspirationRolls {
    if index > 0 {
      detail.RollsAsString += ","
    }
    detail.RollsAsString += strconv.Itoa(roll.Roll(rRand))

  }

  return detail
}
