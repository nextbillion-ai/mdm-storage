package mdmstorage

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

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

func (s TaskState) String() string {
	res := ""
	switch {
	case s == TPending:
		res = "pending"
	case s == TRunning:
		res = "running"
	case s == TPartiallySucceeded:
		res = "partially_succeeded"
	case s == TFullySucceeded:
		res = "fully_succeeded"
	case s == TFailed:
		res = "failed"
	}
	return res
}

type Task struct {
	ID                    uint32    `gorm:"column:id;primaryKey"`
	TaskID                string    `gorm:"column:task_id"`
	NumOfChunks           uint16    `gorm:"column:num_of_chunks"`
	OutputAddr            string    `gorm:"column:output_addr"`
	OriginalReq           string    `gorm:"column:original_req"`
	ExtractedParams       string    `gorm:"column:extracted_params"`
	State                 TaskState `gorm:"column:state"`
	ResourceAllocatorMeta string    `gorm:"column:resource_allocator_meta"`
	Area                  string    `gorm:"column:area"`
	CreateAt              time.Time `gorm:"column:created_at"`
	PickedUpAt            time.Time `gorm:"column:picked_up_at"`
	FinishedAt            time.Time `gorm:"column:finished_at"`
	RetryTimes            uint8     `gorm:"column:retry_times"`
	CDNAddr               string    `gorm:"column:cdn_addr"`
	Meta                  string    `gorm:"column:meta"`
}

func (t *Task) GetPod() []string {
	if strings.Trim(t.ResourceAllocatorMeta, " ") == "" {
		return nil
	}
	return strings.Split(t.ResourceAllocatorMeta, "|")
}

func (t *Task) SetPod(p []string) {
	t.ResourceAllocatorMeta = strings.Join(p, "|")
}

func (t *Task) GetExtractedParams() *ExtractedParams {
	if t.ExtractedParams == "" {
		return nil
	}

	res := new(ExtractedParams)
	err := json.Unmarshal([]byte(t.ExtractedParams), res)
	if err != nil {
		return nil
	}
	return res
}

type Meta struct {
	FailureReason string   `json:"failure_reason,omitempty"`
	Errors        []string `json:"errors,omitempty"`
}

// ExtractedParams is the struct of extracted_params in tasks table
type ExtractedParams struct {
	DepartureTime      uint64 `json:"departure_time,omitempty"`
	Context            string `json:"context,omitempty"`
	Avoid              string `json:"avoid,omitempty"`
	Key                string `json:"key,omitempty"`
	NbGatewayTrackInfo string `json:"nb-gateway-track-info,omitempty"`
	Spliter            string `json:"spliter,omitempty"`

	// TODO: skip approaches for now
	Approaches string `json:"approaches,omitempty"`

	OriginsApproach      string `json:"origins_approach,omitempty"`
	DestinationsApproach string `json:"destinations_approach,omitempty"`

	RouteType   string `json:"route_type,omitempty"`
	TruckSize   string `json:"truck_size,omitempty"`
	TruckWeight uint32 `json:"truck_weight,omitempty"`
	Option      string `json:"option,omitempty"`
	Mode        string `json:"mode,omitempty"`
	Caller      string `json:"caller,omitempty"`
}

func (ep *ExtractedParams) ToJson() (string, error) {
	if ep == nil {
		return "", errors.New("ExtractedParams can't be nil")
	}
	res, err := json.Marshal(ep)
	return string(res), err
}

func ToExtractedParams(value string) (*ExtractedParams, error) {
	res := new(ExtractedParams)
	err := json.Unmarshal([]byte(value), res)
	return res, err
}

type OriginalReq struct {
	Origins       string `json:"origins"`
	Destinations  string `json:"destinations"`
	Mode          string `json:"mode,omitempty"`
	DepartureTime uint64 `json:"departure_time,omitempty"`
	Context       string `json:"context,omitempty"`
	Avoid         string `json:"avoid,omitempty"`
	Approaches    string `json:"approaches,omitempty"`

	OriginsApproach      string `json:"origins_approach,omitempty"`
	DestinationsApproach string `json:"destinations_approach,omitempty"`

	RouteType   string `json:"route_type,omitempty"`
	TruckSize   string `json:"truck_size,omitempty"`
	TruckWeight uint32 `json:"truck_weight,omitempty"`
	Option      string `json:"option,omitempty"`
	Area        string `json:"area"`
	Key         string `json:"-"`
}

func (t *Task) GetOriginalReq() *OriginalReq {
	if t.OriginalReq == "" {
		return nil
	}

	res := new(OriginalReq)
	err := json.Unmarshal([]byte(t.OriginalReq), res)
	if err != nil {
		return nil
	}
	return res
}

// TableName overrides the table name used by Task to `profiles`
func (Task) TableName() string {
	return "mdm.tasks"
}
