package events

import (
	eventsv1http "github.com/edgequota/edgequota-go/gen/http/events/v1"
)

// Accepted returns a PublishEventsResponse indicating all events were accepted.
func Accepted(count int) eventsv1http.PublishEventsResponse {
	return eventsv1http.PublishEventsResponse{Accepted: int64(count)}
}
