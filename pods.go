package mdmstorage

import (
	"encoding/json"
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
	ResourceMetaStr string    `gorm:"column:resource_meta"`
	State           PodState  `gorm:"column:state"`
	RetryTimes      uint8     `gorm:"column:retry_times"`
	Meta            string    `gorm:"column:meta"`
	CreateAt        time.Time `gorm:"column:created_at"`
	UpdateAt        time.Time `gorm:"column:update_at"`
	FinishedAt      time.Time `gorm:"column:finished_at"`

	ResourceMeta    *ResourceMeta `gorm:"-"`
	CurrentJobCount int           `gorm:"-"`
}

func (p *PodInfo) Match(mode, area, option string) bool {
	return p.Mode == mode && p.Area == area && p.Option == option
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
	Asset struct {
		Type    string `json:"type"`
		Release string `json:"release"`
	} `json:"asset"`
	App struct {
		Area     string `json:"area"`
		Mode     string `json:"mode"`
		Replicas int    `json:"replicas"`
		Parallel bool   `json:"parallel"`
		Version  string `json:"version"`
		Storage  string `json:"storage"`
	} `json:"app"`
}

func (p *PodInfo) TableName() string {
	return "mdm.pods"
}
