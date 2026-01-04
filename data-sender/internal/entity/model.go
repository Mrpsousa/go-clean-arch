package entity

import (
	"time"
)

type RabbitMsg struct {
	CreatedAt    time.Time
	ExameName    string
	PacientName string
	PacientNumbPhone string
	DocImagePath string
	
}
