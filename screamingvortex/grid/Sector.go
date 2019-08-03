package grid

import (
  "github.com/kyleady/SectorsOfInequity/screamingvortex/config"
  "math/rand"
  "fmt"
  "strconv"
)

type Sector struct {
  Grid [][]System
  config *config.GridConfig
}

func (sector *Sector) Randomize(gridConfig *config.GridConfig) {
  sector.config = gridConfig

  LogSector(sector)
  sector.createGrid()
  LogSector(sector)
  sector.populateGrid()
  LogSector(sector)
  sector.connectSystems()
  LogSector(sector)
  sector.trimToLargestBlob()
  LogSector(sector)
}

func (sector *Sector) createGrid() {
  sector.Grid = make([][]System, sector.config.Height)
  for i := range sector.Grid {
    sector.Grid[i] = make([]System, sector.config.Width)
  }
}

func (sector *Sector) populateGrid() {
  for i, row := range sector.Grid {
    for j := range row {
      sector.createSystem(i, j)
    }
  }
}

func (sector *Sector) createSystem(i int, j int) {
  system := new(System)
  if roll := rand.Float64(); roll < sector.config.PopulationRate {
    system.InitializeAt(i, j)
  } else {
    system.SetToVoidSpace()
  }

  sector.Grid[i][j] = *system
}

func (sector *Sector) connectSystems() {
  reach := sector.config.ConnectionRange

  for i, row := range sector.Grid {
    for j, system := range row {
      for r_i := -reach; r_i <= reach; r_i++ {
        for r_j := -reach; r_j <= reach; r_j++ {
          if roll := rand.Float64(); roll < sector.connectionChance(r_i, r_j) {
            targetSystem := sector.GetSystem(i+r_i, j+r_j)
            system.ConnectTo(targetSystem)
          }
        }
      }
    }
  }
}

func (sector *Sector) connectionChance(r_i int, r_j int) float64 {
  abs_i := r_i
  if r_i < 0 {
    abs_i = - r_i
  }

  abs_j := r_j
  if r_j < 0 {
    abs_j = - r_j
  }

  reach := abs_i
  if abs_j > abs_i {
    reach = abs_j
  }

  chance := sector.config.ConnectionRate
  for i := 0; i < reach - 1; i++ {
    chance *= sector.config.RangeRateMultiplier
  }

  return chance
}

func (sector *Sector) GetSystem(i int, j int) *System {
  if i < 0 || i >= sector.config.Height {
    return nil
  }

  if j < 0 || j >= sector.config.Width {
    return nil
  }

  return &sector.Grid[i][j]
}

func (sector *Sector) labelBlobsAndGetSizes() []int {
  currentLabel := 0
  blobSizes := make([]int, 0)
  for i, row := range sector.Grid {
    fmt.Println("<==>")
    for j, system := range row {

      system = sector.Grid[i][j]
      blobSize := system.LabelBlob(currentLabel)
      fmt.Println(i, j, currentLabel, blobSize)
      if blobSize > 0 {
        blobSizes = append(blobSizes, blobSize)
        currentLabel++
      }
    }
  }

  return blobSizes
}

func (sector *Sector) trimToLargestBlob() {
  blobSizes := sector.labelBlobsAndGetSizes()
  largestBlobSize := 0
  largestBlobLabel := -1
  for blobLabel, blobSize := range blobSizes {
    if blobSize > largestBlobSize {
      largestBlobSize = blobSize
      largestBlobLabel = blobLabel
    }
  }

  for _, row := range sector.Grid {
    for _, system := range row {
      system.VoidNonMatchingLabel(largestBlobLabel)
    }
  }
}




func LogSector(sector *Sector) {
  systems := 0
  output := ""
  for _, row := range sector.Grid {
    for _, system := range row {
      //output += systemConnectionsToString(&system)
      output += " " + strconv.Itoa(system.Label())
      if !system.IsVoidSpace() {
        systems++
      }
    }
    output += "\n"
  }

  fmt.Println(output)
  fmt.Println("Systems: " + strconv.Itoa(systems))
}

func systemConnectionsToString(system *System) string {
  spaces := 3
  connection_output := ""
  if !system.IsVoidSpace() {
    connections := len(system.Routes)
    if connections < 10 {
      spaces -= 1
    } else if connections < 100 {
      spaces -= 2
    } else {
      spaces = 0
    }

    connection_output = strconv.Itoa(connections)
  }

  spacing := ""
  for i := 0; i < spaces; i++ {
      spacing += " "
  }

  return spacing + connection_output
}
