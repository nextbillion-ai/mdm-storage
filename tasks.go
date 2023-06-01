package mdmstorage

import "time"

type TaskState int

const (
	// TPending means the task has just been created, waiting for processing.
	TPending TaskState = iota
	// TRunning means the task has been picked up by the MDM Director and started processing.
	TRunning
	// TPartiallySucceeded means the whole task has been finished but some chunks failed to be processed.
	TPartiallySucceeded
	// TFullySucceeded means the whole task has been finished. All chunks have been successfully processed.
	TFullySucceeded
	// TFailed means the task failed for some reason.
	TFailed
)

type Task struct {
	ID                    uint32    `gorm:"column:id;primaryKey"`
	TaskID                string    `gorm:"column:task_id"`
	NumOfChunks           uint16    `gorm:"column:num_of_chunks"`
	UnitSize              uint16    `gorm:"column:unit_size"`
	OutputAddr            string    `gorm:"column:output_addr"`
	OriginalReq           string    `gorm:"column:original_req"`
	State                 TaskState `gorm:"column:state"`
	ResourceAllocatorMeta string    `gorm:"column:resource_allocator_meta"`
	CreateAt              time.Time `gorm:"column:created_at"`
	PickedUpAt            time.Time `gorm:"column:picked_up_at"`
	FinishedAt            time.Time `gorm:"column:finished_at"`
}
