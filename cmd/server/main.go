package main

import (
	"os"

	"github.com/rodrwan/bank/pkg/services/graph"
)

func main() {
	port := os.Getenv("PORT")
	graph.NewServer(port)
}
