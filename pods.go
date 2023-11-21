package mdmstorage

import (
	"strings"
	"time"
)

type PodState int

const (
	PApplying PodState = iota
	PRunning
	PFailed
	PRemoved
)

// PodInfo table pods
type PodInfo struct {
	ID              int       `gorm:"column:id"`
	Name            string    `gorm:"column:name"`
	CPU             int       `gorm:"column:cpu"`
	Memory          int       `gorm:"column:memory"`
	Storage         string    `gorm:"column:storage"`
	Option          string    `gorm:"column:option"`
	Mode            string    `gorm:"column:mode"`
	Area            string    `gorm:"column:area"`
	CurrentJobCount int       `gorm:"column:current_job_count"`
	CurrentJob      string    `gorm:"column:current_job"`
	ResourceMeta    string    `gorm:"column:resource_meta"`
	State           PodState  `gorm:"column:state"`
	RetryTimes      uint8     `gorm:"column:retry_times"`
	Meta            string    `gorm:"column:meta"`
	CreateAt        time.Time `gorm:"column:created_at"`
	FinishedAt      time.Time `gorm:"column:finished_at"`
}

func (p *PodInfo) AvailableQuota() int {
	return p.CPU - p.CurrentJobCount
}

func (p *PodInfo) AppendCurrentJob(job string) {
	currentJob := strings.Split(p.CurrentJob, "|")
	currentJob = append(currentJob, job)
	p.CurrentJob = strings.Join(currentJob, "|")
	p.CurrentJobCount++
}

func (p *PodInfo) RemoveCurrentJob(taskID string) int {
	currentJob := strings.Split(p.CurrentJob, "|")

	removeCount := 0
	newCurrentJob := make([]string, 0)
	for _, v := range currentJob {
		if strings.HasPrefix(v, taskID) {
			removeCount++
			p.CurrentJobCount -= 1
			continue
		}
		newCurrentJob = append(newCurrentJob, v)
	}
	p.CurrentJob = strings.Join(currentJob, "|")
	return removeCount
}

func (PodInfo) TableName() string {
	return "mdm.pods"
}

/*
-- Table: mdm.pods

DROP TABLE IF EXISTS mdm.pods;

create sequence mdm.pods_id_seq;
alter sequence mdm.pods_id_seq owner to fangzhou;

CREATE TABLE IF NOT EXISTS mdm.pods
(
    id bigint NOT NULL DEFAULT nextval('mdm.pods_id_seq'::regclass),
    name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    cpu integer,
    memory integer,
    option character varying(255) COLLATE pg_catalog."default" NOT NULL,
    mode character varying(255) COLLATE pg_catalog."default" NOT NULL,
    area character varying(255) COLLATE pg_catalog."default" NOT NULL,
    current_job_count integer,
    current_job  text COLLATE pg_catalog."default",
    resource_meta text COLLATE pg_catalog."default",
    state integer,
    retry_times integer,
    meta text COLLATE pg_catalog."default",
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    finished_at timestamp without time zone,
    CONSTRAINT pods_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS mdm.pods
    OWNER to fangzhou;
*/
