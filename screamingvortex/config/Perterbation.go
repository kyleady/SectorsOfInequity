package config

import "database/sql"
import "math/rand"
import "strings"
import "fmt"
import "regexp"

import "screamingvortex/utilities"

type Perterbation struct {
  Id int64 `sql:"id"`

  FlagsString sql.NullString `sql:"flags"`
  MutedFlagsString sql.NullString `sql:"muted_flags"`
  RequiredFlagsString sql.NullString `sql:"required_flags"`
  RejectedFlagsString sql.NullString `sql:"rejected_flags"`

  flags []string
  mutedFlags []string
  requiredFlags []string
  rejectedFlags []string

  Configs []*AssetConfig

  Manager *ConfigManager
  Rand *rand.Rand
}

func CreateEmptyPerterbation(client utilities.ClientInterface, rRand *rand.Rand) *Perterbation {
  perterbation := new(Perterbation)
  if client != nil {
    perterbation.Manager = CreateEmptyManager(client)
  }
  perterbation.Rand = rRand
  return perterbation
}

func LoadPerterbation(manager *ConfigManager, perterbation *Perterbation) {
  if perterbation.FlagsString.String != "" {
    perterbation.flags = strings.Split(perterbation.FlagsString.String, ",")
  } else {
    perterbation.flags = make([]string, 0)
  }

  if perterbation.MutedFlagsString.String != "" {
    perterbation.mutedFlags = strings.Split(perterbation.MutedFlagsString.String, ",")
  } else {
    perterbation.mutedFlags = make([]string, 0)
  }

  if perterbation.RequiredFlagsString.String != "" {
    perterbation.requiredFlags = strings.Split(perterbation.RequiredFlagsString.String, ",")
  } else {
    perterbation.requiredFlags = make([]string, 0)
  }

  if perterbation.RejectedFlagsString.String != "" {
    perterbation.rejectedFlags = strings.Split(perterbation.RejectedFlagsString.String, ",")
  } else {
    perterbation.rejectedFlags = make([]string, 0)
  }

  perterbation.Configs = FetchManyAssetConfigs(manager, perterbation.Id, perterbation.TableName(""), "configs")
}

func (perterbation *Perterbation) TableName(perterbationType string) string {
  return "plan_perterbation"
}

func (perterbation *Perterbation) GetId() *int64 {
  panic("GetId() not implemented. Config should not be editted.")
}

func (basePerterbation *Perterbation) AddInspiration(inspirationId int64) (*Inspiration, *Perterbation) {
  return basePerterbation.addInspiration(inspirationId)
}

func (basePerterbation *Perterbation) addInspiration(inspirationId int64) (*Inspiration, *Perterbation) {
  inspiration := basePerterbation.Manager.GetInspiration(inspirationId)
  newPerterbation := basePerterbation.Copy()
  for _, perterbationId := range inspiration.PerterbationIds {
    newPerterbation = newPerterbation.AddPerterbation(perterbationId)
  }

  return inspiration, newPerterbation
}

func (basePerterbation *Perterbation) Copy() *Perterbation {
  return basePerterbation.addPerterbation(CreateEmptyPerterbation(nil, nil), false)
}

func (basePerterbation *Perterbation) AddPerterbation(perterbationId int64) *Perterbation {
  modifyingPerterbation := basePerterbation.Manager.GetPerterbation(perterbationId)
  return basePerterbation.addPerterbation(modifyingPerterbation, false)
}

func (basePerterbation *Perterbation) addPerterbation(modifyingPerterbation *Perterbation, isSatellite bool) *Perterbation {
  if !basePerterbation.HasFlags(modifyingPerterbation.requiredFlags, modifyingPerterbation.rejectedFlags) {
    return basePerterbation.Copy()
  }

  newPerterbation := new(Perterbation)
  newPerterbation.Rand = basePerterbation.Rand
  newPerterbation.Manager = basePerterbation.Manager

  newPerterbation.flags = basePerterbation.CombineFlags(modifyingPerterbation)

  newPerterbation.Configs = StackAssetConfigs(basePerterbation.Configs, modifyingPerterbation.Configs)

  return newPerterbation
}

func (perterbation *Perterbation) GetConfig(typeId int64) *AssetConfig {
  for _, config := range perterbation.Configs {
    if config.TypeId == typeId {
      return config
    }
  }

  return CreateEmptyConfigAsset(typeId)
}

func FetchManyPerterbationIds(manager *ConfigManager, parentId int64, tableName string, valueName string) []int64 {
  ids := make([]int64, 0)
  examplePerterbation := new(Perterbation)
  manager.Client.FetchManyToManyChildIds(&ids, parentId, tableName, examplePerterbation.TableName(""), valueName, "", false)
  return ids
}

func (basePerterbation *Perterbation) CombineFlags(perterbation *Perterbation) []string {
  newFlags := make([]string, 0)
  mutedPatterns := make([]*regexp.Regexp, 0)
  for _, mutedFlag := range append(basePerterbation.mutedFlags, perterbation.mutedFlags...) {
    re := regexp.MustCompile(mutedFlag)
    mutedPatterns = append(mutedPatterns, re)
  }

  for _, flagToAdd := range append(basePerterbation.flags, perterbation.flags...) {
    addFlag := true
    for _, mutedPattern := range mutedPatterns {
      if mutedPattern.FindString(flagToAdd) != "" {
        addFlag = false
        break
      }
    }

    if addFlag {
      newFlags = append(newFlags, flagToAdd)
    }
  }

  return newFlags
}

func (perterbation *Perterbation) HasFlags(requiredFlags []string, rejectedFlags []string) bool {
  for _, requiredFlag := range requiredFlags {
    hasFlag := false
    requiredPattern := regexp.MustCompile(requiredFlag)
    for _, activeFlag := range perterbation.flags {
      if requiredPattern.FindString(activeFlag) != "" {
        hasFlag = true
        break
      }
    }

    if !hasFlag {
      return false
    }
  }

  for _, rejectedFlag := range rejectedFlags {
    rejectedPattern := regexp.MustCompile(rejectedFlag)
    for _, activeFlag := range perterbation.flags {
      if rejectedPattern.FindString(activeFlag) != "" {
        return false
      }
    }
  }

  return true
}

func (perterbation *Perterbation) getObject(address []*InspirationKey) (*AssetConfig, *Inspiration, *GroupConfig, *InspirationExtra, *InspirationTable) {
  if len(address) == 0 {
    panic("Empty Address!")
  }

  tmpPerterbation := perterbation
  var assetConfig *AssetConfig
  var inspiration *Inspiration
  var groupConfig *GroupConfig
  var inspirationExtra *InspirationExtra
  var inspirationTable *InspirationTable
  for _, key := range address {
    if tmpPerterbation != nil && key.Type == "AssetConfig" {
      assetConfig = tmpPerterbation.GetConfig(key.Index)
      tmpPerterbation = nil
    } else if assetConfig != nil && key.Type == "InspirationTable" {
      inspirationTable = assetConfig.GetInspirationTable(key.Key)
      assetConfig = nil
    } else if assetConfig != nil && key.Type == "GroupConfig" {
      groupConfig = assetConfig.GetGroupConfig(key.Key)
      assetConfig = nil
    } else if inspirationTable != nil && key.Type == "Inspiration" {
      inspiration = inspirationTable.GetInspiration(key.Key, key.Index != 0, perterbation)
      inspirationTable = nil
    } else if inspiration != nil && key.Type == "InspirationTable" {
      inspirationTable = inspiration.GetInspirationTable(key.Key)
      inspiration = nil
    } else if groupConfig != nil && key.Type == "InspirationExtra" {
      inspirationExtra = groupConfig.GetInspirationExtra(key.Key, key.Index)
      groupConfig = nil
    } else if inspirationExtra != nil && key.Type == "InspirationTable" {
      inspirationTable = inspirationExtra.GetInspirationTable(key.Key)
      inspirationExtra = nil
    } else {
      fmt.Print("Keys\n")
      for _, logKey := range address {
        fmt.Printf("%+v\n", logKey)
      }

      panic(fmt.Sprintf("Invalid Key: %+v", key))
    }
  }

  return assetConfig, inspiration, groupConfig, inspirationExtra, inspirationTable
}

func (perterbation *Perterbation) GetInspirationTable(address []*InspirationKey) *InspirationTable {
  _, _, _, _, inspirationTable := perterbation.getObject(address)
  return inspirationTable
}

func (perterbation *Perterbation) GetInspirationExtras(address []*InspirationKey) []*InspirationExtra {
  _, _, groupConfig, _, _ := perterbation.getObject(address)
  inspirationExtras := groupConfig.Extras
  for _, inspirationExtra := range inspirationExtras {
    inspirationExtra.SetAddress(address)
  }

  return inspirationExtras
}

func (perterbation *Perterbation) GetGroupConfigKeys(address []*InspirationKey) []*InspirationKey {
  assetConfig, _, _, _, _ := perterbation.getObject(address)
  return assetConfig.GetGroupConfigKeys()
}

func (perterbation *Perterbation) GetGroupConfig(address []*InspirationKey) *GroupConfig {
  _, _, groupConfig, _, _ := perterbation.getObject(address)
  return groupConfig
}

func (perterbation *Perterbation) GetInspirationTableNames(address []*InspirationKey) []string {
  assetConfig, inspiration, _, inspirationExtra, _ := perterbation.getObject(address)
  if inspiration != nil {
    return inspiration.GetInspirationTableNames()
  } else if inspirationExtra != nil {
    return inspirationExtra.GetInspirationTableNames()
  } else if assetConfig != nil {
      return assetConfig.GetInspirationTableNames()
  } else {
    panic("Address did not point to a list of InspirationTables.")
  }
}

func (perterbation *Perterbation) Print(indent int) {
  for i := 0; i < indent; i++ {
    fmt.Print(" ")
  }
  fmt.Print("PERTERBATION:\n")
  for i := 0; i < indent; i++ {
    fmt.Print(" ")
  }
  fmt.Printf("{Id:%d, flags:%+v}\n", perterbation.Id, perterbation.flags)

  for _, assetConfig := range perterbation.Configs {
    assetConfig.Print(indent+2)
  }
}
