package main

import (
    "screamingvortex/asset"
    "screamingvortex/config"
    "screamingvortex/utilities"
    "screamingvortex/messages"

    "math/rand"
    "time"
    "encoding/json"
    "log"
    "net/http"
)

func extractIds(writer http.ResponseWriter, req *http.Request) (int64, int64, int64) {
  decoder := json.NewDecoder(req.Body)
  calixisMsg := new(messages.FromCalixis)
  inputErr := decoder.Decode(calixisMsg)
  if inputErr != nil {
      http.Error(writer, "Input\n" + inputErr.Error(), http.StatusInternalServerError)
      return -1, -1, -1
  }

  return calixisMsg.PerterbationId, calixisMsg.TypeId, calixisMsg.JobId
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

func assetHandler(writer http.ResponseWriter, req *http.Request) {
    perterbationId, typeId, jobId := extractIds(writer, req)
    if perterbationId == -1 && typeId == -1 && jobId == -1 { return }

    go func() {
      client := createClient()
      client.Open()
      defer client.Close()

      t := time.Now()
      randSource := rand.NewSource(t.UnixNano())
      rRand := rand.New(randSource)

      basePerterbation := config.CreateEmptyPerterbation(client, rRand)
      perterbation := basePerterbation.AddPerterbation(perterbationId)

      job := utilities.CreateJob(jobId, 1)
      rolledAsset := asset.RollAsset(perterbation, typeId, "", 1)
      job.Step(client)
      rolledAsset.SaveTo(client)
      job.Complete(client, rolledAsset.Id)
    }()

    writer.Write([]byte("OK"))
}

func main() {
    http.HandleFunc("/asset", assetHandler)
    log.Fatal(http.ListenAndServe(":5000", nil))
}
