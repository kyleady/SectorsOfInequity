package asset

import "database/sql"
import "fmt"
import "strings"

import "screamingvortex/config"
import "screamingvortex/utilities"

type Detail struct {
  Id int64 `sql:"id"`
  ParentDetailId sql.NullInt64 `sql:"parent_detail_id"`
  childDetailGroups [][]*Detail
  Inspirations []*config.Inspiration
  InspirationTables []*config.InspirationTable
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
  client.SaveMany2ManyLinks(detail, &detail.Inspirations, "", "", "inspirations", false)
  client.SaveMany2ManyLinks(detail, &detail.InspirationTables, "", "", "inspiration_tables", false)
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

func RollDetail(inspirationTable *config.InspirationTable, perterbation *config.Perterbation) (*Detail, *config.Perterbation) {
  inspirationIds := inspirationTable.RollOnce(perterbation)
  if len(inspirationIds) == 0 {
    return nil, perterbation
  }

  return NewDetail(inspirationTable, perterbation, inspirationIds)
}

func NewDetail(inspirationTable *config.InspirationTable, perterbation *config.Perterbation, inspirationIds []int64) (*Detail, *config.Perterbation) {
  inspirations := make([]*config.Inspiration, len(inspirationIds))
  for i, inspirationId := range inspirationIds {
    inspirations[i], perterbation = perterbation.AddInspiration(inspirationId)
  }

  detail := new(Detail)
  detail.Inspirations = inspirations
  detail.InspirationTables = inspirationTable.ConstituentParts
  detail.RollsAsString = ""
  detail.ParentDetailId =  sql.NullInt64{Int64: 0, Valid: false}
  detail.childDetailGroups = [][]*Detail{}

  stackedRollGroups := []*config.InspirationTable{}
  for _, inspiration := range inspirations {
    stackedRollGroups = config.StackInspirationTables(inspiration.InspirationRolls, stackedRollGroups)
  }

  rollsAsKvPairs := []string{}
  for _, rollGroup := range stackedRollGroups {
    rollsAsKvPairs = append(rollsAsKvPairs, fmt.Sprintf("\"%s\":%d",
      rollGroup.Name,
      config.RollAll(rollGroup.CountRolls, perterbation),
    ))
  }

  detail.RollsAsString = fmt.Sprintf("{%s}", strings.Join(rollsAsKvPairs, ","))

  stackedInspirationTables := []*config.InspirationTable{}
  for _, inspiration := range inspirations {
    stackedInspirationTables = config.StackInspirationTables(inspiration.InspirationTables, stackedInspirationTables)
  }

  for _, inspirationTable := range stackedInspirationTables {
    numberOfChildDetails := config.RollAll(inspirationTable.CountRolls, perterbation)
    var childDetailGroup []*Detail
    for childDetailCount := 0; childDetailCount < numberOfChildDetails; childDetailCount++ {
      childDetail, childPerterbation := RollDetail(inspirationTable, perterbation)
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

func RollDetails(inspirationTables []*config.InspirationTable, perterbation *config.Perterbation) ([]*Detail, *config.Perterbation) {
  newPerterbation := perterbation
  details := []*Detail{}
  detail := new(Detail)
  for _, inspirationTable := range inspirationTables {
    numberOfRollsOnTheTable := inspirationTable.RollCount(newPerterbation)
    for i := 0; i < numberOfRollsOnTheTable; i++ {
      detail, newPerterbation = RollDetail(inspirationTable, newPerterbation)
      if detail != nil {
        details = append(details, detail)
      }
    }

    for _, extraInspiration := range inspirationTable.ExtraInspirations {
      extrasToAdd := config.RollAll(extraInspiration.Weights, newPerterbation)
      for i := 0; i < extrasToAdd; i++ {
        detail, newPerterbation = NewDetail(inspirationTable, newPerterbation, extraInspiration.Values)
        if detail != nil {
          details = append(details, detail)
        }
      }
    }
  }

  return details, newPerterbation
}
