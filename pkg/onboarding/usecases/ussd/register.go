package ussd

import (
	"context"
	"strconv"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

const (
	// InitialState ...
	InitialState = 0
	// GetFirstNameState ...
	GetFirstNameState = 1
	// GetLastNameState ...
	GetLastNameState = 2
	// GetDOBState ...
	GetDOBState = 3
	// GetPINState ...
	GetPINState = 4
	// SaveRecord ...
	SaveRecord = 5
	// WelcomeMenu ...
	WelcomeMenu = 6
	// RegisterInput ...
	RegisterInput = "1"
	// BuyCoverInput ...
	BuyCoverInput = "2"
	// RegWantToBuyCoverInput ...
	RegWantToBuyCoverInput = "1"
	//RegOptOutInput ...
	RegOptOutInput = "2"
	//RegChangePINInput ...
	RegChangePINInput = "3"

	//date layout
	// layoutISO = "01-02-2006"
)

var userFirstName string
var userLastName string
var date string

// HandleUserRegistration ...
func (u *Impl) HandleUserRegistration(ctx context.Context, session *domain.USSDLeadDetails, userResponse string) string {
	if userResponse == "" && session.Level == InitialState {
		// contact := CRMDomain.CRMContact{
		// 	Properties: CRMDomain.ContactProperties{
		// 		Phone:                 session.PhoneNumber,
		// 		FirstChannelOfContact: CRMDomain.ChannelOfContactUssd,
		// 	},
		// }
		// err := u.CreateCRMContact(ctx, contact)
		// if err != nil {
		// 	return "END something wrong happened"
		// }
		resp := "CON Welcome to Be.Well.\r\n"
		resp += "1. Register\r\n"
		resp += "2. I want a cover\r\n"
		return resp
	}

	if userResponse == RegisterInput && session.Level == InitialState {
		err := u.UpdateSessionLevel(ctx, GetFirstNameState, session.SessionID)
		if err != nil {
			return "END something wrong happened"
		}
		resp := "CON Please enter your firstname(e.g.\r\n"
		resp += "John).\r\n"
		return resp
	}

	if userResponse == BuyCoverInput && session.Level == InitialState {
		resp := "END We have recorded your request.\r\n"
		resp += "and one of the representatives will.\r\n"
		resp += "reach out to you. Thank you.\r\n"
		return resp

	}
	if session.Level == GetFirstNameState {
		err := utils.ValidateUSSDInput(userResponse)
		if err != nil {
			return "CON Invalid name. Please enter a valid name (e.g John)"
		}

		isLetter := utils.IsLetter(userResponse)
		if !isLetter {
			return "CON Invalid name. Please enter a valid name (e.g John)"
		}
		userFirstName = userResponse

		//Update CRM Firstname
		// var CRMContactProperties CRMDomain.ContactProperties
		// if userFirstName != "" {
		// 	CRMContactProperties.FirstName = userFirstName
		// }
		// err = u.UpdateCRMContactDetails(CRMContactProperties, session, ctx)
		// if err != nil {
		// 	return err.Error()
		// }

		err = u.UpdateSessionLevel(ctx, GetLastNameState, session.SessionID)
		if err != nil {
			return "END something wrong happened"
		}
		resp := "CON Please enter your lastname(e.g.\r\n"
		resp += "Doe).\r\n"
		return resp

	}

	if session.Level == GetLastNameState {
		err := utils.ValidateUSSDInput(userResponse)
		if err != nil {
			return "CON Invalid name. Please enter a valid name (e.g John)"
		}

		isLetter := utils.IsLetter(userResponse)
		if !isLetter {
			return "CON Invalid name. Please enter a valid name (e.g John)"
		}

		userLastName = userResponse

		//Update CRM Lastname
		// var CRMContactProperties CRMDomain.ContactProperties
		// if userLastName != "" {
		// 	CRMContactProperties.FirstName = userLastName
		// }

		// err = u.UpdateCRMContactDetails(CRMContactProperties, session, ctx)
		// if err != nil {
		// 	return err.Error()
		// }

		err = u.UpdateSessionLevel(ctx, GetDOBState, session.SessionID)
		if err != nil {
			return err.Error()
		}

		resp := "CON Please enter your date of birth in.\r\n"
		resp += "DDMMYYYY format e.g 14031996 for.\r\n"
		resp += "14th March 1992.\r\n"
		return resp
	}

	if session.Level == GetDOBState {
		err := utils.ValidateDateDigits(userResponse)
		if err != nil {
			return "CON The date of birth you entered is not valid, please try again in DDMMYYYY format e.g 14031996"
		}

		err = utils.ValidateDateLength(userResponse)
		if err != nil {
			return "CON The date of birth you entered is not valid, please try again in DDMMYYYY format e.g 14031996"
		}
		resp := utils.ValidateYearOfBirth(userResponse)
		if err != nil {
			return resp
		}

		date = userResponse

		// //Update CRM Contact DOB
		// day := date[0:2]
		// month := date[2:4]
		// year := date[4:8]

		// myDOB := day + "-" + month + "-" + year
		// userDOB, _ := time.Parse(layoutISO, myDOB)

		// var CRMContactProperties CRMDomain.ContactProperties
		// if date != "" {
		// 	CRMContactProperties.DateOfBirth = &userDOB
		// }

		// err = u.UpdateCRMContactDetails(CRMContactProperties, session, ctx)
		// if err != nil {
		// 	return err.Error()
		// }

		err = u.UpdateSessionLevel(ctx, GetPINState, session.SessionID)
		if err != nil {
			return err.Error()
		}
		return "CON Please enter a 4 digit PIN to secure your account"
	}

	if session.Level == GetPINState {
		// TODO FIXME check for empty response
		_, err := u.onboardingRepository.UpdateSessionPIN(ctx, session.SessionID, userResponse)
		if err != nil {
			return "END something wrong happened"
		}
		err = u.UpdateSessionLevel(ctx, SaveRecord, session.SessionID)
		if err != nil {
			return err.Error()
		}
		return "CON Please enter a 4 digit PIN again to confirm"

	}

	if session.Level == SaveRecord {
		if userResponse != session.PIN {
			err := u.UpdateSessionLevel(ctx, 4, session.SessionID)
			if err != nil {
				return "END something wrong happened"
			}
			return "CON Please enter a 4 digit PIN to secure your account"
		}
		day, _ := strconv.Atoi(date[0:2])
		month, _ := strconv.Atoi(date[2:4])
		year, _ := strconv.Atoi(date[4:8])
		dateofBirth := &base.Date{
			Month: month,
			Day:   day,
			Year:  year,
		}
		updateInput := &dto.UserProfileInput{
			DateOfBirth: dateofBirth,
			FirstName:   &userFirstName,
			LastName:    &userLastName,
		}

		err := u.CreateUsddUserProfile(ctx, session.PhoneNumber, session.PIN, updateInput)
		if err != nil {
			return "END something wrong happened"
		}
		// SEND CRM DATA
		data := struct {
			DateOfBirth base.Date
			FirstName   string
			LastName    string
			PhoneNumber string
		}{
			DateOfBirth: *dateofBirth,
			FirstName:   userFirstName,
			LastName:    userLastName,
			PhoneNumber: session.PhoneNumber,
		}

		if err := u.onboardingRepository.StageCRMPayload(ctx, dto.CRMStagingPayload{
			CRMUpdateContactPayload: data,
		}); err != nil {
			return "End something wrong happened"
		}

		err = u.UpdateSessionLevel(ctx, WelcomeMenu, session.SessionID)
		if err != nil {
			return "END something wrong happened"
		}
		resp := "CON Thanks for signing for Be.Well.\r\n"
		resp += "1. I want a cover.\r\n"
		resp += "2. Opt out from marketing messages.\r\n"
		resp += "3. Change PIN."
		return resp

	}

	if session.Level == WelcomeMenu {
		switch userResponse {
		case RegWantToBuyCoverInput:
			resp := "END We have recorded your request.\r\n"
			resp += "and one of the representatives will.\r\n"
			resp += "reach out to you. Thank you.\r\n"
			return resp
		case RegOptOutInput:
			option := "STOP"
			err := u.profile.SetOptOut(ctx, option, session.PhoneNumber)
			if err != nil {
				return err.Error()
			}
			resp := "CON We have successfully opted you.\r\n"
			resp += "marketing messages.\r\n"
			resp += "0. Go back home."
			return resp
		case RegChangePINInput:
			err := u.UpdateSessionLevel(ctx, UserPINState, session.SessionID)
			if err != nil {
				return "END something is wrong"
			}
			return u.HandleChangePIN(ctx, session, userResponse)

		default:
			return "CON Invalid choice. Please try again.\n1. Opt out from marketing messages \n2. Change PIN"
		}
	}
	return "END invalid input"
}

// CreateCRMContact ...
// func (u *Impl) CreateCRMContact(ctx context.Context, contact CRMDomain.CRMContact) error {

// 	bs, err := json.Marshal(contact)
// 	if err != nil {
// 		return err
// 	}

// 	err = u.pubsub.PublishToPubsub(
// 		ctx,
// 		u.pubsub.AddPubSubNamespace(pubsubmessaging.CreateCRMContact),
// 		bs,
// 	)
// 	if err != nil {
// 		log.Printf("unable to publish to Pub/Sub to create CRM contact: %v", err)
// 	}
// 	return nil

// }

// UpdateCRMContactDetails updates CRM contact
// func (u *Impl) UpdateCRMContactDetails(payload CRMDomain.ContactProperties, sessionDetails *domain.USSDLeadDetails, ctx context.Context) error {

// 	bs, err := json.Marshal(dto.UpdateContactPSMessage{
// 		Properties: payload,
// 		Phone:      sessionDetails.PhoneNumber,
// 	})
// 	if err != nil {
// 		return nil
// 	}

// 	err = u.pubsub.PublishToPubsub(
// 		ctx,
// 		u.pubsub.AddPubSubNamespace(pubsubmessaging.UpdateCRMContact),
// 		bs,
// 	)
// 	if err != nil {
// 		log.Printf("unable to publish to Pub/Sub to create CRM contact: %v", err)
// 	}
// 	return nil
// }
