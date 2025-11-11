package main

import (
	"log"
	"time"
)

func main() {
	if RunLoop {
		for {
			log.Println("Starting trading step...")

			startTime := time.Now()

			if err := RunTradingStep(Symbol); err != nil {
				log.Printf("Error: %v\n", err)
			}

			log.Printf("Cost %.0f seconds\n", time.Now().Sub(startTime).Seconds())

			log.Printf("Sleeping %.0f minutes...\n", TimeSlice.Minutes())

			time.Sleep(TimeSlice)
		}
	} else {
		log.Println("Starting trading step...")

		startTime := time.Now()

		if err := RunTradingStep(Symbol); err != nil {
			log.Printf("Error: %v\n", err)
		}

		log.Printf("Cost %.0f seconds\n", time.Now().Sub(startTime).Seconds())
	}
}
