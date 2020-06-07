package asset

import "database/sql"
import "math/rand"
import "strconv"

import "github.com/kyleady/SectorsOfInequity/screamingvortex/config"
import "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"

type Detail struct {
  Id int64 `sql:"id"`
  ParentDetailId sql.NullInt64 `sql:"parent_detail_id"`
  childDetailGroups [][]*Detail
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
  detail.SaveChildren(client)
}

func (detail *Detail) SaveChildren(client utilities.ClientInterface) {
  for _, childDetailGroup := range detail.childDetailGroups {
    for _, childDetail := range childDetailGroup {
      childDetail.ParentDetailId.Valid = true
      childDetail.ParentDetailId.Int64 = detail.Id
      childDetail.SaveTo(client)
    }
  }
}

func newDetail(inspiration *config.Inspiration, rRand *rand.Rand) *Detail {
  detail := new(Detail)
  detail.InspirationId = inspiration.Id
  detail.RollsAsString = ""
  detail.ParentDetailId =  sql.NullInt64{Int64: 0, Valid: false}
  detail.childDetailGroups = [][]*Detail{}

  return detail
}

func RollDetail(weightedInspirations []*config.WeightedValue, perterbation *config.Perterbation) (*Detail, *config.Perterbation) {
  if len(weightedInspirations) == 0 {
    return nil, perterbation
  }

  inspirationId := config.RollWeightedValues(weightedInspirations, perterbation.Rand)
  inspiration, newPerterbation := perterbation.AddInspiration(inspirationId)
  detail := newDetail(inspiration, newPerterbation.Rand)
  for index, roll := range inspiration.InspirationRolls {
    if index > 0 {
      detail.RollsAsString += ","
    }
    detail.RollsAsString += strconv.Itoa(roll.Roll(newPerterbation.Rand))
  }

  for _, nestedInspiration := range inspiration.NestedInspirations {
    numberOfChildDetails := config.RollAll(nestedInspiration.CountRolls, newPerterbation.Rand)
    var childDetailGroup []*Detail
    for childDetailCount := 0; childDetailCount < numberOfChildDetails; childDetailCount++ {
      childDetail, childPerterbation := RollDetail(nestedInspiration.WeightedInspirations, newPerterbation)

      if childDetail != nil {
        childDetailGroup = append(childDetailGroup, childDetail)
        newPerterbation = childPerterbation
      }
    }

    if len(childDetailGroup) > 0 {
      detail.childDetailGroups = append(detail.childDetailGroups, childDetailGroup)
    }
  }

  return detail, newPerterbation
}
