package dto

type DashboardStatsResponse struct {
	UserCount     int64 `json:"user_count"`
	RoleCount     int64 `json:"role_count"`
	TodayLogCount int64 `json:"today_log_count"`
}
