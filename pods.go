package mdmstorage

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type PodState int

const (
	PApplying PodState = iota
	PRunning
	PFailed
	PRemoved
	PCached
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
	ResourceMetaStr string    `gorm:"column:resource_meta"`
	State           PodState  `gorm:"column:state"`
	RetryTimes      uint8     `gorm:"column:retry_times"`
	Meta            string    `gorm:"column:meta"`
	CreateAt        time.Time `gorm:"column:created_at"`
	UpdateAt        time.Time `gorm:"column:update_at"`
	FinishedAt      time.Time `gorm:"column:finished_at"`

	ResourceMeta *ResourceMeta `gorm:"-"`
}

func (p *PodInfo) Match(mode, area, option string) bool {
	return p.Mode == mode && p.Area == area && p.Option == option
}

func (p *PodInfo) GetCurrentJob() []string {
	if strings.Trim(p.CurrentJob, " ") == "" {
		return make([]string, 0)
	}
	return strings.Split(p.CurrentJob, "|")
}

func (p *PodInfo) AppendCurrentJob(taskID string, chunkID uint32) {
	job := fmt.Sprintf("%v::%v", taskID, chunkID)
	if strings.Trim(p.CurrentJob, " ") == "" {
		p.CurrentJob = job
		return
	}
	currentJobs := strings.Split(p.CurrentJob, "|")
	currentJobs = append(currentJobs, job)
	p.CurrentJob = strings.Join(currentJobs, "|")
}

func (p *PodInfo) RemoveCurrentJob(taskID string, chunkID uint32) int {
	job := fmt.Sprintf("%v::%v", taskID, chunkID)
	currentJob := strings.Split(p.CurrentJob, "|")
	removeCount := 0
	newCurrentJob := make([]string, 0)
	for _, v := range currentJob {
		if v == job {
			removeCount++
			continue
		}
		newCurrentJob = append(newCurrentJob, v)
	}
	p.CurrentJob = strings.Join(newCurrentJob, "|")
	return removeCount
}

func (p *PodInfo) GetResourceMeta() *ResourceMeta {
	if p.ResourceMetaStr == "" {
		return nil
	}
	if p.ResourceMeta == nil {
		res := new(ResourceMeta)
		err := json.Unmarshal([]byte(p.ResourceMetaStr), res)
		if err != nil {
			return nil
		}
		p.ResourceMeta = res
	}
	return p.ResourceMeta
}

type ResourceMeta struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Action    string `json:"action"`
	Spec      struct {
		Asset struct {
			Type    string `json:"type"`
			Release string `json:"release"`
		} `json:"asset"`
		App struct {
			Replicas int    `json:"replicas"`
			Parallel bool   `json:"parallel"`
			Version  string `json:"version"`
			Storage  string `json:"storage"`
		} `json:"app"`
	} `json:"spec"`
}

func (p *PodInfo) TableName() string {
	return "mdm.pods"
}
