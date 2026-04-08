package main

import "time"

type Queue struct {
  ID          string    `json:"id"`
  Name        string    `json:"name"`
  Description string    `json:"description"`
  StartTime   string    `json:"startTime"`   // "12:00"
  MaxParticipants int   `json:"maxParticipants"`
  HasPriority bool      `json:"hasPriority"`
  PriorityCount int     `json:"priorityCount"`
  InitialPriority int   `json:"initialPriority"`
  AnonymousChat bool    `json:"anonymousChat"`
  SystemNotifications bool `json:"systemNotifications"`
  SwapPositions bool    `json:"swapPositions"`
  Admins      []string  `json:"admins"`
  CreatedAt   time.Time `json:"createdAt"`
  Participants []string `json:"-"` // список ID участников (пока не используется)
}

type CreateQueueRequest struct {
  Name        string   `json:"queueName"`
  Description string   `json:"description"`
  StartTime   string   `json:"startTime"`
  MaxParticipants int  `json:"maxParticipants"`
  HasPriority bool     `json:"hasPriority"`
  PriorityCount int    `json:"priorityCount"`
  InitialPriority int  `json:"initialPriority"`
  AnonymousChat bool   `json:"anonymousChat"`
  SystemNotifications bool `json:"systemNotifications"`
  SwapPositions bool   `json:"swapPositions"`
  Admins      []string `json:"admins"`
}

type QueueResponse struct {
  ID   string `json:"id"`
  Name string `json:"name"`
  Link string `json:"link"`
}