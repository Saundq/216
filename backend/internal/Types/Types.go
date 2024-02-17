package Types

import "github.com/gofrs/uuid"

type Status string

const (
	ALIVE        Status = "Alive"
	NOTAVAILABLE Status = "Not available"
)

type Calculator struct {
	Id         uuid.UUID
	Name       string
	Task       uuid.UUID
	Expression string
	Status     Status
	HeartBeat  int
}
