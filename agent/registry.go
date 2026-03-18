package agent

import (
	"errors"
	"fmt"
)

var ErrAgentNotFound = errors.New("agent not found")
var ErrAgentExists = errors.New("agent already exists")

// Registry holds all agents in memory.
// Persistence is handled by the store package.
type Registry struct {
	Agents []*Agent `json:"agents"`
}

func NewRegistry() *Registry {
	return &Registry{Agents: make([]*Agent, 0)}
}

func (r *Registry) Get(id string) (*Agent, error) {
	for _, a := range r.Agents {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, fmt.Errorf("%w: %s", ErrAgentNotFound, id)
}

func (r *Registry) Add(a *Agent) error {
	for _, existing := range r.Agents {
		if existing.ID == a.ID {
			return fmt.Errorf("%w: %s", ErrAgentExists, a.ID)
		}
	}
	r.Agents = append(r.Agents, a)
	return nil
}

func (r *Registry) Update(a *Agent) error {
	for i, existing := range r.Agents {
		if existing.ID == a.ID {
			r.Agents[i] = a
			return nil
		}
	}
	return fmt.Errorf("%w: %s", ErrAgentNotFound, a.ID)
}

func (r *Registry) Delete(id string) error {
	for i, a := range r.Agents {
		if a.ID == id {
			r.Agents = append(r.Agents[:i], r.Agents[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("%w: %s", ErrAgentNotFound, id)
}

// List returns agents filtered by status. No filter = all agents.
func (r *Registry) List(statuses ...Status) []*Agent {
	if len(statuses) == 0 {
		result := make([]*Agent, len(r.Agents))
		copy(result, r.Agents)
		return result
	}
	statusSet := make(map[Status]bool)
	for _, s := range statuses {
		statusSet[s] = true
	}
	var result []*Agent
	for _, a := range r.Agents {
		if statusSet[a.Status] {
			result = append(result, a)
		}
	}
	return result
}
