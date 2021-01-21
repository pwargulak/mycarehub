package resources

import (
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"

	"gitlab.slade360emr.com/go/base"
)

// BusinessPartnerEdge is used to serialize GraphQL Relay edges for organization
type BusinessPartnerEdge struct {
	Cursor *string                 `json:"cursor"`
	Node   *domain.BusinessPartner `json:"node"`
}

// BusinessPartnerConnection is used to serialize GraphQL Relay connections for organizations
type BusinessPartnerConnection struct {
	Edges    []*BusinessPartnerEdge `json:"edges"`
	PageInfo *base.PageInfo         `json:"pageInfo"`
}

// AccountRecoveryPhonesResponse  payload sent back to the frontend when recovery an account
type AccountRecoveryPhonesResponse struct {
	MaskedPhoneNumbers   []string `json:"maskedPhoneNumbers"`
	UnMaskedPhoneNumbers []string `json:"unmaskedPhoneNumbers"`
}

// OKResp is used to return OK responses in inter-service calls
type OKResp struct {
	Status   string      `json:"status,omitempty"`
	Response interface{} `json:"response,omitempty"`
}

// NewOKResp a shortcut to create an instance of OKResp
func NewOKResp(rawResponse interface{}) *OKResp {
	return &OKResp{
		Status:   "OK",
		Response: rawResponse,
	}
}

// CreatedUserResponse is used to return a created user
type CreatedUserResponse struct {
	UID         string `json:"uid,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	PhotoURL    string `json:"photo_url,omitempty"`
	ProviderID  string `json:"provider_id,omitempty"`
}

// SupplierLogin is the response returned after the user has successfully login to edi
type SupplierLogin struct {
	Branches *BranchConnection `json:"branches,omitempty"`
	Supplier *base.Supplier    `json:"supplier,omitempty"`
}
