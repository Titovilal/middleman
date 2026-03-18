package agent

import "time"

type Status string

const (
	StatusActive    Status = "active"
	StatusIdle      Status = "idle"
	StatusWorking   Status = "working"
	StatusDiscarded Status = "discarded"
)

// Agent is the registry record for a single AI agent instance.
type Agent struct {
	ID            string `json:"id"`
	ConnectorName string `json:"connector_name"`
	SessionID     string `json:"session_id"`
	WorkDir       string `json:"work_dir"`
	Status        Status `json:"status"`

	// Briefing is the initial context given at spawn time.
	Briefing string `json:"briefing"`

	// KnownContext is a free-text summary the Middleman maintains about
	// what this agent knows: files read, decisions made, etc.
	KnownContext string `json:"known_context,omitempty"`

	// Checkpoints is an ordered list of save points (oldest first).
	Checkpoints []CheckpointRecord `json:"checkpoints"`

	// TaskLog is an append-only record of every delegation.
	TaskLog []TaskRecord `json:"task_log"`

	CreatedAt    time.Time `json:"created_at"`
	LastActiveAt time.Time `json:"last_active_at"`
}

// CheckpointRecord is a persisted snapshot of a session state.
type CheckpointRecord struct {
	Label string `json:"label"`

	// TurnIndex is the assistant turn count at checkpoint time.
	TurnIndex int `json:"turn_index"`

	// NativeRef is connector-specific:
	// Claude: UUID of the last assistant message.
	// Gemini: numeric session index as string.
	NativeRef string `json:"native_ref"`

	CreatedAt time.Time `json:"created_at"`
}

// TaskRecord is one entry in the agent's task log.
type TaskRecord struct {
	TaskID      string    `json:"task_id"`
	Prompt      string    `json:"prompt"`
	Response    string    `json:"response"`
	IsError     bool      `json:"is_error"`
	ErrorDetail string    `json:"error_detail,omitempty"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
}

// LatestCheckpoint returns the most recent checkpoint, or nil if none.
func (a *Agent) LatestCheckpoint() *CheckpointRecord {
	if len(a.Checkpoints) == 0 {
		return nil
	}
	cp := a.Checkpoints[len(a.Checkpoints)-1]
	return &cp
}

// CheckpointByLabel returns the checkpoint with the given label, or nil.
func (a *Agent) CheckpointByLabel(label string) *CheckpointRecord {
	for i := range a.Checkpoints {
		if a.Checkpoints[i].Label == label {
			return &a.Checkpoints[i]
		}
	}
	return nil
}
