package lokalise

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
)

const (
	pathSegments = "segments"
)

// SegmentationService supports List, Retrieve and Update commands
type SegmentationService struct {
	BaseService

	listOpts     SegmentsListOptions
	retrieveOpts SegmentsRetrieveOptions
}

type Segment struct {
	SegmentNumber             int64               `json:"segment_number"`
	LanguageIso               string              `json:"language_iso"`
	ModifiedAt                string              `json:"modified_at"`
	ModifiedAtTimestamp       int64               `json:"modified_at_timestamp"`
	ModifiedBy                int64               `json:"modified_by"`
	ModifiedByEmail           string              `json:"modified_by_email"`
	Value                     string              `json:"value"`
	IsFuzzy                   bool                `json:"is_fuzzy"`
	IsReviewed                bool                `json:"is_reviewed"`
	ReviewedBy                int64               `json:"reviewed_by"`
	Words                     int64               `json:"words"`
	CustomTranslationStatuses []TranslationStatus `json:"custom_translation_statuses"`
}

type SegmentsResponse struct {
	WithProjectID
	Segments []Segment   `json:"segments"`
	Errors   []ErrorKeys `json:"error,omitempty"`
}

type SegmentResponse struct {
	WithProjectID
	KeyID       int64   `json:"key_id"`
	LanguageISO string  `json:"language_iso"`
	Segment     Segment `json:"segment"`
}

type SegmentUpdateRequest struct {
	Value                      string  `json:"value"` // could be string or json for plural keys.
	IsFuzzy                    *bool   `json:"is_fuzzy,omitempty"`
	IsReviewed                 *bool   `json:"is_reviewed,omitempty"`
	CustomTranslationStatusIds []int64 `json:"custom_translation_status_ids,omitempty"`
}

func (s *SegmentationService) List(projectID string, keyID int64, languageIso string) (r SegmentsResponse, err error) {
	resp, err := s.getWithOptions(
		s.Ctx(),
		fmt.Sprintf("%s/%s/%s/%d/%s/%s", pathProjects, projectID, pathKeys, keyID, pathSegments, languageIso),
		&r,
		s.ListOpts(),
	)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (s *SegmentationService) Retrieve(
	projectID string,
	keyID int64,
	languageIso string,
	segmentNumber int64,
) (r SegmentResponse, err error) {
	resp, err := s.getWithOptions(
		s.Ctx(),
		segmentPath(projectID, keyID, languageIso, segmentNumber),
		&r,
		s.RetrieveOpts(),
	)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (s *SegmentationService) Update(
	projectID string,
	keyID int64,
	languageIso string,
	segmentNumber int64,
	updateRequest SegmentUpdateRequest,
) (r SegmentResponse, err error) {
	resp, err := s.put(s.Ctx(), segmentPath(projectID, keyID, languageIso, segmentNumber), &r, updateRequest)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func segmentPath(projectID string, keyID int64, languageIso string, segmentNumber int64) string {
	return fmt.Sprintf(
		"%s/%s/%s/%d/%s/%s/%d",
		pathProjects,
		projectID,
		pathKeys,
		keyID,
		pathSegments,
		languageIso,
		segmentNumber,
	)
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Additional methods
// _____________________________________________________________________________________________________________________

type SegmentsListOptions struct {
	// Possible values are 1 and 0.
	DisableReferences  uint8 `url:"disable_references,omitempty"`
	FilterIsReviewed   uint8 `url:"filter_is_reviewed,omitempty"`
	FilterUnverified   uint8 `url:"filter_unverified,omitempty"`
	FilterUntranslated uint8 `url:"filter_untranslated,omitempty"`

	FilterQAIssues string `url:"filter_qa_issues,omitempty"`
}

func (s *SegmentationService) ListOpts() SegmentsListOptions        { return s.listOpts }
func (s *SegmentationService) SetListOptions(o SegmentsListOptions) { s.listOpts = o }
func (s *SegmentationService) WithListOptions(o SegmentsListOptions) *SegmentationService {
	s.listOpts = o
	return s
}

func (options SegmentsListOptions) Apply(req *resty.Request) {
	v, _ := query.Values(options)
	req.SetQueryString(v.Encode())
}

type SegmentsRetrieveOptions struct {
	DisableReferences uint8 `url:"disable_references,omitempty"`
}

func (options SegmentsRetrieveOptions) Apply(req *resty.Request) {
	v, _ := query.Values(options)
	req.SetQueryString(v.Encode())
}

func (s *SegmentationService) RetrieveOpts() SegmentsRetrieveOptions        { return s.retrieveOpts }
func (s *SegmentationService) SetRetrieveOptions(o SegmentsRetrieveOptions) { s.retrieveOpts = o }
func (s *SegmentationService) WithRetrieveOptions(o SegmentsRetrieveOptions) *SegmentationService {
	s.retrieveOpts = o
	return s
}
