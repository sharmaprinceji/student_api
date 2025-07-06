package scheduler

import (
	"log"
	"net/http"
	"github.com/robfig/cron/v3"
)

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
