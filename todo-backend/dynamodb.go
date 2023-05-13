package main

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/rifqoi/aws-project/todo-backend/common"
)

type TodoRepository interface {
	GetTasks(ctx context.Context) ([]Todo, error)
	AddTask(ctx context.Context, task Todo) error
	UpdateTask(ctx context.Context, updatedTask Todo) (*Todo, error)
	RemoveTask(id int) Todo
}

type DynamoRepo struct {
	db *dynamodb.Client
}

func NewDynamoRepo(db *dynamodb.Client) *DynamoRepo {
	if db == nil {
		panic("db is nil")
	}

	return &DynamoRepo{db}
}
func (dr *DynamoRepo) GetTasks(ctx context.Context) ([]Todo, error) {
	out, err := dr.db.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String("todo"),
	})

	if err != nil {
		return nil, common.WrapErrorf(err, common.ErrorTypeUnknown, "failed to get items")
	}

	todos := []Todo{}
	for _, item := range out.Items {
		var todo Todo
		err := attributevalue.UnmarshalMap(item, &todo)
		if err != nil {
			return nil, common.WrapErrorf(err, common.ErrorTypeServerError, "failed unmarshal value")
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func (dr *DynamoRepo) AddTask(ctx context.Context, todo Todo) error {
	_, err := dr.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("todo"),
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
			return common.WrapErrorf(err, common.ErrorTypeServerError, "table not found")

		}
		return common.WrapErrorf(err, common.ErrorTypeUnknown, "failed to inserting item to dynamodb")
	}

	return nil
}

func (dr *DynamoRepo) UpdateTask(ctx context.Context, updatedTodo Todo) (*Todo, error) {

	return nil, nil
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
// 	out, _ := json.MarshalIndent(todos, "  ", "    ")
// 	fmt.Println(string(out))
//
// }
