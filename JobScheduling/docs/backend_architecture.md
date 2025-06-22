# Job Scheduling System - Backend Architecture

## Overview
The backend is built using Go and follows clean architecture principles with clear separation of concerns. It implements a job scheduling system that supports both fixed-interval and cron-style scheduling.

## Architectural Layers

### 1. Presentation Layer (HTTP/API Layer)
- **Controllers Package**: Handles HTTP requests and responses
  - `JobController`: Manages job and schedule-related endpoints
  - Primary responsibility: Request parsing, response formatting, and routing to appropriate services

### 2. Domain Layer (Core Business Logic)
- **Interfaces Package**: Defines core contracts
  ```go
  // Core service interfaces
  type JobService interface {
      CreateJob(ctx context.Context, job models.Job) (*models.Job, error)
      GetJob(ctx context.Context, id string) (*models.Job, error)
      GetAllJobs(ctx context.Context) ([]models.Job, error)
  }

  type ScheduleService interface {
      ScheduleFixedJob(ctx context.Context, jobID string, schedule models.FixedSchedule) error
      ScheduleCustomJob(ctx context.Context, jobID string, schedule models.CustomSchedule) error
      ProcessSchedules(ctx context.Context) error
  }
  ```

- **Models Package**: Core business entities
  - `Job`: Represents a schedulable task
  - `FixedSchedule`: Interval-based scheduling
  - `CustomSchedule`: Cron-style scheduling

### 3. Service Layer (Application Logic)
- **Services Package**: Implements business logic
  - `JobService`: Job management operations
  - `ScheduleService`: Schedule processing and management
  - `ValidationService`: Input validation
  - Schedule Processors: Handle different schedule types

### 4. Data Layer (Persistence)
- **Repositories Package**: Data access abstraction
  - `JobRepository`: MongoDB-based job storage
  - `CacheRepository`: Redis-based schedule storage
  - Abstracts database operations behind interfaces

## Key Components

### 1. Dependency Container
```go
type Container struct {
    JobRepository     interfaces.JobRepository
    CacheRepository   interfaces.CacheRepository
    ValidationService interfaces.ValidationService
    JobService        interfaces.JobService
    ScheduleService   interfaces.ScheduleService
    JobController     *controllers.JobController
}
```
- Manages dependency injection
- Centralizes component initialization
- Ensures proper dependency wiring

### 2. Job Management
- Job Creation Flow:
  1. HTTP request → JobController
  2. JobController → JobService
  3. JobService → ValidationService (validates input)
  4. JobService → JobRepository (persists job)
  5. MongoDB storage

### 3. Schedule Management
- Two Schedule Types:
  1. Fixed Schedule (interval-based)
  2. Custom Schedule (cron-based i.e used standard cronExpresion)
- Schedule Processing:
  1. Scan Redis for job keys
  2. Load schedule details
  3. Process using appropriate processor
  4. Update next invocation time

### 4. Data Storage
- **MongoDB**
  - Primary job storage
  - Stores job details and metadata
  - Persistent storage

- **Redis**
  - Schedule information caching
  - Fast access to schedule data
  - Efficient schedule processing

## Design Patterns Used

1. **Repository Pattern**
   - Abstracts data access
   - Makes storage implementation swappable
   - Clear separation from business logic

2. **Strategy Pattern**
   - Different schedule types (Fixed/Custom)
   - Common Schedule interface
   - Runtime strategy selection

3. **Dependency Injection**
   - Components receive dependencies
   - Improved testability
   - Loose coupling

4. **Service Layer Pattern**
   - Encapsulates business logic
   - Coordinates between components
   - Single responsibility principle

## Interaction Flow

1. **Job Creation**
   ```
   Client → Controller → JobService → Validation → MongoDB
   ```

2. **Schedule Creation**
   ```
   Client → Controller → ScheduleService → Validation → Redis
   ```

3. **Schedule Processing**
   ```
   Processor → Redis → Type-specific Processor → Redis Update
   ```

## Error Handling
- Each layer handles appropriate errors
- Validation errors at service layer
- Database errors at repository layer
- HTTP errors at controller layer

## Extensibility Points
1. New schedule types can be added by:
   - Implementing Schedule interface
   - Adding new schedule processor
   - Extending ScheduleService

2. Storage can be changed by:
   - Implementing repository interfaces
   - Updating container configuration

## Performance Considerations
1. Redis caching for fast schedule access
2. Efficient schedule processing using processors
3. MongoDB for reliable job persistence
4. Concurrent schedule processing support

The architecture is designed to be:
- Maintainable
- Testable
- Extensible
- Performance-oriented
- Clearly separated in concerns
