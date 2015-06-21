package main

import (
    "fmt"
    "os"
)

const FILENAME="images-cve.dat"

func writeToFile( imageName string, cveName string) {
    f, err := os.OpenFile(FILENAME, os.O_APPEND|os.O_WRONLY, 0600)
    if err !=nil {
        fmt.Printf("failed opening %s", FILENAME)
        return
    }
    defer f.Close()
    imageAndCve := fmt.Sprintf("%s@%s\r\n", cveName, imageName)
    f.WriteString(imageAndCve)
}

func main() {
    writeToFile("image123", "cve456")
}
