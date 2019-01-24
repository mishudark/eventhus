package eventhus

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
	CommandVersion string    `json:"command_version"`
	AggregateID    string    `json:"aggregate_id"`
	AggregateType  string    `json:"aggregate_type"`
	Type           AlertType `json:"type"`
}
