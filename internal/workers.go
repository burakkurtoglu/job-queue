package internal

import (
	"fmt"
	db "job-queue/database"
	"log"
	"strconv"
	"time"
)

var jobChan = make(chan db.Data, 10)

func dispatcher() {

	for {
		datas, err := db.GetStatusDb()
		if err != nil {
			log.Println("Dispatcher got DB error, trying again next. Error:", err)
			time.Sleep(time.Millisecond * 500)
			continue
		}
		for _, job := range datas {
			if job.STATUS == 0 {
				jobChan <- job
			}
		}
		time.Sleep(time.Second * 2)
	}

}

func worker() {

	for {

		job := <-jobChan
		jobLong, err := strconv.Atoi(job.PAYLOAD)
		if err != nil {
			log.Println("worker couldn't convert payload into int. Error:", err)
			continue
		}

		jobLong /= (job.RETRY_COUNT + 1)

		// 0 : waiting, 1: pending, 2: done, 3: failed
		job.STATUS = 1
		if err := db.UpdateJobStatusDb(job.STATUS, job.RETRY_COUNT, job.ID); err != nil {
			log.Println("Worker got db status update error:", err)
			continue
		}
		if jobLong > 50 && job.RETRY_COUNT > 3 { // FAILED
			time.Sleep(time.Microsecond * 10)
			job.STATUS = 3
			job.RETRY_COUNT += 1
			if err := db.UpdateJobStatusDb(job.STATUS, job.RETRY_COUNT, job.ID); err != nil {
				log.Println("Worker got db status update error:", err)
				continue
			}
			fmt.Printf("Failed to process job, id: %d\n", job.ID)
			continue

		} else if jobLong > 50 && job.RETRY_COUNT < 3 {
			time.Sleep(time.Microsecond * 10)
			job.STATUS = 0
			job.RETRY_COUNT += 1
			if err := db.UpdateJobStatusDb(job.STATUS, job.RETRY_COUNT, job.ID); err != nil {
				log.Println("Worker got db status update error:", err)
				continue
			}
			fmt.Printf("Failed to process job, trying again, id: %d\n", job.ID)

		} else if jobLong <= 50 && job.RETRY_COUNT <= 3 { // done
			time.Sleep(time.Duration(jobLong) * time.Microsecond)
			job.STATUS = 2
			if err := db.UpdateJobStatusDb(job.STATUS, job.RETRY_COUNT, job.ID); err != nil {
				log.Println("Worker got db status update error:", err)
				continue
			}
			fmt.Printf("The job done successfuly, id: %d\n", job.ID)
			continue
		}
	}
}

func StartWorkerPool() {
	go dispatcher()
	go worker()
	go worker()
	go worker()
	go worker()
	go worker()

}
