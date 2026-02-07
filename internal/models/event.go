package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventStatus string

const (
	EventStatusUpcoming  EventStatus = "upcoming"
	EventStatusOngoing   EventStatus = "ongoing"
	EventStatusPast      EventStatus = "past"
	EventStatusCancelled EventStatus = "cancelled"
)

type Event struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title         string             `bson:"title" json:"title"`
	Description   string             `bson:"description" json:"description"`
	Content       string             `bson:"content" json:"content"`
	OrganizerID   primitive.ObjectID `bson:"organizer_id" json:"organizer_id"`
	OrganizerName string             `bson:"organizer_name" json:"organizer_name"`
	Location      string             `bson:"location" json:"location"`
	StartDate     time.Time          `bson:"start_date" json:"start_date"`
	EndDate       time.Time          `bson:"end_date" json:"end_date"`
	Category      PostCategory       `bson:"category" json:"category"`
	Status        EventStatus        `bson:"status" json:"status"`
	MaxAttendees  int                `bson:"max_attendees,omitempty" json:"max_attendees,omitempty"`
	AttendeeCount int                `bson:"attendee_count" json:"attendee_count"`
	Media         []MediaItem        `bson:"media,omitempty" json:"media,omitempty"`
	MediaCount    int                `bson:"media_count" json:"media_count"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

func NewEvent(title, description, content, location string, organizerID primitive.ObjectID, organizerName string, startDate, endDate time.Time, category PostCategory) *Event {
	now := time.Now()

	var status EventStatus
	if startDate.After(now) {
		status = EventStatusUpcoming
	} else if endDate.After(now) {
		status = EventStatusOngoing
	} else {
		status = EventStatusPast
	}

	return &Event{
		ID:            primitive.NewObjectID(),
		Title:         title,
		Description:   description,
		Content:       content,
		OrganizerID:   organizerID,
		OrganizerName: organizerName,
		Location:      location,
		StartDate:     startDate,
		EndDate:       endDate,
		Category:      category,
		Status:        status,
		AttendeeCount: 0,
		MediaCount:    0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func (e *Event) AddMedia(url, mediaType, caption string) {
	media := MediaItem{
		ID:       primitive.NewObjectID(),
		URL:      url,
		Type:     mediaType,
		Caption:  caption,
		Position: len(e.Media),
	}
	e.Media = append(e.Media, media)
	e.MediaCount = len(e.Media)
}

func (e *Event) UpdateStatus() {
	now := time.Now()
	if e.StartDate.After(now) {
		e.Status = EventStatusUpcoming
	} else if e.EndDate.After(now) {
		e.Status = EventStatusOngoing
	} else {
		e.Status = EventStatusPast
	}
}
