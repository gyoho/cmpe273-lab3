package main

import (
    "log"
    "net/http"
    "strconv"
    "math/big"
    "crypto/md5"
    "encoding/hex"
)

func main() {
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
