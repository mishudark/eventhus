package eventhus

import "fmt"

// AlertType defines the alert(error) type while a command is being processed
type AlertType string

// nolint
const (
	AlertLoadingEventts    AlertType = "loading_events"
	AlertReplayingEvents   AlertType = "replying_events"
	AlertProcessingCommand AlertType = "processing_command"
	AlertInvalidID         AlertType = "invalid_id"
	AlertSavingOnStorage   AlertType = "saving_on_storage"
)

// Alert is an error while the command is being processed
type Alert struct {
	CommandID      string    `json:"command_id"`
	CommandType    string    `json:"command_type"`
	CommandVersion int       `json:"command_version"`
	AggregateID    string    `json:"aggregate_id"`
	AggregateType  string    `json:"aggregate_type"`
	Type           AlertType `json:"type"`
}

// NewAlert returns an alert that implements an error interface
func NewAlert(typ AlertType, command Command) Alert {
	return Alert{
		CommandID:      command.GetID(),
		CommandType:    command.GetType(),
		CommandVersion: command.GetVersion(),
		AggregateID:    command.GetAggregateID(),
		AggregateType:  command.GetAggregateType(),
		Type:           typ,
	}
}

func (a *Alert) Error() string {
	return fmt.Sprintf("[%s]: command-id=%s command-type=%s command-version=%d aggregate-id:%s aggregate_type=%s",
		a.Type,
		a.CommandID,
		a.CommandType,
		a.CommandVersion,
		a.AggregateID,
		a.AggregateType)
}

var _ error = (*Alert)(nil)
