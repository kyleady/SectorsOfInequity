package main

import (
    "screamingvortex/asset"
    "screamingvortex/config"
    "screamingvortex/grid"
    "screamingvortex/utilities"
    "screamingvortex/messages"

    "encoding/json"
    "log"
    "net/http"
)

func extractIds(writer http.ResponseWriter, req *http.Request) (int64, int64) {
  decoder := json.NewDecoder(req.Body)
  calixisMsg := new(messages.FromCalixis)
  inputErr := decoder.Decode(calixisMsg)
  if inputErr != nil {
      http.Error(writer, "Input\n" + inputErr.Error(), http.StatusInternalServerError)
      return -1, -1
  }

  return calixisMsg.ConfigId, calixisMsg.JobId
}

func createClient() *utilities.Client {
  return &utilities.Client{
    Local: "client_secrets.json",
    Environment: "dev",
    Region: "us-west-1",
    Resource: "koronus",
    Secret: "root",
  }
}

func gridHandler(writer http.ResponseWriter, req *http.Request) {
    configId, jobId := extractIds(writer, req)
    if configId < 0 { return }

    go func() {
      client := createClient()
      client.Open()
      defer client.Close()

      manager := config.CreateEmptyManager(client)
      gridConfig := config.FetchGrid(manager, configId)

      job := utilities.CreateJob(jobId, 3)
      sectorGrid := new(grid.Sector)
      sectorGrid.Randomize(gridConfig, client, job)
      sectorGrid.SaveTo(client)
      job.Complete(client, sectorGrid.Id)
    }()

    writer.Write([]byte("OK"))
}

func sectorHandler(writer http.ResponseWriter, req *http.Request) {
    configId, jobId := extractIds(writer, req)
    if configId < 0 { return }

    go func() {
      client := createClient()
      client.Open()
      defer client.Close()

      sectorGrid := grid.LoadSectorFrom(client, configId)

      job := utilities.CreateJob(jobId, len(sectorGrid.Systems))
      sectorAsset := asset.RandomSector(sectorGrid, client, job)
      sectorAsset.SaveTo(client)
      job.Complete(client, sectorAsset.Id)
    }()

    writer.Write([]byte("OK"))
}

func main() {
    http.HandleFunc("/grid", gridHandler)
    http.HandleFunc("/sector", sectorHandler)
    log.Fatal(http.ListenAndServe(":5000", nil))
}
