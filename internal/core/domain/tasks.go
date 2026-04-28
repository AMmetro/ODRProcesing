package domain

type Task struct {
	ID           int
	Title        string
	Description  string
	Completed    bool
	AuthorUserId int
	Version      int
}

// func NewTask(
// 	id int,
// 	title string,
// 	description string,
// 	completed bool,
// 	authorUserId int,
// 	version int,
// ) Task {
// 	return Task{
// 		ID:           id,
// 		Title:        title,
// 		Description:  description,
// 		Completed:    completed,
// 		AuthorUserId: authorUserId,
// 		Version:      version,
// 	}
// }
