# job-queue

A background job queue service built with Go and PostgreSQL. HTTP endpoints accept jobs, a dispatcher polls the database and feeds a worker pool via channels, and workers process jobs concurrently with automatic retry logic.

**Live:** https://job-queue-ehrb.onrender.com

## Features

- `POST /jobs` — enqueue a new job
- `GET /jobs/{id}` — query job status
- Worker pool with goroutines processing jobs concurrently
- Automatic retry (up to 3 attempts) before marking a job as failed
- PostgreSQL persistence — jobs survive restarts
- Environment-based configuration, no hardcoded secrets

## How It Works

```
HTTP Request → Handler → PostgreSQL
                              ↑
                         Dispatcher (polls every 2s)
                              ↓
                           Channel
                              ↓
                      Worker Pool (5 workers)
                              ↓
                    Update status in PostgreSQL
```

Job statuses: `0` waiting → `1` pending → `2` done / `3` failed (retried up to 3x)

## Project Structure

```
job-queue/
├── server/         # main.go — HTTP server entry point
├── internal/       # handlers, dispatcher, worker pool
├── database/       # DB connection and queries
├── go.mod
└── go.sum
```

## Running Locally

```bash
git clone https://github.com/burakkurtoglu/job-queue
cd job-queue
```

Create a `.env` file:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=jobDb
DB_SSLMODE=disable
```

Create the jobs table in PostgreSQL:

```sql
CREATE TABLE IF NOT EXISTS jobs (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    payload TEXT,
    status INT DEFAULT 0,
    retry_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

Run:

```bash
go run ./server/
```

## API

**Enqueue a job**
```
POST /jobs
Content-Type: application/json

{"name": "example", "payload": "75"}
```

**Check job status**
```
GET /jobs/{id}
```

## Stack

- Go 1.26
- PostgreSQL
- `net/http` — no external HTTP framework
- `database/sql` + `lib/pq`
- Deployed on Render
