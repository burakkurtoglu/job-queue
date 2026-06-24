package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Data struct {
	ID          int
	NAME        string
	PAYLOAD     string
	STATUS      int
	RETRY_COUNT int
	CREATED_AT  time.Time
	UPDATED_AT  time.Time
}

var db *sql.DB
var err error

func InitDb() {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	passw := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	sslM := os.Getenv("DB_SSLMODE")

	connSTR := fmt.Sprintf("host=%s port=%s user=%s passw0rd=%s dbname=%s sslmode=%s", host, port, user, passw, name, sslM)
	db, err = sql.Open("postgres", connSTR)
	if err != nil {
		log.Fatal(err)
	}

	queryStr := `CREATE TABLE IF NOT EXISTS jobs(
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		payload TEXT,
		status INT DEFAULT 0, -- 0: pending, 1: processing, 2: completed. 3:failed
		retry_count INT DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(queryStr)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertDb(Name string, Payload string) {

	queryStr := `INSERT INTO jobs (name, payload)
	VALUES ($1, $2);`

	_, err = db.Exec(queryStr, Name, Payload)
	if err != nil {
		log.Fatal(err)
	}
}

func GetJobDb(id int) Data {

	queryStr := `SELECT * FROM jobs WHERE id=$1;`
	DbData := db.QueryRow(queryStr, id)
	var data Data
	err = DbData.Scan(&data.ID, &data.NAME, &data.PAYLOAD, &data.STATUS, &data.RETRY_COUNT, &data.CREATED_AT, &data.UPDATED_AT)
	if err != nil {
		//log.Fatal(err)
	}
	return data
}

func GetStatusDb() []Data {
	queryStr := `SELECT * FROM jobs WHERE status=0 ORDER BY created_at LIMIT 5;`
	Job, err := db.Query(queryStr)
	defer Job.Close()
	var datas []Data

	for Job.Next() {
		var data Data
		err = Job.Scan(&data.ID, &data.NAME, &data.PAYLOAD, &data.STATUS, &data.RETRY_COUNT, &data.CREATED_AT, &data.UPDATED_AT)
		if err != nil {

		}
		datas = append(datas, data)
	}
	return datas
}

func UpdateJobStatusDb(status int, retry int, id int) {
	queryStr := `UPDATE jobs SET status=$1, updated_at=CURRENT_TIMESTAMP, retry_count=$2 WHERE id=$3;`
	_, err = db.Exec(queryStr, status, retry, id)
	if err != nil {
	}

	//fmt.Printf("Status of id %d changed successfuly!\n", id)
}
