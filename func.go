package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
)

const (
    FILENAME = "images-cve.dat"
    FORMAT   = "%s@%s\r\n"
)

func writeToFile( imageName string, cveName string) {
    f, err := os.OpenFile(FILENAME, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
    if err !=nil {
        fmt.Printf("failed opening %s:%s", FILENAME, err)
        return
    }
    defer f.Close()
    imageAndCve := fmt.Sprintf(FORMAT, cveName, imageName)
    f.WriteString(imageAndCve)
}


func existInFile(imageName string, cveName string) bool {
    imageAndCve := fmt.Sprintf(FORMAT, cveName, imageName)
    
    f, err := os.OpenFile(FILENAME, os.O_RDONLY, 0)
    if err !=nil {
        fmt.Printf("failed opening %s", FILENAME)
        return false
    }
    defer f.Close()
    
    scanner := bufio.NewScanner(f)   
    
    found := false
    for scanner.Scan() {
        line := scanner.Text()
        if (strings.TrimSpace(line) == strings.TrimSpace(imageAndCve)) {
            fmt.Printf("found!")
            found = true
            break
        }
    }
    return found
}


func main() {
    writeToFile("image123", "cve456")
    exist := existInFile("amir", "123")
    if (exist == true) {
        fmt.Printf("bug!")
    } else {
        fmt.Printf("ok")
    }
    exist = existInFile("image123", "cve456")
    if (exist == false) {
        fmt.Printf("bug!")
    } else {
        fmt.Printf("ok")
    }

}

