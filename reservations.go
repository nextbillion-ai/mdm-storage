package mdmstorage

import (
	"strings"
	"time"
)

type ReservationState int

const (
	RPending ReservationState = iota
	RApplying
	RRunning
	RFinished
	RFailed
)

// Reservation table reservations
type Reservation struct {
	ID          int              `gorm:"column:id"`
	TaskID      string           `gorm:"column:task_id"`
	TaskIDNum   uint32           `gorm:"column:task_id_num"`
	Area        string           `gorm:"column:area"`
	Mode        string           `gorm:"column:mode"`
	Option      string           `gorm:"column:option"`
	Status      ReservationState `gorm:"column:status"`
	ChunkString string           `gorm:"column:chunk"`
	PodString   string           `gorm:"column:pod"`
	Meta        string           `gorm:"column:meta"`
	CreateAt    time.Time        `gorm:"column:created_at"`
	FinishedAt  time.Time        `gorm:"column:finished_at"`
}

func (r *Reservation) Chunk() []string {
	return strings.Split(r.ChunkString, "|")
}

func (r *Reservation) Pod() []string {
	return strings.Split(r.PodString, "|")
}

func (Reservation) TableName() string {
	return "mdm.reservations"
}

/*
-- Table: mdm.reservations

-- DROP TABLE IF EXISTS mdm.reservations;

create sequence mdm.reservations_id_seq;
alter sequence mdm.reservations_id_seq owner to fangzhou;

CREATE TABLE IF NOT EXISTS mdm.reservations
(
    id bigint NOT NULL DEFAULT nextval('mdm.reservations_id_seq'::regclass),
    task_id character varying(255) COLLATE pg_catalog."default" NOT NULL,
    option character varying(255) COLLATE pg_catalog."default" NOT NULL,
    mode character varying(255) COLLATE pg_catalog."default" NOT NULL,
    area character varying(255) COLLATE pg_catalog."default" NOT NULL,
    state integer,
    chunk text COLLATE pg_catalog."default",
    pod text COLLATE pg_catalog."default",
    meta text COLLATE pg_catalog."default",
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    finished_at timestamp without time zone,
    CONSTRAINT reservations_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS mdm.reservations
    OWNER to fangzhou;
*/
