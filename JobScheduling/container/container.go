package container

import (
	"os"
	"strconv"

	"github.com/vijayshahwal/jobScheduling/controllers"
	"github.com/vijayshahwal/jobScheduling/interfaces"
	"github.com/vijayshahwal/jobScheduling/repositories"
	"github.com/vijayshahwal/jobScheduling/services"
)

type Container struct {
	JobRepository     interfaces.JobRepository
	CacheRepository   interfaces.CacheRepository
	ValidationService interfaces.ValidationService
	JobService        interfaces.JobService
	ScheduleService   interfaces.ScheduleService
	JobController     *controllers.JobController
}

// Helper function to get environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func NewContainer() *Container {
	// Repositories
	mongoURL := getEnv("MONGO_URL", "mongodb://localhost:27017")
	mongoDBName := getEnv("MONGO_DB_NAME", "job_scheduler")

	// Get Redis configuration from environment variables
	redisHost := getEnv("REDIS_HOST", "localhost:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	// Repositories
	jobRepo := repositories.NewMongoJobRepository(mongoURL, mongoDBName)
	cacheRepo := repositories.NewRedisCacheRepository(redisHost, redisPassword, redisDB)

	// Services
	validator := services.NewValidationService()
	jobService := services.NewJobService(jobRepo, validator)
	scheduleService := services.NewScheduleService(cacheRepo, jobRepo, validator)

	// Controllers
	jobController := controllers.NewJobController(jobService, scheduleService)

	return &Container{
		JobRepository:     jobRepo,
		CacheRepository:   cacheRepo,
		ValidationService: validator,
		JobService:        jobService,
		ScheduleService:   scheduleService,
		JobController:     jobController,
	}
}
