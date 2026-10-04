package model

import "github.com/MohammadrezaNadirkhanloo/internal/domain/filter"

var TodoSpec = filter.NewSpec("id desc",
	filter.Number("id", "id"),

	filter.Text("title", "title").Search(),
	filter.Text("category", "category").Search(),

	filter.Bool("done", "done"),

	filter.Select("priority", "priority",
		PriorityLow, PriorityMedium, PriorityHigh),

	filter.DateTime("dueDate", "due_date"),

	filter.Tags("tags", "tags").ArrayColumn().NoSort(),

	filter.DateTime("createdAt", "created_at"),

	filter.Number("userId", "user_id").NoFilter().NoSort(),
)

var UserSpec = filter.NewSpec("id desc",
	filter.Number("id", "id"),
	filter.Text("username", "username").Search(),
	filter.DateTime("createdAt", "created_at"),
)
