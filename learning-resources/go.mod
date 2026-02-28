module github.com/learnbot/learning-resources

go 1.24.0

require (
	github.com/google/uuid v1.6.0
	github.com/learnbot/database v0.0.0
	github.com/lib/pq v1.11.2
)

replace github.com/learnbot/database => ../database
