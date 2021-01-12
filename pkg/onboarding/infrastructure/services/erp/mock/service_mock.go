package mock

import (
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

// FakeServiceERP is an `ERP` service mock .
type FakeServiceERP struct {
	FetchERPClientFn    func() *base.ServerClient
	CreateERPSupplierFn func(method string, path string, payload map[string]interface{}, partner domain.PartnerType) error
}

// FetchERPClient ...
func (f *FakeServiceERP) FetchERPClient() *base.ServerClient {
	return f.FetchERPClientFn()
}

// CreateERPSupplier ...
func (f *FakeServiceERP) CreateERPSupplier(method string, path string, payload map[string]interface{}, partner domain.PartnerType) error {
	return f.CreateERPSupplierFn(method, path, payload, partner)
}
