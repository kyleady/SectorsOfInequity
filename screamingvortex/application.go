package main

import (
    "github.com/kyleady/SectorsOfInequity/screamingvortex/asset"
    "github.com/kyleady/SectorsOfInequity/screamingvortex/config"
    "github.com/kyleady/SectorsOfInequity/screamingvortex/grid"
    "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"
    "github.com/kyleady/SectorsOfInequity/screamingvortex/messages"

    "encoding/json"
    "log"
    "net/http"
)

func extractConfigId(writer http.ResponseWriter, req *http.Request) int64 {
  decoder := json.NewDecoder(req.Body)
  calixisMsg := new(messages.FromCalixis)
  inputErr := decoder.Decode(calixisMsg)
  if inputErr != nil {
      http.Error(writer, "Input\n" + inputErr.Error(), http.StatusInternalServerError)
      return -1
  }

  return calixisMsg.ConfigId
}

func createClient() *utilities.Client {
  return &utilities.Client{
    Environment: "dev",
    Region: "us-west-1",
    Resource: "koronus",
    Secret: "root",
  }
}

func gridHandler(writer http.ResponseWriter, req *http.Request) {
    id := extractConfigId(writer, req)
    if id < 0 { return }
    client := createClient()
    client.Open()
    defer client.Close()

    gridConfig := config.LoadGridFrom(client, id)

    sectorGrid := new(grid.Sector)
    sectorGrid.Randomize(gridConfig)
    sectorGrid.SaveTo(client)

    writer.Write([]byte("OK"))
}

func sectorHandler(writer http.ResponseWriter, req *http.Request) {
    id := extractConfigId(writer, req)
    if id < 0 { return }
    client := createClient()
    client.Open()
    defer client.Close()

    sectorGrid := grid.LoadSectorFrom(client, id)

    sectorAsset := asset.RandomSector(sectorGrid, client)
    sectorAsset.SaveTo(client)

    writer.Write([]byte("OK"))
}

func main() {
    http.HandleFunc("/grid", gridHandler)
    http.HandleFunc("/sector", sectorHandler)
    log.Fatal(http.ListenAndServe(":5000", nil))
}
