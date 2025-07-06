package scheduler

import (
	"log"
	"net/http"
	"github.com/robfig/cron/v3"
	"time"
)

//this schedular job run by using cron
func StartCronJob() {
	c := cron.New()

	c.AddFunc("@every 1m", func() {
		resp, err := http.Get("http://localhost:5001/api/students")
		if err != nil {
			log.Printf("Failed to call API: %v\n", err)
			return
		}
		defer resp.Body.Close()
		log.Println("Called sechudlar gob call : /api/students. Status:", resp.Status)
	})

	c.Start()
}

//this schedular job by using gorotines
func StartStudentFetchJob() {
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		for range ticker.C {
			resp, err := http.Get("http://localhost:5001/api/students")
			if err != nil {
				log.Printf("Error fetching students: %v\n", err)
				continue
			}
			log.Printf("Student fetch job ran. Status: %s\n", resp.Status)
			resp.Body.Close()
		}
	}()
}

