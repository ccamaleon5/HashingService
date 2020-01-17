package main

import (
    "fmt"
    "log"
    "bytes"
    "time"
    "os"
    "io/ioutil"
	"net/http"
    "encoding/json"
    "github.com/go-openapi/strfmt"
    lib "github.com/ccamaleon5/HashingService/lib"
    model "github.com/ccamaleon5/HashingService/model"
)

func main(){
	setupRoutes()
}

func validateHash(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    enableCors(&w)
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    // Parse our multipart form, 10 << 20 specifies a maximum
    // upload of 10 MB files.
    parseErr := r.ParseMultipartForm(10 << 20)
    if parseErr != nil{
        fmt.Println("error:",parseErr)
        http.Error(w, "failed to parse multipart message", http.StatusBadRequest)
        return
    }
    
    file, handler, err := r.FormFile("media")
    if err != nil {
        fmt.Println("Error Retrieving the File")
        fmt.Println(err)
        return
    }
    defer file.Close()
    fmt.Printf("Uploaded File: %+v\n", handler.Filename)
    fmt.Printf("File Size: %+v\n", handler.Size)
    fmt.Printf("MIME Header: %+v\n", handler.Header)

    hash := lib.Hash(file)

    response:="{\"hash\":\""+hash+"\"}"

    _, err = file.Seek(0, os.SEEK_SET)
    if err != nil {
        fmt.Println(err)
    }

    w.Write([]byte(response)) 
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    enableCors(&w)
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    // Parse our multipart form, 10 << 20 specifies a maximum
    // upload of 10 MB files.
    parseErr := r.ParseMultipartForm(10 << 20)
    if parseErr != nil{
        fmt.Println("error:",parseErr)
        http.Error(w, "failed to parse multipart message", http.StatusBadRequest)
        return
    }
    
    file, handler, err := r.FormFile("media")
    if err != nil {
        fmt.Println("Error Retrieving the File")
        fmt.Println(err)
        return
    }
    defer file.Close()
    fmt.Printf("Uploaded File: %+v\n", handler.Filename)
    fmt.Printf("File Size: %+v\n", handler.Size)
    fmt.Printf("MIME Header: %+v\n", handler.Header)

    metadata, errMeta := getMetadata(r)
	if errMeta != nil {
		http.Error(w, "failed to get metadata", http.StatusBadRequest)
		return
	}
	log.Println("Metadata:",string(metadata))

    res := model.Metadata{}
    json.Unmarshal(metadata, &res)
    
    res.Document = lib.Hash(file)

    _, err = file.Seek(0, os.SEEK_SET)
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println("metadata json:",res)

    responseCredential := createCredential(&res) 

    fmt.Println("responseCredential:",responseCredential)

    w.Write([]byte(responseCredential)) 
}

func getMetadata(r *http.Request) ([]byte, error) {
	f, _, err := r.FormFile("metadata")
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata form file: %v", err)
	}

	metadata, errRead := ioutil.ReadAll(f)
	if errRead != nil {
		return nil, fmt.Errorf("failed to read metadata: %v", errRead)
	}

	return metadata, nil
}

func enableCors(w *http.ResponseWriter) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
    (*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Access-Control-Allow-Origin, Authorization, X-Requested-With")
}

func setupRoutes() {
	fmt.Println("Init Hashing Server")
    http.HandleFunc("/upload", uploadFile)
    http.HandleFunc("/validate", validateHash)
    http.ListenAndServe(":9000", nil)
}

func createCredential(metadata *model.Metadata)(string){
    credentials := make([]*model.CredentialSubject, 0, 50)
    credentialSubject := model.CredentialSubject{}
    credentialSubject.Type = "DocumentHashCredential"
    
    fmt.Println("ffff:",time.Now().UTC().Format("2006-01-02T15:04:05Z"))
    issuanceDate, err := strfmt.ParseDateTime(time.Now().UTC().Format("2006-01-02T15:04:05Z"))
    if err != nil{
        fmt.Println("Error:",err)
    }

    expirationTime, err := time.Parse("2006-01-02T15:04:05Z", metadata.ExpirationDate)
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println("expirationTime:",expirationTime.UTC().Format("2006-01-02T15:04:05Z"))

    expirationDate, err := strfmt.ParseDateTime(expirationTime.UTC().Format("2006-01-02T15:04:05Z"))
    if err != nil{
        fmt.Println("Error:",err)
    }

    credentialSubject.IssuanceDate = issuanceDate
    credentialSubject.ExpirationDate = expirationDate

    fmt.Println("credentialSubject:",credentialSubject.IssuanceDate)    

    credentialSubject.Content = &metadata
    credentials = append(credentials, &credentialSubject)
    jsonValue, _ := json.Marshal(credentials)
    fmt.Println("#####REQUEST####", string(jsonValue))
    
    timeout := time.Duration(10 * time.Second)
    client := http.Client{
        Timeout: timeout,
    }

    req, err := http.NewRequest("POST", "http://localhost:8000/v1/credential",  bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-type","application/json")
    req.Header.Set("accept","application/json")

    response, err := client.Do(req)

    if err != nil {
        fmt.Printf("The HTTP request failed with error %s\n", err)
    } else {
        data, _ := ioutil.ReadAll(response.Body)
        fmt.Println(string(data))
        return string(data)
    }
    return "{}"
}