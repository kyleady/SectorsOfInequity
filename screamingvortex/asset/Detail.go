package asset

import "database/sql"
import "math/rand"
import "fmt"
import "strings"

import "screamingvortex/config"
import "screamingvortex/utilities"

type Detail struct {
  Id int64 `sql:"id"`
  ParentDetailId sql.NullInt64 `sql:"parent_detail_id"`
  childDetailGroups [][]*Detail
  Inspirations []*config.Inspiration
  NestedInspirations []*config.NestedInspiration
  RollsAsString string `sql:"rolls"`
}

func (detail *Detail) TableName(detailType string) string {
  return "plan_detail"
}

func (detail *Detail) GetId() *int64 {
  return &detail.Id
}

func (detail *Detail) SaveTo(client utilities.ClientInterface) {
  fmt.Printf("\n%+v", detail)
  client.Save(detail, "")
  detail.SaveChildren(client)
}

func (detail *Detail) SaveChildren(client utilities.ClientInterface) {
  client.SaveMany2ManyLinks(detail, &detail.Inspirations, "", "", "inspirations", false)
  client.SaveMany2ManyLinks(detail, &detail.NestedInspirations, "", "", "nested_inspirations", false)
  for _, childDetailGroup := range detail.childDetailGroups {
    for _, childDetail := range childDetailGroup {
      childDetail.ParentDetailId.Valid = true
      childDetail.ParentDetailId.Int64 = detail.Id
      childDetail.SaveTo(client)
    }
  }
}

func (detail *Detail) GetName() string {
  if len(detail.Inspirations) <= 0 {
    return ""
  } else {
    return detail.Inspirations[0].Name
  }
}

func baseDetail(inspirations []*config.Inspiration, rRand *rand.Rand) *Detail {
  detail := new(Detail)
  detail.Inspirations = inspirations
  detail.RollsAsString = ""
  detail.ParentDetailId =  sql.NullInt64{Int64: 0, Valid: false}
  detail.childDetailGroups = [][]*Detail{}

  return detail
}

func NewSatelliteDetail(inspirationIds []int64, perterbation *config.Perterbation) (*Detail, *config.Perterbation) {
  return newDetail(inspirationIds, perterbation, true)
}

func NewDetail(inspirationIds []int64, perterbation *config.Perterbation) (*Detail, *config.Perterbation) {
  return newDetail(inspirationIds, perterbation, false)
}

func newDetail(inspirationIds []int64, perterbation *config.Perterbation, isSatellite bool) (*Detail, *config.Perterbation) {
  inspirations := make([]*config.Inspiration, len(inspirationIds))
  for i, inspirationId := range inspirationIds {
    if isSatellite {
      inspirations[i], perterbation = perterbation.AddSatellitedInspiration(inspirationId)
    } else {
      inspirations[i], perterbation = perterbation.AddInspiration(inspirationId)
    }
  }

  detail := baseDetail(inspirations, perterbation.Rand)
  stackedRollGroups := []*config.NestedInspiration{}
  for _, inspiration := range inspirations {
    stackedRollGroups = config.StackNestedInspirations(inspiration.InspirationRolls, stackedRollGroups)
  }

  rollsAsKvPairs := []string{}
  for _, rollGroup := range stackedRollGroups {
    rollsAsKvPairs = append(rollsAsKvPairs, fmt.Sprintf("\"%s\":%d",
      rollGroup.Name,
      config.RollAll(rollGroup.CountRolls, perterbation),
    ))
  }

  detail.RollsAsString = fmt.Sprintf("{%s}", strings.Join(rollsAsKvPairs, ","))

  stackedNestedInspirations := []*config.NestedInspiration{}
  for _, inspiration := range inspirations {
    stackedNestedInspirations = config.StackNestedInspirations(inspiration.NestedInspirations, stackedNestedInspirations)
  }

  for _, nestedInspiration := range stackedNestedInspirations {
    numberOfChildDetails := config.RollAll(nestedInspiration.CountRolls, perterbation)
    var childDetailGroup []*Detail
    for childDetailCount := 0; childDetailCount < numberOfChildDetails; childDetailCount++ {
      childDetail, childPerterbation := RollDetail(nestedInspiration.WeightedInspirations, perterbation)
      childDetail.NestedInspirations = nestedInspiration.ConstituentParts

      if childDetail != nil {
        childDetailGroup = append(childDetailGroup, childDetail)
        perterbation = childPerterbation
      }
    }

    if len(childDetailGroup) > 0 {
      detail.childDetailGroups = append(detail.childDetailGroups, childDetailGroup)
    }
  }

  return detail, perterbation
}

func RollSatelliteDetail(weightedInspirations []*config.WeightedValue, perterbation *config.Perterbation) (*Detail, *config.Perterbation) {
  if len(weightedInspirations) == 0 {
    return nil, perterbation
  }

  inspirationIds := config.RollWeightedValues(weightedInspirations, perterbation)
  return NewSatelliteDetail(inspirationIds, perterbation)
}

func RollDetail(weightedInspirations []*config.WeightedValue, perterbation *config.Perterbation) (*Detail, *config.Perterbation) {
  if len(weightedInspirations) == 0 {
    return nil, perterbation
  }

  inspirationIds := config.RollWeightedValues(weightedInspirations, perterbation)
  return NewDetail(inspirationIds, perterbation)
}
