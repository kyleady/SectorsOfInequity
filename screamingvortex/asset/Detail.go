package asset

import "math/rand"
import "strconv"

import "github.com/kyleady/SectorsOfInequity/screamingvortex/config"
import "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"

type Detail struct {
  Id int64 `sql:"id"`
  InspirationId int64 `sql:"inspiration_id"`
  RollsAsString string `sql:"rolls"`
}

func (detail *Detail) TableName(detailType string) string {
  return "plan_detail"
}

func (detail *Detail) GetId() *int64 {
  return &detail.Id
}

func (detail *Detail) SaveTo(client utilities.ClientInterface) {
  client.Save(detail, "")
}

func RandomDetail(inspiration *config.Inspiration, rRand *rand.Rand) *Detail {
  detail := new(Detail)
  detail.InspirationId = inspiration.Id
  detail.RollsAsString = ""
  for index, roll := range inspiration.InspirationRolls {
    if index > 0 {
      detail.RollsAsString += ","
    }
    detail.RollsAsString += strconv.Itoa(roll.Roll(rRand))
  }

  return detail
}
