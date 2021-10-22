package dto

import (
	"net/url"
	"time"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/enums"
	"gorm.io/datatypes"
)

// FacilityInput describes the facility input
type FacilityInput struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Active      bool   `json:"active"`
	County      string `json:"county"`
	Description string `json:"description"`
}

// FacilityFilterInput is used to supply filter parameters for healthcare facility filter inputs
type FacilityFilterInput struct {
	Search  *string `json:"search"`
	Name    *string `json:"name"`
	MFLCode *string `json:"code"`
}

// ToURLValues transforms the filter input to `url.Values`
func (i *FacilityFilterInput) ToURLValues() (values url.Values) {
	vals := url.Values{}
	if i.Search != nil {
		vals.Add("search", *i.Search)
	}
	if i.Name != nil {
		vals.Add("name", *i.Name)
	}
	if i.MFLCode != nil {
		vals.Add("code", *i.MFLCode)
	}
	return vals
}

// FacilitySortInput is used to supply sort input for healthcare facility list queries
type FacilitySortInput struct {
	Name    *enumutils.SortOrder `json:"name"`
	MFLCode *enumutils.SortOrder `json:"code"`
}

// ToURLValues transforms the filter input to `url.Values`
func (i *FacilitySortInput) ToURLValues() (values url.Values) {
	vals := url.Values{}
	if i.Name != nil {
		if *i.Name == enumutils.SortOrderAsc {
			vals.Add("order_by", "name")
		} else {
			vals.Add("order_by", "-name")
		}
	}
	if i.MFLCode != nil {
		if *i.Name == enumutils.SortOrderAsc {
			vals.Add("code", "number")
		} else {
			vals.Add("code", "-number")
		}
	}
	return vals
}

// MetricInput reprents the metrics data structure input
type MetricInput struct {

	// TODO Metric types should be a controlled list i.e enum
	Type enums.MetricType `json:"metric_type"`

	// this will vary by context
	// should not identify the user (there's a UID field)
	// focus on the actual event
	Payload datatypes.JSON `gorm:"column:payload"`

	Timestamp time.Time `json:"time"`

	// a user identifier, can be hashed for anonymity
	// with a predictable one way hash
	UID string `json:"uid"`
}

// PINInput represents the PIN input data structure
type PINInput struct {
	PIN          string          `json:"pin"`
	ConfirmedPin string          `json:"confirmed_pin"`
	Flavour      feedlib.Flavour `json:"flavour"`
}

// LoginInput represents the Login input data structure
type LoginInput struct {
	UserID  string          `json:"user_id"`
	PIN     string          `json:"pin"`
	Flavour feedlib.Flavour `json:"flavour"`
}

// StaffProfileInput contains input required to register a staff
type StaffProfileInput struct {
	StaffNumber string

	// Facilities []*domain.Facility // TODO: needs at least one

	// A UI switcher optionally toggles the default
	// TODO: the list of facilities to switch between is strictly those that the user is assigned to
	DefaultFacilityID *string // TODO: required, FK to facility

	// // there is nothing special about super-admin; just the set of roles they have
	// Roles []domain.RoleType `json:"roles"` // TODO: roles are an enum (controlled list), known to both FE and BE

	// Addresses []*domain.UserAddress
}

// ClientProfileInput is used to supply the client profile input
type ClientProfileInput struct {
	ClientType enums.ClientType `json:"clientType"`
}

// UserInput is used to supply user input for registration
type UserInput struct {
	Username string // @handle, also globally unique; nickname

	DisplayName string // user's preferred display name

	// TODO Consider making the names optional in DB; validation in frontends
	FirstName  string // given name
	MiddleName string
	LastName   string

	UserType enums.UsersType // TODO enum; e.g client, health care worker

	Gender enumutils.Gender // TODO enum; genders; keep it simple

	Contacts []*ContactInput // TODO: validate, ensure

	// // for the preferred language list, order matters
	Languages []enumutils.Language // TODO: turn this into a slice of enums, start small (en, sw)
	Flavour   feedlib.Flavour
}

// ContactInput contains input required to register a user
type ContactInput struct {
	Type    enums.ContactType
	Contact string //TODO Validate: phones are E164, emails are valid
	Active  bool
	//   a user may opt not to be contacted via this contact
	//   e.g if it's a shared phone owned by a teenager
	OptedIn bool
}
