package state

import (
	"encoding/json"
	"fmt"
	"os"
)

type FeatureStatus string

const (
	StatusPending   FeatureStatus = "pending"
	StatusCompleted FeatureStatus = "completed"
)

type Feature struct {
	Name                  string          `json:"name"`
	Files                 []string        `json:"files"`
	SubFeatures           []string        `json:"sub_features,omitempty"`
	ImplementationDetails []string        `json:"implementation_details,omitempty"`
	Dependencies          []string        `json:"dependencies,omitempty"`
	Status                FeatureStatus   `json:"status"`
}

type DocTornState struct {
	Version           int       `json:"version"`
	ProjectRoot       string    `json:"project_root"`
	Features          []Feature `json:"features"`
	Order             []string  `json:"order"`
	CurrentFeature    string    `json:"current_feature"`
	MetaDocsGenerated bool      `json:"meta_docs_generated"`
}

const stateFileName = ".doc-torn-state.json"

func stateFilePath() string {
	return stateFileName
}

func Read() (*DocTornState, error) {
	data, err := os.ReadFile(stateFilePath())
	if err != nil {
		return nil, fmt.Errorf("reading state: %w", err)
	}
	var s DocTornState
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing state: %w", err)
	}
	if err := s.validate(); err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *DocTornState) Write() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling state: %w", err)
	}
	return os.WriteFile(stateFilePath(), data, os.FileMode(0644))
}

func (s *DocTornState) validate() error {
	validStatuses := map[FeatureStatus]bool{
		StatusPending:   true,
		StatusCompleted: true,
	}
	for _, f := range s.Features {
		if !validStatuses[f.Status] {
			return fmt.Errorf("feature %q has invalid status %q", f.Name, f.Status)
		}
	}
	return nil
}

func (s *DocTornState) Complete(name string) error {
	found := false
	for i, f := range s.Features {
		if f.Name == name {
			s.Features[i].Status = StatusCompleted
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("feature %q not found in state", name)
	}
	// Find next pending feature
	s.CurrentFeature = ""
	for _, orderName := range s.Order {
		for _, f := range s.Features {
			if f.Name == orderName && f.Status == StatusPending {
				s.CurrentFeature = orderName
				return s.Write()
			}
		}
	}
	return s.Write()
}

func (s *DocTornState) AllCompleted() bool {
	for _, f := range s.Features {
		if f.Status != StatusCompleted {
			return false
		}
	}
	return true
}
