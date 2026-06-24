package internal

import (
	"fmt"
	db "job-queue/database"
	"strconv"
	"time"
)

var jobChan = make(chan db.Data, 10)

func dispatcher() {

	for {
		datas := db.GetStatusDb()
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
			continue
		}

		jobLong /= (job.RETRY_COUNT + 1)

		// 0 : waiting, 1: pending, 2: done, 3: failed
		job.STATUS = 1
		db.UpdateJobStatusDb(job.STATUS, job.RETRY_COUNT, job.ID)
		if jobLong > 50 && job.RETRY_COUNT > 3 { // FAILED
			time.Sleep(time.Microsecond * 10)
			job.STATUS = 3
			job.RETRY_COUNT += 1
			db.UpdateJobStatusDb(job.STATUS, job.RETRY_COUNT, job.ID)
			fmt.Printf("Failed to process job, id: %d\n", job.ID)
			continue

		} else if jobLong > 50 && job.RETRY_COUNT < 3 {
			time.Sleep(time.Microsecond * 10)
			job.STATUS = 0
			job.RETRY_COUNT += 1
			db.UpdateJobStatusDb(job.STATUS, job.RETRY_COUNT, job.ID)
			fmt.Printf("Failed to process job, trying again, id: %d\n", job.ID)

		} else if jobLong <= 50 && job.RETRY_COUNT <= 3 { // done
			time.Sleep(time.Duration(jobLong) * time.Microsecond)
			job.STATUS = 2
			db.UpdateJobStatusDb(job.STATUS, job.RETRY_COUNT, job.ID)
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
