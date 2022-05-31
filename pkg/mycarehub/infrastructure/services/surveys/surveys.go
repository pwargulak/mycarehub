package surveys

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/serverutils"
)

var (
	surveysSystemEmail    = serverutils.MustGetEnvVar("SURVEYS_SYSTEM_EMAIL")
	surveysSystemPassword = serverutils.MustGetEnvVar("SURVEYS_SYSTEM_PASSWORD")
)

// Surveys is the interface that defines the methods that are required to access the surveys client
type Surveys interface {
	MakeRequest(ctx context.Context, payload domain.RequestHelperPayload) (*http.Response, error)
	ListSurveyForms(ctx context.Context, projectID int) ([]*domain.SurveyForm, error)
	GetSurveyForm(ctx context.Context, projectID int, formID string) (*domain.SurveyForm, error)
	GeneratePublicAccessLink(ctx context.Context, input dto.SurveyLinkInput) (*dto.SurveyPublicLink, error)
	GetSubmissions(ctx context.Context, input dto.VerifySurveySubmissionInput) ([]domain.Submission, error)
	DeletePublicAccessLink(ctx context.Context, input dto.VerifySurveySubmissionInput) error
	ListSubmitters(ctx context.Context, projectID int, formID string) ([]domain.Submitter, error)
}

// Impl implements the Surveys interface
type Impl struct {
	client domain.SurveysClient
}

// NewSurveysImpl returns a new Impl
func NewSurveysImpl(client domain.SurveysClient) Surveys {
	return &Impl{
		client: client,
	}
}

// MakeRequest performs a http request and returns a response
func (s *Impl) MakeRequest(ctx context.Context, payload domain.RequestHelperPayload) (*http.Response, error) {
	client := s.client.HTTPClient

	// A GET or DELETE request should not send data when doing a request. We should use query parameters
	// instead of having a request body. In some cases where a GET request has an empty body {},
	// it might result in status code 400 with the error:
	//  `Your client has issued a malformed or illegal request. That’s all we know.`
	if payload.Method == http.MethodGet || payload.Method == http.MethodDelete {
		req, reqErr := http.NewRequestWithContext(ctx, payload.Method, payload.Path, nil)
		if reqErr != nil {
			return nil, reqErr
		}

		req.SetBasicAuth(surveysSystemEmail, surveysSystemPassword)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Extended-Metadata", "true")
		return client.Do(req)
	}

	encoded, err := json.Marshal(payload.Body)
	if err != nil {
		return nil, err
	}

	p := bytes.NewBuffer(encoded)
	req, reqErr := http.NewRequestWithContext(ctx, payload.Method, payload.Path, p)
	if reqErr != nil {
		return nil, reqErr
	}

	req.SetBasicAuth(surveysSystemEmail, surveysSystemPassword)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

// ListSurveyForms returns a list of survey forms
func (s *Impl) ListSurveyForms(ctx context.Context, projectID int) ([]*domain.SurveyForm, error) {

	payload := domain.RequestHelperPayload{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("%s/v1/projects/%s/forms", s.client.BaseURL, strconv.Itoa(projectID)),
	}

	resp, err := s.MakeRequest(ctx, payload)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("surveys: ListSurveyForms error: status code: %d", resp.StatusCode)
	}

	var surveyForms []*dto.SurveyForm
	err = json.NewDecoder(resp.Body).Decode(&surveyForms)
	if err != nil {
		return nil, err
	}

	var surveyFormsDomain []*domain.SurveyForm
	for _, surveyForm := range surveyForms {
		surveyFormsDomain = append(surveyFormsDomain, &domain.SurveyForm{
			ProjectID: surveyForm.ProjectID,
			XMLFormID: surveyForm.XMLFormID,
			Name:      surveyForm.Name,
			EnketoID:  surveyForm.EnketoID,
		})
	}

	return surveyFormsDomain, nil
}

// GetSurveyForm returns a survey form
func (s *Impl) GetSurveyForm(ctx context.Context, projectID int, formID string) (*domain.SurveyForm, error) {

	payload := domain.RequestHelperPayload{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("%s/v1/projects/%s/forms/%s", s.client.BaseURL, strconv.Itoa(projectID), formID),
	}

	resp, err := s.MakeRequest(ctx, payload)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("surveys: GetSurveyForm error: status code: %d", resp.StatusCode)
	}

	var surveyForm dto.SurveyForm
	err = json.NewDecoder(resp.Body).Decode(&surveyForm)

	if err != nil {
		return nil, err
	}

	return &domain.SurveyForm{
		ProjectID: surveyForm.ProjectID,
		XMLFormID: surveyForm.XMLFormID,
		Name:      surveyForm.Name,
		EnketoID:  surveyForm.EnketoID,
	}, nil
}

// GeneratePublicAccessLink returns a survey public link
func (s *Impl) GeneratePublicAccessLink(ctx context.Context, input dto.SurveyLinkInput) (*dto.SurveyPublicLink, error) {
	payload := domain.RequestHelperPayload{
		Method: http.MethodPost,
		Path:   fmt.Sprintf("%s/v1/projects/%s/forms/%s/public-links", s.client.BaseURL, strconv.Itoa(input.ProjectID), input.FormID),
		Body:   map[string]interface{}{"once": input.OnceOnly, "displayName": input.DisplayName},
	}

	resp, err := s.MakeRequest(ctx, payload)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("surveys: GeneratePublicAccessLink error: status code: %d", resp.StatusCode)
	}

	var surveyPublicLink dto.SurveyPublicLink
	err = json.NewDecoder(resp.Body).Decode(&surveyPublicLink)
	if err != nil {
		return nil, err
	}

	return &surveyPublicLink, nil
}

// GetSubmissions returns a list of all survey submissions
func (s *Impl) GetSubmissions(ctx context.Context, input dto.VerifySurveySubmissionInput) ([]domain.Submission, error) {
	url := fmt.Sprintf("%s/v1/projects/%v/forms/%s/submissions", serverutils.MustGetEnvVar("SURVEYS_BASE_URL"), input.ProjectID, input.FormID)
	payload := domain.RequestHelperPayload{
		Method: http.MethodGet,
		Path:   url,
	}
	resp, reqErr := s.MakeRequest(ctx, payload)
	if reqErr != nil {
		helpers.ReportErrorToSentry(reqErr)
		return nil, reqErr
	}
	defer resp.Body.Close()
	respBody, respErr := ioutil.ReadAll(resp.Body)
	if respErr != nil {
		helpers.ReportErrorToSentry(respErr)
		return nil, respErr
	}

	var submissions []domain.Submission
	err := json.Unmarshal(respBody, &submissions)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to unmarshal submissions: %w", err)
	}

	return submissions, nil
}

// DeletePublicAccessLink deletes the survey public link
func (s *Impl) DeletePublicAccessLink(ctx context.Context, input dto.VerifySurveySubmissionInput) error {
	url := fmt.Sprintf("%s/v1/projects/%v/forms/%s/public-links/%v", serverutils.MustGetEnvVar("SURVEYS_BASE_URL"), input.ProjectID, input.FormID, input.SubmitterID)
	payload := domain.RequestHelperPayload{
		Method: http.MethodDelete,
		Path:   url,
	}
	_, reqErr := s.MakeRequest(ctx, payload)
	if reqErr != nil {
		helpers.ReportErrorToSentry(reqErr)
		return reqErr
	}

	return nil
}

// ListSubmitters returns a a listing of all known submitting actors to a given Form. Each Actor that has submitted to the given Form will be returned once.
func (s *Impl) ListSubmitters(ctx context.Context, projectID int, formID string) ([]domain.Submitter, error) {
	url := fmt.Sprintf("%s/v1/projects/%v/forms/%s/submissions/submitters", serverutils.MustGetEnvVar("SURVEYS_BASE_URL"), projectID, formID)
	payload := domain.RequestHelperPayload{
		Method: http.MethodGet,
		Path:   url,
	}
	resp, reqErr := s.MakeRequest(ctx, payload)
	if reqErr != nil {
		helpers.ReportErrorToSentry(reqErr)
		return nil, reqErr
	}
	defer resp.Body.Close()
	respBody, respErr := ioutil.ReadAll(resp.Body)
	if respErr != nil {
		helpers.ReportErrorToSentry(respErr)
		return nil, respErr
	}

	var submitters []domain.Submitter
	err := json.Unmarshal(respBody, &submitters)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to unmarshal submitters: %w", err)
	}

	return submitters, nil
}
