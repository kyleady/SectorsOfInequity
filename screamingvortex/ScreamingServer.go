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

    client := &utilities.Client{
      Environment: "dev",
      Region: "us-west-1",
      Resource: "koronus",
      Secret: "root",
    }
    client.Open()
    defer client.Close()

    gridConfig := new(config.GridConfig)
    client.Fetch(gridConfig, calixisMsg.ConfigId)

    sector := new(grid.Sector)
    sector.Randomize(gridConfig)

    jsSector, outputErr := json.Marshal(*sector)
    if outputErr != nil {
      http.Error(writer, "Output\n" + outputErr.Error(), http.StatusInternalServerError)
      return
    }

    writer.Header().Set("Content-Type", "application/json")
    writer.Write(jsSector)
}

func main() {
    http.HandleFunc("/grid", gridHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
