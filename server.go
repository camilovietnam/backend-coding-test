package main

import (
    "encoding/json"
    "log"
    "net/http"
    "io/ioutil"
    "time"

    "github.com/gorilla/mux"
)

type User struct {
    Id int64
    Login string
}

type Repository struct {
    Id int64
    Name string
    Description string
    stargazersCount int64
    Owner User
}

func getRepositories(w http.ResponseWriter, r *http.Request) {
    var repositories []Repository
    var outData []Repository
    var username = mux.Vars(r)["username"]
    var url = "https://api.github.com/users/" + username + "/repos"

    httpClient := http.Client{
        Timeout: time.Second * 10,
    }

    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        log.Fatal(err)
    }

    res, getErr := httpClient.Do(req)

    if getErr != nil {
        log.Fatal(getErr)
    }

    if res.Body != nil {
        defer res.Body.Close()
    }

    body, readError := ioutil.ReadAll(res.Body)

    if readError != nil {
        log.Fatal(readError)
    }

    jsonErr := json.Unmarshal(body, &repositories)

    if jsonErr != nil {
        log.Fatal(jsonErr)
    }

    for _, repositoryData := range repositories {
         var trimmedData = Repository {
            Id: repositoryData.Id,
            Name: repositoryData.Name,
            Description: repositoryData.Description,
//             stargazersCount: 0,
            Owner: User{
                Id: repositoryData.Owner.Id,
                Login: repositoryData.Owner.Login,
            },
//             ownerLogin: "logimn",
         }

         outData = append(outData, trimmedData)
    }

    jsonResponse, _ := json.Marshal(outData)
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}

func main() {
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/{username}/repositories", getRepositories).Methods("GET")

    log.Println("Server starts in port 5000")
    log.Fatal(http.ListenAndServe(":5000", router))
}