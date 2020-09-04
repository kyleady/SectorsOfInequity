package config

import "database/sql"
import "strings"

type Roll struct {
  RequiredFlagsString sql.NullString `sql:"required_flags"`
  RejectedFlagsString sql.NullString `sql:"rejected_flags"`

  DiceCount int `sql:"dice_count"`
  DiceSize int `sql:"dice_size"`
  Base int `sql:"base"`
  Multiplier int `sql:"multiplier"`
  KeepHighest int `sql:"keep_highest"`
  Minimum sql.NullInt64 `sql:"minimum"`
  Maximum sql.NullInt64 `sql:"maximum"`

  requiredFlags []string
  rejectedFlags []string
}

func (roll *Roll) TableName(rollType string) string {
  return "plan_roll"
}

func (roll *Roll) GetId() *int64 {
  panic("GetId() not implemented. Config should not be editted.")
}

func (roll *Roll) Roll(perterbation *Perterbation) int {
    if !perterbation.HasFlags(roll.requiredFlags, roll.rejectedFlags) {
      return 0
    }

    result := 0
    rRand := perterbation.Rand
    diceFrequency := make(map[int]int)

    for i := 0; i < roll.DiceCount; i++ {
      diceFrequency[rRand.Intn(roll.DiceSize) + 1]++
    }

    diceToIgnore := 0
    dieStart := 1
    dieFinish := roll.DiceSize + 1
    dieIncrement := 1
    if roll.KeepHighest > 0 {
      diceToIgnore = roll.DiceCount - roll.KeepHighest
    } else if roll.KeepHighest < 0 {
      diceToIgnore = roll.DiceCount + roll.KeepHighest
      dieStart = roll.DiceSize
      dieFinish = 1 - 1
      dieIncrement = -1
    }

    for die := dieStart; die != dieFinish; die += dieIncrement {
      if diceToIgnore > diceFrequency[die] {
        diceToIgnore -= diceFrequency[die]
        diceFrequency[die] = 0
      } else if diceToIgnore > 0 {
        diceFrequency[die] -= diceToIgnore
        diceToIgnore = 0
        result += die * diceFrequency[die]
      } else {
        result += die * diceFrequency[die]
      }
    }

    result *= roll.Multiplier

    result += roll.Base

    if roll.Minimum.Valid && result < int(roll.Minimum.Int64) {
      result = int(roll.Minimum.Int64)
    }

    if roll.Maximum.Valid && result > int(roll.Maximum.Int64) {
      result = int(roll.Maximum.Int64)
    }

    return result
}

func RollAll(rolls []*Roll, perterbation *Perterbation) int {
  result := 0
  for _, roll := range rolls {
    result += roll.Roll(perterbation)
  }

  return result
}

func FetchManyRolls(manager *ConfigManager, parentId int64, tableName string, valueName string) []*Roll {
  rolls := make([]*Roll, 0)
  rollTableName := new(Roll).TableName("")
  manager.Client.FetchMany(&rolls, parentId, tableName, rollTableName, valueName, "", false)
  for _, roll := range rolls {
    if roll.RequiredFlagsString.String != "" {
      roll.requiredFlags = strings.Split(roll.RequiredFlagsString.String, ",")
    } else {
      roll.requiredFlags = make([]string, 0)
    }

    if roll.RejectedFlagsString.String != "" {
      roll.rejectedFlags = strings.Split(roll.RejectedFlagsString.String, ",")
    } else {
      roll.rejectedFlags = make([]string, 0)
    }
  }

  return rolls
}

func CreateConstantRoll(base int) *Roll {
  return &Roll{
    RequiredFlagsString: sql.NullString{Valid: false, String: ""},
    RejectedFlagsString: sql.NullString{Valid: false, String: ""},
    DiceCount: 0,
    DiceSize: 0,
    Base: base,
    Multiplier: 1,
    KeepHighest: 0,
    Minimum: sql.NullInt64{Valid: false, Int64: 0},
    Maximum: sql.NullInt64{Valid: false, Int64: 0},
    requiredFlags: []string{},
    rejectedFlags: []string{},
  }
}
