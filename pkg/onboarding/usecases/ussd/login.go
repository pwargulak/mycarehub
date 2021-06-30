package ussd

import (
	"context"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

// HandleLogin represents the workflow for authenticating a user
func (u *Impl) HandleLogin(ctx context.Context, session *domain.USSDLeadDetails, userResponse string) string {
	switch userResponse {
	case EmptyInput:
		resp := "CON Welcome to Be.Well.Please enter\r\n"
		resp += "your PIN to continue(enter 00 if\r\n"
		resp += "you forgot your PIN)\r\n"
		return resp

	case ForgotPINInput:
		err := u.UpdateSessionLevel(ctx, ForgotPINVerifyDate, session.SessionID)
		if err != nil {
			return "END Something went wrong. Please try again."
		}
		resp := "CON Please enter your date of birth in\r\n"
		resp += "DDMMYYYY format e.g 14031996 for\r\n"
		resp += "14th March 1992\r\n"
		resp += "to be able to reset PIN\r\n"
		return resp

	default:
		isLoggedIn, err := u.LoginInUser(ctx, session.PhoneNumber, userResponse, base.FlavourConsumer)
		if err != nil {
			return "END Something went wrong. Please try again."
		}
		if !isLoggedIn {
			resp := "CON The PIN you entered is not correct\r\n"
			resp += "Please try again (enter 00 if you\r\n"
			resp += "forgot your PIN)"
			return resp
		}
		err = u.UpdateSessionLevel(ctx, HomeMenuState, session.SessionID)
		if err != nil {
			return "END Something went wrong. Please try again."
		}
		userResponse := ""
		return u.HandleHomeMenu(ctx, HomeMenuState, session, userResponse)
	}

}

// LoginInUser authenticates a user to allow them proceed to the home menu
func (u *Impl) LoginInUser(
	ctx context.Context,
	phone string,
	PIN string,
	flavour base.Flavour,
) (bool, error) {

	profile, err := u.onboardingRepository.GetUserProfileByPrimaryPhoneNumber(
		ctx,
		phone,
		false,
	)
	if err != nil {
		return false, err
	}

	PINData, err := u.onboardingRepository.GetPINByProfileID(ctx, profile.ID)
	if err != nil {
		return false, err
	}
	if PINData == nil {
		return false, err
	}
	matched := u.pinExt.ComparePIN(PIN, PINData.Salt, PINData.PINNumber, nil)
	if !matched {
		return false, nil

	}
	return true, nil

}
