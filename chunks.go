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

func (s ChunkState) String() string {
	res := ""
	switch {
	case s == CPending:
		res = "pending"
	case s == CResourceCreating:
		res = "resource_creating"
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
	Meta             string     `gorm:"column:meta"`
}

// TableName overrides the table name used by Task to `profiles`
func (Chunk) TableName() string {
	return "mdm.chunks"
}

/*
-- Table: mdm.chunks

-- DROP TABLE IF EXISTS mdm.chunks;

CREATE TABLE IF NOT EXISTS mdm.chunks
(
    id bigint NOT NULL DEFAULT nextval('mdm.chunks_id_seq'::regclass),
    chunk_index integer,
    task_id character varying(255) COLLATE pg_catalog."default" NOT NULL,
    state integer,
    resource_info text COLLATE pg_catalog."default",
    retry_times integer,
    finished_at timestamp without time zone,
    origins text COLLATE pg_catalog."default",
    destinations text COLLATE pg_catalog."default",
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    origin_index text COLLATE pg_catalog."default",
    destination_index text COLLATE pg_catalog."default",
    started_at timestamp without time zone,
    meta text COLLATE pg_catalog."default",
    CONSTRAINT chunks_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS mdm.chunks
    OWNER to fangzhou;
*/
