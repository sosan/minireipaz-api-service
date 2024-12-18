package repos

import (
	"minireipaz/pkg/domain/models"
	"time"

	"github.com/google/uuid"
)

type WorkflowService interface {
	CreateWorkflow(workflowFrontend *models.WorkflowFrontend) (created bool, exist bool, workflow *models.Workflow)
	GetWorkflow(userID, workflowID *string) (newWorkflow *models.Workflow, exist bool)
	GetAllWorkflows(userID *string) (allWorkflows []models.Workflow, err error)
	UpdateWorkflow(workflow *models.Workflow) (updated bool, exist bool)
	ValidateWorkflowGlobalUUID(uuid *string) bool
	ValidateUserWorkflowUUID(worklfowID, name *string) bool
}

type WorkflowRedisRepoInterface interface {
	Create(workflow *models.Workflow) (created bool, exist bool)
	Update(worflow *models.Workflow) (updated bool, exist bool)
	Remove(workflow *models.Workflow) (removed bool)
	ValidateWorkflowGlobalUUID(uuid *string) bool
	ValidateUserWorkflowUUID(userID, name *string) bool
	GetByUUID(id uuid.UUID) (*models.Workflow, error)
	AcquireLock(key, value string, expiration time.Duration) (locked bool, err error)
	RemoveLock(key string) bool
}

type WorkflowBrokerRepository interface {
	Create(workflow *models.Workflow) (sended bool)
	Update(workflow *models.Workflow) (sended bool)
}

type WorkflowHTTPRepository interface {
	GetWorkflowDataByID(userID, workflowID *string, limitCount uint64) (*models.InfoWorkflow, error)
	GetAllWorkflows(userID *string, limitCount uint64) (*models.InfoWorkflow, error)
}
