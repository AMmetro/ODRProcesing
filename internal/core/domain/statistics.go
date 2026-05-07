package domain

type StatisticsSummary struct {
	TotalUsers     int `json:"total_users"`
	TotalTasks     int `json:"total_tasks"`
	CompletedTasks int `json:"completed_tasks"`
	PendingTasks   int `json:"pending_tasks"`
	// TasksPerUser   float64 `json:"tasks_per_user"`
}

func NewStatisticsSummary(
	totalUsers int,
	totalTasks int,
	completedTasks int,
	pendingTasks int,
) StatisticsSummary {
	// tasksPerUser := 0.0
	// if totalUsers > 0 {
	// 	tasksPerUser = float64(totalTasks) / float64(totalUsers)
	// }

	return StatisticsSummary{
		TotalUsers:     totalUsers,
		TotalTasks:     totalTasks,
		CompletedTasks: completedTasks,
		PendingTasks:   pendingTasks,
		// TasksPerUser:   tasksPerUser,
	}
}
