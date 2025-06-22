

*Job Scheduling System*

*Prerequisites*

1. Install the latest stable version of Golang and Docker.
2. Verify that the Docker engine is running.
3. Ensure cURL is installed.
4. For detailed architecture, refer to docs/backend_architecture.md.

*Steps to Execute the Code*

*Step 1: Clone the Repository*
```
bash
git clone git@github.com:vijayshahwal/jobScheduling.git
```

*Step 2: Navigate to the Project Directory*
```
bash
cd jobScheduling/
```

*Step 3: Initialize Go Modules*
```
bash
go mod init github.com/vijayshahwal/jobScheduling
```

*Step 4: Tidy Go Modules*
```
bash
go mod tidy
```

*Step 5: Run Redis Container*
```
bash
docker run -d --name redis -p 6379:6379 redis:latest
```

*Step 6: Run MongoDB Container*
```
bash
docker run -d --name mongo -p 27017:27017 mongo:latest
```

*Step 7: Run the Application*
```
bash
MONGO_URL="mongodb://localhost:27017" MONGO_DB_NAME="job_scheduler" REDIS_HOST="localhost:6379" REDIS_PASSWORD="" REDIS_DB="0" go run main.go
```

*Step 8: Create Jobs*
Open a new integrated terminal in the same folder and execute the following commands:

*Create Jobs*
```
bash
curl --location 'http://localhost:8080/job' \
--header 'Content-Type: application/json' \
--data '{ "id": "123", "name": "process 3", "description":"send whatsapp notification", "priority": 6 }'

curl --location 'http://localhost:8080/job' \
--header 'Content-Type: application/json' \
--data '{ "id": "124", "name": "process 4", "description":"send sms notification", "priority": 6 }'

curl --location 'http://localhost:8080/job' \
--header 'Content-Type: application/json' \
--data '{ "id": "125", "name": "process 5", "description":"send push notification", "priority": 6 }'

curl --location 'http://localhost:8080/job' \
--header 'Content-Type: application/json' \
--data '{ "id": "126", "name": "process 6", "description":"send email notification", "priority": 6 }'
```

*Schedule Jobs*
```
bash
curl --location 'http://localhost:8080/job/123/schedule/fixed' \
--header 'Content-Type: application/json' \
--data '{ "jobId": "123", "minutes": "1", "hours": "0", "daily": "0" }'

curl --location 'http://localhost:8080/job/124/schedule/fixed' \
--header 'Content-Type: application/json' \
--data '{ "jobId": "124", "minutes": "1", "hours": "0", "daily": "0" }'

curl --location 'http://localhost:8080/job/125/schedule/custom' \
--header 'Content-Type: application/json' \
--data '{ "jobId": "125", "minutes": "*/2", "hours": "*", "dayOfMonth": "*", "month": "*", "dayOfWeek": "*" }'
```
