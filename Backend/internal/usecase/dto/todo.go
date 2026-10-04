package dto

import (
	"strings"
	"time"

	"github.com/MohammadrezaNadirkhanloo/internal/domain/model"
	"github.com/lib/pq"
)

type CreateTodoInput struct {
	Title    string     `json:"title"    binding:"required,min=1,max=200"`
	DueDate  *time.Time `json:"dueDate"`
	Category *string    `json:"category" binding:"omitempty,max=50"`
	Priority string     `json:"priority" binding:"omitempty,oneof=low medium high"`
	Tags     []string   `json:"tags"     binding:"omitempty,max=20,dive,min=1,max=30"`
}

type UpdateTodoInput struct {
	Title    *string    `json:"title"    binding:"omitempty,min=1,max=200"`
	Done     *bool      `json:"done"`
	DueDate  *time.Time `json:"dueDate"`
	Category *string    `json:"category" binding:"omitempty,max=50"`
	Priority *string    `json:"priority" binding:"omitempty,oneof=low medium high"`
	Tags     *[]string  `json:"tags"     binding:"omitempty,max=20,dive,min=1,max=30"`
}

type TodoOutput struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	Done      bool       `json:"done"`
	DueDate   *time.Time `json:"dueDate"`
	Category  *string    `json:"category"`
	Priority  string     `json:"priority"`
	Tags      []string   `json:"tags"`
	CreatedAt time.Time  `json:"createdAt"`
}

func ToTodoModel(in CreateTodoInput) model.Todo {
	priority := in.Priority
	if priority == "" {
		priority = model.PriorityMedium
	}

	tags := pq.StringArray{}
	if len(in.Tags) > 0 {
		tags = pq.StringArray(in.Tags)
	}

	return model.Todo{
		Title:    in.Title,
		DueDate:  in.DueDate,
		Category: emptyToNil(in.Category),
		Priority: priority,
		Tags:     tags,
	}
}

func emptyToNil(s *string) *string {
	if s == nil {
		return nil
	}
	t := strings.TrimSpace(*s)
	if t == "" {
		return nil
	}
	return &t
}

func TodoChanges(in UpdateTodoInput) map[string]any {
	changes := make(map[string]any, 6)

	if in.Title != nil {
		changes["title"] = *in.Title
	}
	if in.Done != nil {
		changes["done"] = *in.Done
	}
	if in.DueDate != nil {
		changes["due_date"] = *in.DueDate
	}
	if in.Category != nil {
		if v := emptyToNil(in.Category); v == nil {
			changes["category"] = nil
		} else {
			changes["category"] = *v
		}
	}
	if in.Priority != nil {
		changes["priority"] = *in.Priority
	}
	if in.Tags != nil {
		changes["tags"] = pq.StringArray(*in.Tags)
	}
	return changes
}

func ToTodoOutput(m model.Todo) TodoOutput {
	tags := []string(m.Tags)
	if tags == nil {
		tags = []string{}
	}

	return TodoOutput{
		ID:        m.ID,
		Title:     m.Title,
		Done:      m.Done,
		DueDate:   m.DueDate,
		Category:  m.Category,
		Priority:  m.Priority,
		Tags:      tags,
		CreatedAt: m.CreatedAt,
	}
}
