PREREQISITES:
1. Please install latest stable version of Golang and docker
2. Verify that docker engine is running
3. cURL is installed
4. For Detailed architecture refer to docs/backend_architecture.md


Steps to execute the code:
1. git clone git@github.com:vijayshahwal/jobScheduling.git
2. cd jobScheduling/
3. go mod init github.com/vijayshahwal/jobScheduling
4. go mod tidy
5. docker run -d --name redis -p 6379:6379 -e REDIS_HOST="localhost:6379" -e REDIS_PASSWORD="" -e REDIS_DB="0" redis:latest
6. docker run -d --name mongo -p 27017:27017 -e MONGO_URL="mongodb://localhost:27017" -e MONGO_DB_NAME="job_scheduler" mongo:latest
7. MONGO_URL="mongodb://localhost:27017" MONGO_DB_NAME="job_scheduler" REDIS_HOST="localhost:6379" REDIS_PASSWORD="" REDIS_DB="0" go run main.go
8. open new integrated terminal in same folder and execute the below command

## --- createJob --- ##
curl --location 'http://localhost:8080/job' \
--header 'Content-Type: application/json' \
--data '{
    "id": "123",
    "name": "process 3",
    "description":"send watsapp notification",
    "priority": 6
}'

curl --location 'http://localhost:8080/job' \
--header 'Content-Type: application/json' \
--data '{
    "id": "124",
    "name": "process 4",
    "description":"send sms notification",
    "priority": 6
}'

curl --location 'http://localhost:8080/job' \
--header 'Content-Type: application/json' \
--data '{
    "id": "125",
    "name": "process 5",
    "description":"send push notification",
    "priority": 6
}'

curl --location 'http://localhost:8080/job' \
--header 'Content-Type: application/json' \
--data '{
    "id": "126",
    "name": "process 6",
    "description":"send email notification",
    "priority": 6
}'


## ---- Schedule JOB --- ##

curl --location 'http://localhost:8080/job/123/schedule/fixed' \
--header 'Content-Type: application/json' \
--data '{
    "jobId": "123",
    "minutes": "1",
    "hours": "0",
    "daily": "0"
}'

curl --location 'http://localhost:8080/job/124/schedule/fixed' \
--header 'Content-Type: application/json' \
--data '{
    "jobId": "124",
    "minutes": "1",
    "hours": "0",
    "daily": "0"
}'

curl --location 'http://localhost:8080/job/125/schedule/custom' \
--header 'Content-Type: application/json' \
--data '{
    "jobId": "125",
    "minutes": "*/2",
    "hours": "*",
    "dayOfMonth": "*",
    "month": "*",
    "dayOfWeek": "*"
}'
