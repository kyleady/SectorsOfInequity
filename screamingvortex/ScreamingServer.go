package main

import (
    "github.com/kyleady/SectorsOfInequity/screamingvortex/config"
    "github.com/kyleady/SectorsOfInequity/screamingvortex/grid"

    "encoding/json"
    "log"
    "net/http"
)

func gridHandler(writer http.ResponseWriter, req *http.Request) {
    log.Print(req)
    decoder := json.NewDecoder(req.Body)
    gridConfig := new(config.GridConfig)
    gridConfig.SetToDefault()
    inputErr := decoder.Decode(gridConfig)
    if inputErr != nil {
        http.Error(writer, "Input\n" + inputErr.Error(), http.StatusInternalServerError)
        return
    }

    sector := new(grid.Sector)
    sector.Randomize(gridConfig)
    jsSector, outputErr := json.Marshal(&sector)
    if outputErr != nil {
      http.Error(writer, "Output\n" + outputErr.Error(), http.StatusInternalServerError)
      return
    }

    writer.Header().Set("Content-Type", "application/json")
    writer.Write(jsSector)
}

func main() {
    http.HandleFunc("/grid/", gridHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))


}
