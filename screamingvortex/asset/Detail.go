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

func RollDetail(inspirationTableAddress []*config.InspirationKey, perterbation *config.Perterbation) (*Detail, *config.Perterbation) {
  inspirationTable := perterbation.GetInspirationTable(inspirationTableAddress)
  inspirationIds := inspirationTable.RollOnce(perterbation)
  if len(inspirationIds) == 0 {
    return nil, perterbation
  }

  return newDetail(inspirationTableAddress, perterbation, inspirationIds, 0)
}

func NewDetail(inspirationTableAddress []*config.InspirationKey, perterbation *config.Perterbation, inspirationIds []int64) (*Detail, *config.Perterbation) {
  return newDetail(inspirationTableAddress, perterbation, inspirationIds, 1)
}

func newDetail(inspirationTableAddress []*config.InspirationKey, perterbation *config.Perterbation, inspirationIds []int64, isExtra int64) (*Detail, *config.Perterbation) {
  inspirationTable := perterbation.GetInspirationTable(inspirationTableAddress)
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
  inspirationAddress := append(inspirationTableAddress, &config.InspirationKey{Type: "Inspiration", Key: inspirations[0].Name, Index: isExtra})
  details, newPerterbation := RollDetails(inspirationAddress, perterbation)
  if len(details) > 0 {
    detail.childDetailGroups = append(detail.childDetailGroups, details)
    perterbation = newPerterbation
  }

  return detail, perterbation
}

func RollDetails(tablesAddress []*config.InspirationKey, perterbation *config.Perterbation) ([]*Detail, *config.Perterbation) {
  newPerterbation := perterbation
  details := []*Detail{}
  rolledDetails := []*Detail{}
  rolledInspirationTableNames := []string{}
  failSafe := 0
  for true {
    failSafe++
    if failSafe > 1000 {
      fmt.Println("\n\n")
      config.LogAddress(tablesAddress)
      panic(fmt.Sprintf("Not exiting!\nrolledInspirationTableNames:\n%+v", rolledInspirationTableNames))
    }

    inspirationTableAddress := GetUnrolledInspirationTableAddress(tablesAddress, rolledInspirationTableNames, newPerterbation)
    if len(inspirationTableAddress) == 0 {
      break
    }

    rolledDetails, newPerterbation = RollInspirationTable(inspirationTableAddress, newPerterbation)
    details = append(details, rolledDetails...)
    rolledInspirationTableNames = append(rolledInspirationTableNames, inspirationTableAddress[len(inspirationTableAddress)-1].Key)
  }

  return details, newPerterbation
}

func GetUnrolledInspirationTableAddress(tablesAddress []*config.InspirationKey, rolledInspirationTableNames []string, perterbation *config.Perterbation) []*config.InspirationKey {
  inspirationTableNames := perterbation.GetInspirationTableNames(tablesAddress)
  unrolledTableName := ""
  inspirationTableFound := false
  for _, inspirationTableName := range inspirationTableNames {
    hasBeenRolled := false
    for _, rolledInspirationTableName := range rolledInspirationTableNames {
      if inspirationTableName == rolledInspirationTableName {
        hasBeenRolled = true
        break
      }
    }

    if !hasBeenRolled {
      unrolledTableName = inspirationTableName
      inspirationTableFound = true
      break
    }
  }

  if !inspirationTableFound {
    return []*config.InspirationKey{}
  }

  inspirationTableAddress := append(tablesAddress, &config.InspirationKey{
    Type: "InspirationTable",
    Key: unrolledTableName,
  })

  return inspirationTableAddress
}

func RollInspirationTable(tableAddress []*config.InspirationKey, perterbation *config.Perterbation) ([]*Detail, *config.Perterbation) {
  details := []*Detail{}
  detail := new(Detail)
  newPerterbation := perterbation
  inspirationTable := perterbation.GetInspirationTable(tableAddress)
  numberOfRollsOnTheTable := inspirationTable.RollCount(newPerterbation)
  for i := 0; i < numberOfRollsOnTheTable; i++ {
    detail, newPerterbation = RollDetail(tableAddress, newPerterbation)
    if detail != nil {
      details = append(details, detail)
    }
  }

  for _, extraInspiration := range inspirationTable.ExtraInspirations {
    extrasToAdd := config.RollAll(extraInspiration.Weights, newPerterbation)
    for i := 0; i < extrasToAdd; i++ {
      detail, newPerterbation = NewDetail(tableAddress, newPerterbation, extraInspiration.Values)
      if detail != nil {
        details = append(details, detail)
      }
    }
  }

  return details, newPerterbation
}
