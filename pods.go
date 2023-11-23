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
	UpdateAt        time.Time `gorm:"column:update_at"`
	FinishedAt      time.Time `gorm:"column:finished_at"`
}

func (p *PodInfo) Match(task *Task) bool {
	ep := task.GetExtractedParams()
	return p.Mode == ep.Mode && p.Area == task.Area && p.Option == ep.Option
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
	p.CurrentJob = strings.Join(newCurrentJob, "|")
	return removeCount
}

func (PodInfo) TableName() string {
	return "mdm.pods"
}
