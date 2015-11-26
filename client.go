package main

import (
    "fmt"
    "log"
    "net/http"
    "strconv"
    "math/big"
    "crypto/md5"
    "encoding/hex"
    "io/ioutil"
    "github.com/julienschmidt/httprouter"
)

func hash(key_id string) int64 {
    hashFunc := md5.New()
    hashVal := big.NewInt(0)
    hashFunc.Write([]byte(key_id))
    hexstr := hex.EncodeToString(hashFunc.Sum(nil))
    hashVal.SetString(hexstr, 16)
    portNum := hashVal.Int64()
    if portNum < 0 {
        portNum *= -1
    }
    return portNum % 3 + 3000
}

func add(rw http.ResponseWriter, _ *http.Request, param httprouter.Params) {
    key := param.ByName("key_id")
    value := param.ByName("value")
    portNum := hash(key)

    url := "http://localhost:" + strconv.FormatInt(portNum, 10) + "/keys/" + key + "/" + value
    req, err := http.NewRequest("PUT", url, nil)
    if err != nil {
        log.Fatal(err)
    }

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        log.Fatal(err)
    }

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Fatal(err)
    }
	defer res.Body.Close()

    if err != nil {
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "Key needs to be integer\n")
    } else {
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(200)
        fmt.Fprintf(rw, string(body) + " to server: localhost:%d", portNum)
    }
}

func get(rw http.ResponseWriter, _ *http.Request, param httprouter.Params) {
    key := param.ByName("key_id")
    portNum := hash(key)

    url := "http://localhost:" + strconv.FormatInt(portNum, 10) + "/keys/" + key
    res, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Fatal(err)
    }
	defer res.Body.Close()

    if err != nil {
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "Key needs to be integer\n")
    } else {
        rw.Header().Set("Content-Type", "application/json")
        rw.WriteHeader(200)
        fmt.Fprintf(rw, string(body))
    }
}

func test() {
    values := []string {"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
    hashFunc := md5.New()

    for idx, value := range values {
		hashVal := big.NewInt(0)
		hashFunc.Write([]byte(strconv.Itoa(idx)))
		hexstr := hex.EncodeToString(hashFunc.Sum(nil))
		hashVal.SetString(hexstr, 16)
		portNum := hashVal.Int64()
		if portNum < 0 {
			portNum *= -1
		}
        portNum = portNum % 3 + 3000

        url := "http://localhost:" + strconv.FormatInt(portNum, 10) + "/keys/" + strconv.Itoa(idx+1) + "/" + value
        req, err := http.NewRequest("PUT", url, nil)
        if err != nil {
            log.Fatal(err)
        }

        _, err = http.DefaultClient.Do(req)
        if err != nil {
            log.Fatal(err)
        }
	}

    log.Println("End of program")
}

func main() {
    router := httprouter.New()
    router.PUT("/keys/:key_id/:value", add)
    router.GET("/keys/:key_id", get)

    log.Println("Server listening on 8080")
    http.ListenAndServe(":8080", router)
}
