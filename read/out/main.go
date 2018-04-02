package main

import (
    "os"
)

func main() {
    os.Stderr.WriteString("Not supported.\n")
    os.Exit(1)
}
