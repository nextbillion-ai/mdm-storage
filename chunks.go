package mdmstorage

import (
	"encoding/json"
	"fmt"
	"time"
)

type ChunkState int

const (
	// CPending means the task has just been created, waiting for processing.
	CPending ChunkState = iota
	// CRunning means DMD Director successfully got a pod and dispatched the sub-task to it.
	CRunning
	// CSucceeded means the entire task has been finished and has uploaded the output file to the correct Cloudflare address.
	CSucceeded
	// CFailed means the sub-task failed for some reason.
	CFailed
)

func (s ChunkState) String() string {
	res := ""
	switch {
	case s == CPending:
		res = "pending"
	case s == CRunning:
		res = "running"
	case s == CSucceeded:
		res = "successed"
	case s == CFailed:
		res = "failed"
	}
	return res
}

type Chunk struct {
	ID               uint32     `gorm:"column:id;primaryKey"`
	ChunkIndex       uint32     `gorm:"column:chunk_index"`
	TaskID           string     `gorm:"column:task_id"`
	State            ChunkState `gorm:"column:state"`
	ResourceInfo     string     `gorm:"column:resource_info"`
	RetryTimes       uint8      `gorm:"column:retry_times"`
	FinishedAt       time.Time  `gorm:"column:finished_at"`
	StartedAt        time.Time  `gorm:"column:started_at"`
	Origins          string     `gorm:"column:origins"`
	Destinations     string     `gorm:"column:destinations"`
	OriginIndex      string     `gorm:"column:origin_index"`
	DestinationIndex string     `gorm:"column:destination_index"`
	CreateAt         time.Time  `gorm:"column:created_at"`
	MetaStr          string     `gorm:"column:meta"`
	Meta             *Meta      `gorm:"-"`
}

func (c *Chunk) FlattenToString() {
	// flatten meta to string
	if c.Meta == nil {
		return
	}
	res, err := json.Marshal(c.Meta)
	if err != nil {
		c.MetaStr = fmt.Sprintf("failed to marshal task meta: %v", err)
		return
	}
	c.MetaStr = string(res)
}

// TableName overrides the table name used by Task to `profiles`
func (c *Chunk) TableName() string {
	return "mdm.chunks"
}

func (c *Chunk) SetFailureReason(reason string) {
	if c.Meta == nil {
		c.Meta = &Meta{}
	}
	c.MetaStr = string(reason)
}
