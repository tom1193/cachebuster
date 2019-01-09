//Author: tom1193. File busting API for files in:
//https://s3.console.aws.amazon.com/s3/buckets/ph-mode-static/datavis-library/?region=us-east-1&tab=overview
//based on following resources:
//https://medium.com/@adigunhammedolalekan/build-and-deploy-a-secure-rest-api-with-go-postgresql-jwt-and-gorm-6fadf3da505b
//https://www.alexedwards.net/blog/golang-response-snippets#string

package main

import (
    "fmt"
    "log"
    "os"
    "encoding/json"
    "github.com/gorilla/mux"
    "net/http"
    "github.com/rs/cors"
    "github.com/tom1193/cachebuster/filecache"
    "github.com/tom1193/cachebuster/utils"
    "github.com/tom1193/cachebuster/proto"
)

func PostFiles(w http.ResponseWriter, r *http.Request) {
    //POST should be in JSON
    decoder := json.NewDecoder(r.Body)
    var pr proto.PostRequest
    err := decoder.Decode(&pr)
    if err != nil {
        m := utils.Message(false, "Invalid request, cannot parse json body")
        w.WriteHeader(http.StatusBadRequest)
        utils.Respond(w, m)
    }
    //update file cache
    res, status := filecache.UpdateFileCache(pr)
    w.WriteHeader(status)
    if status != http.StatusCreated {
        utils.Respond(w, res)
    }
}

//GET requests should respond in browser-compatible JSON
func RequestFiles(w http.ResponseWriter, r *http.Request) {
    params := r.URL.Query() //returns Values type which is a map[string][]string of the parameters
    res, status := filecache.RequestFileCache(params["filenames"], params["env"][0])
    w.WriteHeader(status)
    utils.Respond(w, res)
}

func EchoFiles(w http.ResponseWriter, r *http.Request) {
    params := r.URL.Query()
    res, status := filecache.EchoFileCache(params["env"][0])
    w.WriteHeader(status)
    utils.Respond(w, res)
}

func main () {
    router := mux.NewRouter()
    //cors allows browsers to call API
    handler := cors.Default().Handler(router)
    //handle api calls
    router.HandleFunc("/post", PostFiles).Methods("POST")
    router.HandleFunc("/get", RequestFiles).Methods("GET")
    router.HandleFunc("/echo", EchoFiles).Methods("GET")
    port := os.Getenv("PORT")
    if port == "" {
        port = "8000" //localhost
    }
    fmt.Println(port)
    log.Fatal(http.ListenAndServe(":" + port, handler))
}
