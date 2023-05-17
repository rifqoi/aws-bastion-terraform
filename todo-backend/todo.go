package main

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/rifqoi/aws-project/todo-backend/common"
)

type Todo struct {
	ID     uuid.UUID `dynamodbav:"id,omitempty"`
	Task   string    `dynamodbav:"task"`
	Status bool      `dynamodbav:"status"`
}

func NewTodo(task string, status bool) Todo {
	return Todo{
		ID:     uuid.New(),
		Task:   task,
		Status: status,
	}
}

// UnmarshalDynamoDBAttributeValue is called when using attributevalue.Unmarshal method
// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue#Unmarshaler
func (t *Todo) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	mapValue, ok := av.(*types.AttributeValueMemberM)
	if !ok {
		return nil
	}

	for k, v := range mapValue.Value {
		if k == "id" {
			idString, ok := v.(*types.AttributeValueMemberS)
			if !ok {
				return common.NewErrorf(common.ErrorTypeServerError, "failed to unmarshal uuid, type casting error")
			}
			id, err := uuid.Parse(idString.Value)
			if err != nil {
				return common.WrapErrorf(err, common.ErrorTypeServerError, "failed to unmarshal uuid")
			}

			t.ID = id
		} else if k == "task" {
			taskString, ok := v.(*types.AttributeValueMemberS)
			if !ok {
				return common.NewErrorf(common.ErrorTypeServerError, "failed to unmarshal task, type casting error")
			}

			t.Task = taskString.Value
		} else if k == "status" {
			statusString, ok := v.(*types.AttributeValueMemberBOOL)
			if !ok {
				return common.NewErrorf(common.ErrorTypeServerError, "failed to unmarshal status, type casting error")
			}

			t.Status = statusString.Value
		}
	}

	return nil
}

func (t *Todo) GetTask() string {
	return t.Task
}

func (t *Todo) SetDone() bool {
	t.Status = true
	return t.Status
}

func (t *Todo) SetNotDone() bool {
	t.Status = false
	return t.Status
}

func (t *Todo) ToggleDone() {
	t.Status = !t.Status
}

func (t *Todo) GetStatus() bool {
	return t.Status
}

func (t *Todo) IsDone() bool {
	return t.Status == true
}
