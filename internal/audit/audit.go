package audit

import (
	"wastegen/internal/store"
)

type Recorder struct {
	auditStore *store.AuditStore
}

func NewRecorder(auditStore *store.AuditStore) *Recorder {
	return &Recorder{auditStore: auditStore}
}

func (r *Recorder) Record(event Event) error {
	event = NewEvent(event.FurnaceID, event.EventType, event.Message)
	return r.auditStore.Append(store.AuditRecord{
		ID:        event.ID,
		FurnaceID: event.FurnaceID,
		EventType: event.EventType,
		Message:   event.Message,
		At:        event.At,
	})
}

func (r *Recorder) Events() ([]Event, error) {
	records, err := r.auditStore.List()
	if err != nil {
		return nil, err
	}
	events := make([]Event, 0, len(records))
	for _, record := range records {
		events = append(events, Event{
			ID:        record.ID,
			FurnaceID: record.FurnaceID,
			EventType: record.EventType,
			Message:   record.Message,
			At:        record.At,
		})
	}
	return events, nil
}

func (r *Recorder) ForFurnace(furnaceID string) ([]Event, error) {
	records, err := r.auditStore.ForFurnace(furnaceID)
	if err != nil {
		return nil, err
	}
	events := make([]Event, 0, len(records))
	for _, record := range records {
		events = append(events, Event{
			ID:        record.ID,
			FurnaceID: record.FurnaceID,
			EventType: record.EventType,
			Message:   record.Message,
			At:        record.At,
		})
	}
	return events, nil
}

func (r *Recorder) Recent(limit int) ([]Event, error) {
	records, err := r.auditStore.Recent(limit)
	if err != nil {
		return nil, err
	}
	events := make([]Event, 0, len(records))
	for _, record := range records {
		events = append(events, Event{
			ID:        record.ID,
			FurnaceID: record.FurnaceID,
			EventType: record.EventType,
			Message:   record.Message,
			At:        record.At,
		})
	}
	return events, nil
}

func (r *Recorder) Count() int {
	records, err := r.auditStore.List()
	if err != nil {
		return 0
	}
	return len(records)
}
