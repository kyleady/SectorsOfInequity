package config

import "math/rand"

import "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"

type Roll struct {
  DiceCount int `sql:"dice_count"`
  DiceSize int `sql:"dice_size"`
  Base int `sql:"base"`
  Multiplier int `sql:"multiplier"`
  KeepHighest int `sql:"keep_highest"`
}

func LoadRollFrom(client utilities.ClientInterface, rollType string, id int64) *Roll {
  roll := new(Roll)
  client.Fetch(roll, rollType, id)

  return roll
}

func (roll *Roll) TableName(rollType string) string {
  return "plan_roll"
}

func (roll *Roll) GetId() *int64 {
  panic("GetId() not implemented. Config should not be editted.")
}

func (roll *Roll) Roll(rRand *rand.Rand) int {
    result := 0

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

    return result
}

func RollAll(rolls []*Roll, rRand *rand.Rand) int {
  result := 0
  for _, roll := range rolls {
    result += roll.Roll(rRand)
  }

  return result
}

func FetchAllRolls(client utilities.ClientInterface, rolls *[]*Roll, parentId int64, tableName string, valueName string) {
  rollTableName := new(Roll).TableName("")
  client.FetchMany(rolls, parentId, tableName, rollTableName, valueName, "", false)
}
