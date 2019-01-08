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
    // "encoding/json"
    "github.com/gorilla/mux"
    "net/http"
    "github.com/rs/cors"
    "github.com/tom1193/cachebuster/filecache"
    "github.com/tom1193/cachebuster/utils"
)

func PostFiles(w http.ResponseWriter, r *http.Request) {
    params := r.URL.Query() //returns Values type which is a map[string][]string of the parameters
    fmt.Println(params)
    res, err := filecache.UpdateFileCache(params["filenames"], params["env"][0])
    if err != nil {
        http.Error(w, err.Error(), res)
        return
    }
    w.WriteHeader(res)
}

//GET requests should respond in browser-compatible JSON
func RequestFiles(w http.ResponseWriter, r *http.Request) {
    params := r.URL.Query()
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
