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
    Stargazers_count int64
    Owner User
}

type LocalRepository struct {
    Id int64
    Name string
    Description string
    Stargazers_count int64
    OwnerId int64
    OwnerLogin string
}

func main() {
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/{username}/repositories", getRepositories).Methods("GET")

    log.Println("Server started in port 5000")
    log.Fatal(http.ListenAndServe(":5000", router))
}

func getGithubRepositories(repositoryName string) []Repository{
    var url string = "https://api.github.com/users/" + repositoryName + "/repos"
    var repositories []Repository

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

    return repositories
}

func getRepositories(w http.ResponseWriter, r *http.Request) {
    var username = mux.Vars(r)["username"]
    var repositories []Repository = getGithubRepositories(username)
    var outputData []LocalRepository

    for _, repositoryData := range repositories {
         var data = LocalRepository {
            Id: repositoryData.Id,
            Name: repositoryData.Name,
            Description: repositoryData.Description,
            OwnerLogin: repositoryData.Owner.Login,
            OwnerId: repositoryData.Owner.Id,
            Stargazers_count: repositoryData.Stargazers_count,
         }

         outputData = append(outputData, data)
    }

    jsonResponse, _ := json.MarshalIndent(outputData, "", "");
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}

