package enums

import (
	"fmt"
	"io"
	"strconv"
)

// ServiceRequestType is a list of all the service request types.
type ServiceRequestType string

const (
	//ServiceRequestTypeRedFlag represents a health diary entry
	ServiceRequestTypeRedFlag ServiceRequestType = "RED_FLAG"
	// ServiceRequestTypePinReset represents the client's reset pin service request
	ServiceRequestTypePinReset ServiceRequestType = "PIN_RESET"
	// ServiceRequestTypeStaffPinReset represents the reset pin service request
	ServiceRequestTypeStaffPinReset ServiceRequestType = "STAFF_PIN_RESET"
	// ServiceRequestTypeProfileUpdate represents the profile update service request
	ServiceRequestTypeProfileUpdate ServiceRequestType = "PROFILE_UPDATE"
	// ServiceRequestTypeHomePageHealthDiary represents the homepage healthdiary service request
	ServiceRequestTypeHomePageHealthDiary ServiceRequestType = "HOME_PAGE_HEALTH_DIARY_ENTRY"
)

// AllServiceRequestType is a set of a  valid and known service request types.
var AllServiceRequestType = []ServiceRequestType{
	ServiceRequestTypeRedFlag,
	ServiceRequestTypePinReset,
	ServiceRequestTypeStaffPinReset,
	ServiceRequestTypeProfileUpdate,
	ServiceRequestTypeHomePageHealthDiary,
}

// IsValid returns true if a request type is valid
func (m ServiceRequestType) IsValid() bool {
	switch m {
	case ServiceRequestTypeRedFlag,
		ServiceRequestTypePinReset,
		ServiceRequestTypeStaffPinReset,
		ServiceRequestTypeProfileUpdate,
		ServiceRequestTypeHomePageHealthDiary:
		return true
	}
	return false
}

func (m ServiceRequestType) String() string {
	return string(m)
}

// UnmarshalGQL converts the supplied value to a request type.
func (m *ServiceRequestType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = ServiceRequestType(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid SortDataType", str)
	}
	return nil
}

// MarshalGQL writes the sort type to the supplied
func (m ServiceRequestType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}
