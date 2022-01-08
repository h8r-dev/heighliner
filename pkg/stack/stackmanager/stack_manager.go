package stackmanager

import (
	"errors"

	"github.com/h8r-dev/heighliner/pkg/models"
	"github.com/h8r-dev/heighliner/pkg/stack/stackhandler"
)

type StackManager interface {
	GetStack(id string) *models.Stack
	ListStackNames() []string
	InstantiateStack(id string, stackID string, params map[string]interface{}) (*models.Application, error)
}

type stackManager struct {
	stacks   map[string]*models.Stack
	handlers map[string]stackhandler.StackHandler
}

func (s *stackManager) ListStackNames() []string {
	var names []string
	for id := range s.stacks {
		names = append(names, id)
	}
	return names
}

func (s *stackManager) GetStack(id string) *models.Stack {
	return s.stacks[id]
}

func (s *stackManager) InstantiateStack(id string, stackID string, params map[string]interface{}) (*models.Application, error) {
	_, ok := s.handlers[stackID]
	if !ok {
		return nil, errors.New("stack id not found")
	}
	app := &models.Application{
		ID:      id,
		StackID: stackID,
	}
	return app, nil
}

func New() StackManager {
	return &stackManager{
		stacks:   defaultStacks(),
		handlers: defaultHandlers(),
	}
}

func defaultHandlers() map[string]stackhandler.StackHandler {
	return map[string]stackhandler.StackHandler{
		"microservice": nil,
	}
}

func defaultStacks() map[string]*models.Stack {
	return map[string]*models.Stack{
		"microservice": {
			JSONSchema: `{
  "$id": "https://example.com/person.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Person",
  "type": "object",
  "properties": {
    "firstName": {
      "type": "string",
      "description": "The person's first name."
    },
    "lastName": {
      "type": "string",
      "description": "The person's last name."
    },
    "age": {
      "description": "Age in years which must be equal to or greater than zero.",
      "type": "integer",
      "minimum": 0
    }
  }
}
`,
		},
	}
}
