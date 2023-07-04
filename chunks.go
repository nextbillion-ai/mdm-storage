package mdmstorage

import "time"

type ChunkState int

const (
	// CPending means the task has just been created, waiting for processing.
	CPending ChunkState = iota
	// CResourceCreating means the task has been picked up by the MDM Director and started processing.
	CResourceCreating
	// CRunning means DMD Director successfully got a pod and dispatched the sub-task to it.
	CRunning
	// CSucceeded means the entire task has been finished and has uploaded the output file to the correct Cloudflare address.
	CSucceeded
	// CFailed means the sub-task failed for some reason.
	CFailed
)

type Chunk struct {
	ID           uint32     `gorm:"column:id;primaryKey"`
	ChunkIndex   uint32     `gorm:"column:chunk_index"`
	TaskID       string     `gorm:"column:task_id"`
	State        ChunkState `gorm:"column:state"`
	ResourceInfo string     `gorm:"column:resource_info"`
	RetryTimes   uint8      `gorm:"column:retry_times"`
	FinishedAt   time.Time  `gorm:"column:finished_at"`
	Origins      string     `gorm:"column:origins"`
	Destinations string     `gorm:"column:destinations"`
}

// TableName overrides the table name used by Task to `profiles`
func (Chunk) TableName() string {
	return "mdm.chunks"
}
