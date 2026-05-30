package handler

import (
	"github.com/gsystes/backend/internal/communication/websocket"
	domainRepo "github.com/gsystes/backend/internal/domain/repository"
)

type EventBroadcaster struct {
	hub      *websocket.Hub
	userRepo domainRepo.UserRepository
	roleRepo domainRepo.RoleRepository
	logRepo  domainRepo.OperationLogRepository
}

func NewEventBroadcaster(
	hub *websocket.Hub,
	userRepo domainRepo.UserRepository,
	roleRepo domainRepo.RoleRepository,
	logRepo domainRepo.OperationLogRepository,
) *EventBroadcaster {
	return &EventBroadcaster{
		hub:      hub,
		userRepo: userRepo,
		roleRepo: roleRepo,
		logRepo:  logRepo,
	}
}

func (b *EventBroadcaster) BroadcastStats() {
	if b.hub == nil || b.hub.ClientCount() == 0 {
		return
	}
	userCount, _ := b.userRepo.Count()
	roleCount, _ := b.roleRepo.Count()
	todayLogCount, _ := b.logRepo.CountToday()

	b.hub.BroadcastStatUpdate(&websocket.StatUpdatePayload{
		UserCount:     userCount,
		RoleCount:     roleCount,
		TodayLogCount: todayLogCount,
	})
}

func (b *EventBroadcaster) SendNotification(username, title, message string) {
	if b.hub == nil {
		return
	}
	b.hub.BroadcastNotification(&websocket.NotificationPayload{
		Username: username,
		Title:    title,
		Message:  message,
	})
}
