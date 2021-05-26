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
    Stargazers_count int64
}

type LocalRepository struct {
    Id int64
    Name string
    Description string
    StargazersCount int64
    OwnerLogin string
    OwnerId int64
    Stargazers_count int64
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
    var repositories []Repository
    var outData []LocalRepository
    var username = mux.Vars(r)["username"]


    repositories = getGithubRepositories(username);

    for _, repositoryData := range repositories {
         var trimmedData = LocalRepository {
            Id: repositoryData.Id,
            Name: repositoryData.Name,
            Description: repositoryData.Description,
            StargazersCount: 0,
            OwnerLogin: repositoryData.Owner.Login,
            OwnerId: repositoryData.Owner.Id,
            Stargazers_count: repositoryData.Stargazers_count,
         }

         outData = append(outData, trimmedData)
    }

    jsonResponse, _ := json.MarshalIndent(outData, "", "");
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}

func main() {
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/{username}/repositories", getRepositories).Methods("GET")

    log.Println("Server starts in port 5000")
    log.Fatal(http.ListenAndServe(":5000", router))
}