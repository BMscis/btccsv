package timelog

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type TimeWriter struct {
	measures *os.File

	wmeasures *csv.Writer

	start time.Time
}

func NewWriter() *TimeWriter {
	var err error
	t := new(TimeWriter)
	// Tools for time measurement
	t.measures, err = os.Create("measure.csv")
	if err != nil {
		log.Fatalln("Error creating measure.csv", err)
	}
	t.wmeasures = csv.NewWriter(t.measures)

	t.start = time.Now()

	return t
}

func (t *TimeWriter) Close() {
	t.wmeasures.Flush()

	t.measures.Close()
}

func (t *TimeWriter) Log_time(height uint32) {
	if height%100 == 0 {
		err := t.wmeasures.Write([]string{fmt.Sprint(height), time.Now().Sub(t.start).String()})
		if err != nil {
			log.Fatalln("Error writing time at height: "+fmt.Sprint(height), err)
		}
	}
}
