package otp

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"

	"gopkg.in/yaml.v2"
)

const otpService = "otp"

// OTP service endpoints
const (
	SendRetryOtp = "internal/send_retry_otp/"
	SendOtp      = "internal/send_otp/"
)

// ServiceOTP represent the business logic required for management of OTP
type ServiceOTP interface {
	GenerateAndSendOTP(
		ctx context.Context,
		phone string,
	) (string, error)
	SendRetryOTP(
		ctx context.Context,
		msisdn string,
		retryStep int,
	) (string, error)
	VerifyOTP(ctx context.Context, phone, OTP string) (bool, error)
}

// ServiceOTPImpl represents OTP usecases
type ServiceOTPImpl struct {
	Otp *base.InterServiceClient
}

// NewOTPService returns new instance of ServiceOTPImpl
func NewOTPService(r repository.OnboardingRepository) ServiceOTP {

	var config base.DepsConfig
	//os file and parse it to go type
	file, err := ioutil.ReadFile(filepath.Clean(base.PathToDepsFile()))
	if err != nil {
		log.Errorf("error occured while opening deps file %v", err)
		os.Exit(1)
	}

	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Errorf("failed to unmarshal yaml config file %v", err)
		os.Exit(1)
	}

	var otpClient *base.InterServiceClient
	otpClient, err = base.SetupISCclient(config, otpService)
	if err != nil {
		log.Panicf("unable to initialize otp inter service client: %s", err)

	}

	return &ServiceOTPImpl{Otp: otpClient}
}

// GenerateAndSendOTP creates a new otp and sends it to the provided phone number.
func (o *ServiceOTPImpl) GenerateAndSendOTP(
	ctx context.Context,
	phone string,
) (string, error) {
	body := map[string]interface{}{
		"msisdn": phone,
	}
	defaultOTP := ""
	resp, err := o.Otp.MakeRequest(http.MethodPost, SendOtp, body)
	if err != nil {
		return defaultOTP, fmt.Errorf("unable to generate and send otp: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return defaultOTP, fmt.Errorf(
			"unable to generate and send otp, with status code %v", resp.StatusCode,
		)
	}
	code, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return defaultOTP, fmt.Errorf("unable to convert response to string: %v", err)
	}

	var OTP string
	err = json.Unmarshal(code, &OTP)
	if err != nil {
		return defaultOTP, fmt.Errorf("failed to unmarshal OTP: %v", err)
	}

	return OTP, nil
}

// SendRetryOTP generates fallback OTPs when Africa is talking sms fails
func (o *ServiceOTPImpl) SendRetryOTP(
	ctx context.Context,
	msisdn string,
	retryStep int,
) (string, error) {
	phoneNumber, err := base.NormalizeMSISDN(msisdn)
	if err != nil {
		return "", fmt.Errorf("unable to normalize the msisdn: %v", err)
	}

	body := map[string]interface{}{
		"msisdn":    phoneNumber,
		"retryStep": retryStep,
	}
	resp, err := o.Otp.MakeRequest(http.MethodPost, SendRetryOtp, body)
	if err != nil {
		return "", fmt.Errorf("unable to generate and send fallback otp: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"unable to generate and send fallback otp, with status code %v",
			resp.StatusCode,
		)
	}

	code, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to convert response to string: %v", err)
	}

	var RetryOTP string
	err = json.Unmarshal(code, &RetryOTP)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal OTP: %v", err)
	}

	return RetryOTP, nil
}

// VerifyOTP takes a phone number and an OTP an checks for the vlidity of the OTP code
func (o *ServiceOTPImpl) VerifyOTP(ctx context.Context, phone, OTP string) (bool, error) {
	return base.VerifyOTP(phone, OTP, o.Otp)
}
