package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/rifqoi/aws-project/todo-backend/common"
)

type TodoRepository interface {
	GetTaskByID(ctx context.Context, id uuid.UUID) (*Todo, error)
	GetTasks(ctx context.Context) ([]Todo, error)
	AddTask(ctx context.Context, task Todo) error
	UpdateTask(ctx context.Context, updatedTask Todo) error
	RemoveTask(ctx context.Context, id uuid.UUID) error
}

type DynamoRepo struct {
	db *dynamodb.Client
}

func NewDynamoRepo(db *dynamodb.Client) TodoRepository {
	if db == nil {
		panic("db is nil")
	}

	return &DynamoRepo{db}
}

const table = "todo"

func (dr *DynamoRepo) GetTaskByID(ctx context.Context, id uuid.UUID) (*Todo, error) {
	out, err := dr.db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id.String()},
		},
	})

	if err != nil {
		var notFoundExc *types.ResourceNotFoundException
		if errors.As(err, &notFoundExc) {
			return nil, common.WrapErrorf(err, common.ErrorTypeServerError, "DynamoRepo.GetTaskByID: task not found")

		}

		return nil, common.WrapErrorf(err, common.ErrorTypeUnknown, "DynamoRepo.GetTaskByID: failed to inserting item to dynamodb")
	}

	var todo Todo
	err = attributevalue.UnmarshalMap(out.Item, &todo)
	if err != nil {
		return nil, common.WrapErrorf(err, common.ErrorTypeUnknown, "DynamoRepo.GetTaskByID: failed to unmarshal value")
	}

	return &todo, nil
}

func (dr *DynamoRepo) GetTasks(ctx context.Context) ([]Todo, error) {
	out, err := dr.db.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(table),
	})

	if err != nil {
		return nil, common.WrapErrorf(err, common.ErrorTypeUnknown, "DynamoRepo.GetTasks: failed to get items")
	}

	todos := []Todo{}
	for _, item := range out.Items {
		var todo Todo
		err := attributevalue.UnmarshalMap(item, &todo)
		if err != nil {
			return nil, common.WrapErrorf(err, common.ErrorTypeServerError, "DynamoRepo.GetTasks: failed unmarshal value")
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func (dr *DynamoRepo) AddTask(ctx context.Context, todo Todo) error {
	_, err := dr.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item: map[string]types.AttributeValue{
			"id":        &types.AttributeValueMemberS{Value: uuid.NewString()},
			"task":      &types.AttributeValueMemberS{Value: todo.GetTask()},
			"status":    &types.AttributeValueMemberBOOL{Value: todo.GetStatus()},
			"createdAt": &types.AttributeValueMemberS{Value: time.Now().String()},
			"updatedAt": &types.AttributeValueMemberS{Value: time.Now().String()},
		},
	})

	if err != nil {
		var notFoundExc *types.ResourceNotFoundException
		if errors.As(err, &notFoundExc) {
			return common.WrapErrorf(err, common.ErrorTypeServerError, "DynamoRepo.AddTask: table not found")

		}
		return common.WrapErrorf(err, common.ErrorTypeUnknown, "DynamoRepo.AddTask: failed to inserting item to dynamodb")
	}

	return nil
}

func (dr *DynamoRepo) UpdateTask(ctx context.Context, updatedTodo Todo) error {
	expr, err := expression.NewBuilder().WithUpdate(
		expression.Set(
			expression.Name("task"),
			expression.Value(updatedTodo.Task),
		).Set(
			expression.Name("status"),
			expression.Value(updatedTodo.Status),
		),
	).Build()
	fmt.Println(updatedTodo)
	fmt.Println(expr.Names())

	if err != nil {
		return common.WrapErrorf(err, common.ErrorTypeUnknown, "DynamoRepo.UpdateTask: failed to build expression")
	}

	_, err = dr.db.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: updatedTodo.ID.String()},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       expr.Condition(),
	})

	if err != nil {
		return common.WrapErrorf(err, common.ErrorTypeServerError, "DynamoRepo.UpdateTask: failed to update item in dynamodb")
	}

	return nil
}

func (dr *DynamoRepo) RemoveTask(ctx context.Context, id uuid.UUID) error {

	_, err := dr.db.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id.String()},
		},
	})

	if err != nil {
		return common.WrapErrorf(err, common.ErrorTypeUnknown, "DynamoRepo.RemoveTask: failed to remove task")
	}

	return nil
}

// func main() {
// 	dyDB := config.NewDynamoDB()
//
// 	dyRepo := NewDynamoRepo(dyDB)
//
// 	// err := dyRepo.AddTask(context.TODO(), NewTodo("makan nasi", false))
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }
//
// 	todos, err := dyRepo.GetTasks(context.TODO())
// 	if err != nil {
// 		fmt.Println(err)
// 	}
//
// 	out, _ := json.MarshalIndent(todos, "  ", "    ")
// 	fmt.Println(string(out))
//
// 	id := todos[1].ID
//
// 	err = dyRepo.RemoveTask(context.TODO(), id)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
//
// 	todos, err = dyRepo.GetTasks(context.TODO())
// 	if err != nil {
// 		fmt.Println(err)
// 	}
//
// 	out, _ = json.MarshalIndent(todos, "  ", "    ")
// 	fmt.Println(string(out))
//
// }
