package job

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
)

const TypeSendWelcomeEmail = "email:welcome"

type SendWelcomeEmailPayload struct {
	UserID uint `json:"user_id"`
}

// NewSendWelcomeEmailTask は dispatch 用のタスクを作る
func NewSendWelcomeEmailTask(userID uint) (*asynq.Task, error) {
	b, err := json.Marshal(SendWelcomeEmailPayload{UserID: userID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeSendWelcomeEmail, b), nil
}

// HandleSendWelcomeEmail は Job の handle() 相当
func HandleSendWelcomeEmail(ctx context.Context, t *asynq.Task) error {
	var p SendWelcomeEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	log.Printf("send welcome email to user %d", p.UserID)
	return nil
}
