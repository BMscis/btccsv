package main

import (
	"fmt"
	"log"
	"os"

	"github.com/wallerprogramm/btc-neo4j/graph"
	"github.com/wallerprogramm/btc-neo4j/neo4jdb"

	"github.com/chainswatch/bclib/btc"
	"github.com/joho/godotenv"
)

func main() {
	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
			fmt.Println("No .env")
		}
	}

	neo4jdb.Add_constraints()

	var c int = 1
	// var height uint32 = 53000
	var height uint32 = 1E5
	//var height uint32 = 10
	if err := btc.LoadFile(0, height, graph.Build, c); err != nil {
		log.Fatal(err)
	}

	neo4jdb.Close()
}
