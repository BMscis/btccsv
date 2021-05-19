package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BMscis/btccsv/grapho"

	"github.com/chainswatch/bclib/btc"
	"github.com/joho/godotenv"
	logs "github.com/sirupsen/logrus"
)

func main() {
	start := time.Now()

	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
			fmt.Println("No .env")
		}
	}

	var c int = 1

	//var height uint32 = 1E5
	var height uint32 = 564700

	if err := btc.LoadFile(0, height, grapho.Build, c); err != nil {
		log.Fatal(err)
	}

	end := time.Now()
	fmt.Println("------------------------------------------------")
	logs.Infof("Started at %v\n", start)
	logs.Infof("Ended at %v\n", end)
	logs.Infof("Parsing until height=%v took: %v\n", height, end.Sub(start))

	grapho.W.Close()
	grapho.T.Close()
}
