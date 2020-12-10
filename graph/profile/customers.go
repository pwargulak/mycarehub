package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"gitlab.slade360emr.com/go/base"
)

const (
	customerAPIPath        = "/api/business_partners/customers/"
	active                 = true
	country                = "KEN" // Anticipate worldwide expansion
	isCustomer             = true
	customerType           = PartnerTypesConsumer
	customerCollectionName = "customers"
)

// SaveCustomerToFireStore persists customer data to firestore
func (s Service) SaveCustomerToFireStore(customer Customer) error {
	ctx := context.Background()
	_, _, err := s.firestoreClient.Collection(s.GetCustomerCollectionName()).Add(ctx, customer)
	return err
}

// GetCustomerCollectionName creates a suffixed customer collection name
func (s Service) GetCustomerCollectionName() string {
	suffixed := base.SuffixCollection(customerCollectionName)
	return suffixed
}

// AddCustomer makes a call to our own ERP and creates a customer account for the consumer users
// that is used for transacting on Be.Well
func (s Service) AddCustomer(ctx context.Context, uid *string, name string) (*Customer, error) {
	s.checkPreconditions()

	userUID, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
	}

	profile, err := s.ParseUserProfileFromContextOrUID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("unable to read user profile: %w", err)
	}

	collection := s.firestoreClient.Collection(s.GetCustomerCollectionName())
	query := collection.Where("userprofile.verifiedIdentifiers", "array-contains", userUID)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	if len(docs) > 1 {
		if base.IsDebug() {
			log.Printf("uid %s has more than one customer records (it has %d)", userUID, len(docs))
		}
	}

	if len(docs) == 0 {
		currency, err := base.FetchDefaultCurrency(s.erpClient)
		if err != nil {
			return nil, fmt.Errorf("unable to fetch orgs default currency: %v", err)
		}
		payload := map[string]interface{}{
			"active":        active,
			"partner_name":  name,
			"country":       country,
			"currency":      *currency.ID,
			"is_customer":   isCustomer,
			"customer_type": customerType,
		}

		content, marshalErr := json.Marshal(payload)
		if marshalErr != nil {
			return nil, fmt.Errorf("unable to marshal to JSON: %v", marshalErr)
		}
		newCustomer := Customer{
			UserProfile: *profile,
		}

		if err := base.ReadRequestToTarget(s.erpClient, "POST", customerAPIPath, "", content, &newCustomer); err != nil {
			return nil, fmt.Errorf("unable to make request to the ERP: %v", err)
		}

		if err := s.SaveCustomerToFireStore(newCustomer); err != nil {
			return nil, fmt.Errorf("unable to add customer to firestore: %v", err)
		}

		profile.HasCustomerAccount = true
		profileDsnap, err := s.RetrieveUserProfileFirebaseDocSnapshot(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve firebase user profile: %v", err)
		}

		if err = base.UpdateRecordOnFirestore(
			s.firestoreClient, s.GetUserProfileCollectionName(), profileDsnap.Ref.ID, profile,
		); err != nil {
			return nil, fmt.Errorf("unable to update user profile: %v", err)
		}

		return &newCustomer, nil
	}
	dsnap := docs[0]
	customer := &Customer{}
	err = dsnap.DataTo(customer)
	if err != nil {
		return nil, fmt.Errorf("unable to read customer: %w", err)
	}

	return customer, nil
}

// AddCustomerKYC persists information that is relevant to knowing our customers
func (s Service) AddCustomerKYC(ctx context.Context, input CustomerKYCInput) (*CustomerKYC, error) {
	s.checkPreconditions()

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
	}

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(
		ctx, uid, s.GetCustomerCollectionName(), "userprofile.verifiedIdentifiers")
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve customer: %v", err)
	}

	if dsnap == nil {
		return nil, fmt.Errorf("customer not found")
	}

	customer := &Customer{}
	err = dsnap.DataTo(customer)
	if err != nil {
		return nil, fmt.Errorf("unable to read customer data: %v", err)
	}

	customer.CustomerKYC.KRAPin = input.KRAPin
	customer.CustomerKYC.Occupation = input.Occupation
	customer.CustomerKYC.IDNumber = input.IDNumber
	customer.CustomerKYC.Address = input.Address
	customer.CustomerKYC.City = input.City

	beneficiaries := customer.CustomerKYC.Beneficiary

	var foundBeneficiariesNames []string
	for _, b := range beneficiaries {
		foundBeneficiariesNames = append(foundBeneficiariesNames, b.Name)
	}

	for _, beneficiary := range input.Beneficiary {
		beneficiaryData := &Beneficiary{
			Name:         beneficiary.Name,
			Msisdns:      beneficiary.Msisdns,
			Emails:       beneficiary.Emails,
			Relationship: beneficiary.Relationship,
			DateOfBirth:  beneficiary.DateOfBirth,
		}
		if !base.StringSliceContains(foundBeneficiariesNames, beneficiaryData.Name) {
			beneficiaries = append(beneficiaries, beneficiaryData)
		}
	}

	customer.CustomerKYC.Beneficiary = beneficiaries

	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetCustomerCollectionName(), dsnap.Ref.ID, customer,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to update customer with customer KYC info: %v", err)
	}

	customerKYC := customer.CustomerKYC
	return &customerKYC, nil
}

// UpdateCustomer updates a customerKYC information in firestore
func (s Service) UpdateCustomer(ctx context.Context, input CustomerKYCInput) (*Customer, error) {
	s.checkPreconditions()

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
	}

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(
		ctx, uid, s.GetCustomerCollectionName(), "userprofile.verifiedIdentifiers")
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve customer: %v", err)
	}

	customer := &Customer{}
	err = dsnap.DataTo(customer)
	if err != nil {
		return nil, fmt.Errorf("unable to read customer: %w", err)
	}

	if input.KRAPin != "" {
		customer.CustomerKYC.KRAPin = input.KRAPin
	}

	if input.Occupation != "" {
		customer.CustomerKYC.Occupation = input.Occupation
	}

	if input.IDNumber != "" {
		customer.CustomerKYC.IDNumber = input.IDNumber
	}

	if input.City != "" {
		customer.CustomerKYC.City = input.City
	}

	if input.Address != "" {
		customer.CustomerKYC.Address = input.Address
	}

	beneficiaries := customer.CustomerKYC.Beneficiary
	if input.Beneficiary != nil {
		for _, beneficiary := range input.Beneficiary {
			beneficiaryData := &Beneficiary{
				Name:         beneficiary.Name,
				Msisdns:      beneficiary.Msisdns,
				Emails:       beneficiary.Emails,
				Relationship: beneficiary.Relationship,
				DateOfBirth:  beneficiary.DateOfBirth,
			}
			beneficiaries = append(beneficiaries, beneficiaryData)
		}
		customer.CustomerKYC.Beneficiary = beneficiaries
	}

	err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetCustomerCollectionName(), dsnap.Ref.ID, customer,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to update customer with customer KYC info: %v", err)
	}

	return customer, nil
}

// FindCustomer fetches a customer by their UID
func (s Service) FindCustomer(ctx context.Context, uid string) (*Customer, error) {
	s.checkPreconditions()

	dsnap, err := s.RetrieveFireStoreSnapshotByUID(
		ctx, uid, s.GetCustomerCollectionName(), "userprofile.verifiedIdentifiers")
	if err != nil {
		return nil, fmt.Errorf("unable to retreive doc snapshot by uid: %v", err)
	}

	if dsnap == nil {
		if dsnap == nil {
			return nil, fmt.Errorf("a user with the UID %s does not have a customer's account", uid)
		}
	}

	customer := &Customer{}
	err = dsnap.DataTo(customer)
	if err != nil {
		return nil, fmt.Errorf("unable to read customer: %v", err)
	}
	return customer, nil
}

// SuspendCustomer flips the active boolean on the erp partner from true to false
// consequently logically deleting the account
func (s Service) SuspendCustomer(ctx context.Context, uid string) (bool, error) {
	s.checkPreconditions()

	err := s.DeleteUser(ctx, uid)
	if err != nil {
		return false, fmt.Errorf("error deleting user: %v", err)
	}

	collection := s.firestoreClient.Collection(s.GetCustomerCollectionName())
	query := collection.Where("userprofile.verifiedIdentifiers", "array-contains", uid)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return false, err
	}
	if len(docs) == 0 {
		return false, nil
	}

	dsnap := docs[0]
	customer := &Customer{}
	err = dsnap.DataTo(customer)
	if err != nil {
		return false, fmt.Errorf("unable to read customer: %w", err)
	}

	payload := map[string]interface{}{
		"active": false,
	}

	content, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return false, fmt.Errorf("unable to marshal to JSON: %v", marshalErr)
	}

	customerPath := fmt.Sprintf("%s%s/", customerAPIPath, customer.CustomerID)
	if err := base.ReadRequestToTarget(s.erpClient, "PATCH", customerPath, "", content, &customer); err != nil {
		return false, fmt.Errorf("unable to make request to the ERP: %v", err)
	}

	if err = base.UpdateRecordOnFirestore(
		s.firestoreClient, s.GetCustomerCollectionName(), dsnap.Ref.ID, customer,
	); err != nil {
		return false, fmt.Errorf("unable to update customer: %v", err)
	}

	return true, nil
}
