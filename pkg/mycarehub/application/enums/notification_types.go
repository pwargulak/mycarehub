package enums

import (
	"fmt"
	"io"
	"strconv"
)

// NotificationType represents a type of notification
type NotificationType string

const (
	// NotificationTypeAppointment represents notifications from appointments
	NotificationTypeAppointment NotificationType = "APPOINTMENT"

	// NotificationTypeServiceRequest represents notifications from service requests
	NotificationTypeServiceRequest NotificationType = "SERVICE_REQUEST"

	// NotificationTypeCommunities represents notifications from communities
	NotificationTypeCommunities NotificationType = "COMMUNITIES"
)

// AllNotificationTypes represents a slice of all possible `NotificationType` values
var AllNotificationTypes = []NotificationType{
	NotificationTypeAppointment,
	NotificationTypeServiceRequest,
	NotificationTypeCommunities,
}

// IsValid returns true if a notification type is valid
func (n NotificationType) IsValid() bool {
	switch n {
	case
		NotificationTypeAppointment,
		NotificationTypeServiceRequest,
		NotificationTypeCommunities:
		return true
	}
	return false
}

// String returns a string representation of the enum
func (n NotificationType) String() string {
	return string(n)
}

// UnmarshalGQL converts the supplied value to a metric type.
func (n *NotificationType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*n = NotificationType(str)
	if !n.IsValid() {
		return fmt.Errorf("%s is not a valid NotificationType", str)
	}
	return nil
}

// MarshalGQL writes the metric type to the supplied writer
func (n NotificationType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(n.String()))
}