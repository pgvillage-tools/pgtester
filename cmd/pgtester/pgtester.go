package main

import (
	"github.com/pgvillage-tools/pgtester/internal/pgtester"
)

func main() {
	pgtester.Initialize()
	pgtester.Handle()
}
