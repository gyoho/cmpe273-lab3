package main

import (
    "fmt"
    "net/http"
    "strconv"
    "bytes"
    "github.com/julienschmidt/httprouter"
    "log"
)

var hashmap map[int]string

func add(rw http.ResponseWriter, _ *http.Request, param httprouter.Params) {
    key, err := strconv.Atoi(param.ByName("key_id"))
    value := param.ByName("value")

    if err != nil {
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "Key needs to be integer\n")
    } else {
        hashmap[key] = value
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(200)
        fmt.Fprintf(rw, "Added key/value.")
    }
}

func get(rw http.ResponseWriter, _ *http.Request, param httprouter.Params) {
    key, err := strconv.Atoi(param.ByName("key_id"))

    if err != nil {
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "Key needs to be integer\n")
    } else {
        value, ok := hashmap[key]
        if ok {
            rw.Header().Set("Content-Type", "application/json")
            rw.WriteHeader(200)

            jsonStr := `{
                "key" : "` + strconv.Itoa(key) + `",
                "value" : "` + value + `"
            }`

            fmt.Fprintf(rw, "%s\n", jsonStr)
        } else {
            rw.Header().Set("Content-Type", "plain/text")
            rw.WriteHeader(400)
            fmt.Fprintf(rw, "Not valid key\n")
        }
    }
}

func list(rw http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
    if len(hashmap) == 0 {
        rw.Header().Set("Content-Type", "application/json")
        rw.WriteHeader(200)
        fmt.Fprintf(rw, "{}")
    } else {
        var buffer bytes.Buffer
        buffer.WriteString("{\n[\n")
        for key, value := range hashmap {
            str := `{
                        "key" : "` + strconv.Itoa(key) + `",
                        "value" : "` + value + `"
                    },` + "\n"
            buffer.WriteString(str)
        }

        jsonStr := buffer.String()
        jsonStr = jsonStr[:len(jsonStr)-2]
        jsonStr = jsonStr + "\n]\n}"

        rw.Header().Set("Content-Type", "application/json")
        rw.WriteHeader(200)
        fmt.Fprintf(rw, "%s\n", jsonStr)
    }
}

func main() {
    hashmap = make(map[int]string)

    router := httprouter.New()
    router.PUT("/keys/:key_id/:value", add)
    router.GET("/keys/:key_id", get)
    router.GET("/keys", list)

    log.Println("Server listening on 3001")
    http.ListenAndServe(":3001", router)
}
