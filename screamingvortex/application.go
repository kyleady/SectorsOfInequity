package main

import (
    "github.com/kyleady/SectorsOfInequity/screamingvortex/config"
    "github.com/kyleady/SectorsOfInequity/screamingvortex/grid"
    "github.com/kyleady/SectorsOfInequity/screamingvortex/utilities"
    "github.com/kyleady/SectorsOfInequity/screamingvortex/messages"

    "encoding/json"
    "log"
    "net/http"
)

func gridHandler(writer http.ResponseWriter, req *http.Request) {
    decoder := json.NewDecoder(req.Body)
    calixisMsg := new(messages.FromCalixis)
    inputErr := decoder.Decode(calixisMsg)
    if inputErr != nil {
        http.Error(writer, "Input\n" + inputErr.Error(), http.StatusInternalServerError)
        return
    }
    id := calixisMsg.ConfigId

    client := &utilities.Client{
      Environment: "dev",
      Region: "us-west-1",
      Resource: "koronus",
      Secret: "root",
    }
    client.Open()
    defer client.Close()

    gridConfig := config.LoadFrom(client, id)

    sectorConfig := new(grid.Sector)
    sectorConfig.Randomize(gridConfig)
    sectorConfig.SaveTo(client)

    writer.Write([]byte("OK"))
}

func main() {
    http.HandleFunc("/grid", gridHandler)
    log.Fatal(http.ListenAndServe(":5000", nil))
}
