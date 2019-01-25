package eventhus

import "fmt"

// FailureType defines the alert(error) type while a command is being processed
type FailureType string

// nolint
const (
	FailureLoadingEventts    FailureType = "loading_events"
	FailureReplayingEvents   FailureType = "replying_events"
	FailureProcessingCommand FailureType = "processing_command"
	FailureInvalidID         FailureType = "invalid_id"
	FailureSavingOnStorage   FailureType = "saving_on_storage"
)

// Failure is an error while the command is being processed
type Failure struct {
	CommandID      string      `json:"command_id"`
	CommandType    string      `json:"command_type"`
	CommandVersion int         `json:"command_version"`
	AggregateID    string      `json:"aggregate_id"`
	AggregateType  string      `json:"aggregate_type"`
	Type           FailureType `json:"type"`
}

// NewFailure returns an alert that implements an error interface
func NewFailure(typ FailureType, command Command) Failure {
	return Failure{
		CommandID:      command.GetID(),
		CommandType:    command.GetType(),
		CommandVersion: command.GetVersion(),
		AggregateID:    command.GetAggregateID(),
		AggregateType:  command.GetAggregateType(),
		Type:           typ,
	}
}

func (a *Failure) Error() string {
	return fmt.Sprintf("[%s]: command-id=%s command-type=%s command-version=%d aggregate-id:%s aggregate_type=%s",
		a.Type,
		a.CommandID,
		a.CommandType,
		a.CommandVersion,
		a.AggregateID,
		a.AggregateType)
}

var _ error = (*Failure)(nil)
