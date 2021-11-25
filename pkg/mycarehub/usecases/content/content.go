package content

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/serverutils"
)

var (
	contentAPIEndpoint = serverutils.MustGetEnvVar("CONTENT_API_URL")
)

// IGetContent is used to fetch content from the CMS
type IGetContent interface {
	GetContent(ctx context.Context, categoryID *int, limit string) (*domain.Content, error)
}

// IContentCategoryList groups allthe content category listing methods
type IContentCategoryList interface {
	ListContentCategories(ctx context.Context) ([]*domain.ContentItemCategory, error)
}

// IShareContent is the interface for the ShareContent
type IShareContent interface {
	// TODO: update share count (increment)
	// TODO: add / check entry in ContentShares table
	// TODO: metrics
	ShareContent(ctx context.Context, input dto.ShareContentInput) (bool, error)
}

// IBookmarkContent is used to bookmark content
type IBookmarkContent interface {
	// TODO: update bookmark count (increment)
	// TODO: idempotence, with user ID i.e a user can only bookmark once
	// TODO: add / check entry in ContentBookmarks table
	// TODO: metrics
	BookmarkContent(ctx context.Context, userID string, contentID int) (bool, error)
}

// IUnBookmarkContent is used to unbookmark content
type IUnBookmarkContent interface {
	// TODO: update bookmark count (decrement)
	// TODO: idempotence, with user ID i.e a user can only remove something they bookmarked
	// TODO: remove entry from ContentBookmarks table if it exists...be forgiving (idempotence)
	// TODO: metrics
	UnBookmarkContent(ctx context.Context, userID string, contentID int) (bool, error)
}

// UseCasesContent holds the interfaces that are implemented within the content service
type UseCasesContent interface {
	IGetContent
	IContentCategoryList
	IShareContent
	IBookmarkContent
	IUnBookmarkContent
}

// UseCasesContentImpl represents content implementation
type UseCasesContentImpl struct {
	Update infrastructure.Update
	Query  infrastructure.Query
}

// NewUseCasesContentImplementation initializes a new contents service
func NewUseCasesContentImplementation(
	update infrastructure.Update,
	query infrastructure.Query,
) *UseCasesContentImpl {
	return &UseCasesContentImpl{
		Update: update,
		Query:  query,
	}
}

// GetContent fetches content from wagtail CMS. The category ID is optional and it is used to return content based
// on the category it belongs to. The limit field describes how many items will be rendered on the front end side.
func (u UseCasesContentImpl) GetContent(ctx context.Context, categoryID *int, limit string) (*domain.Content, error) {
	params := url.Values{}
	params.Add("type", "content.ContentItem")
	params.Add("limit", limit)
	params.Add("order", "-first_published_at")
	params.Add("fields", "'*")
	if categoryID != nil {
		params.Add("category", strconv.Itoa(*categoryID))
	}

	getContentEndpoint := fmt.Sprintf(contentAPIEndpoint + "/?" + params.Encode())
	var contentItems *domain.Content
	resp, err := utils.MakeRequest(ctx, http.MethodGet, getContentEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request")
	}

	dataResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %v", err)
	}

	err = json.Unmarshal(dataResponse, &contentItems)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return contentItems, nil
}

// ListContentCategories gets the list of all content categories
func (u *UseCasesContentImpl) ListContentCategories(ctx context.Context) ([]*domain.ContentItemCategory, error) {
	return u.Query.ListContentCategories(ctx)
}

// ShareContent enables a user to share a content
func (u *UseCasesContentImpl) ShareContent(ctx context.Context, input dto.ShareContentInput) (bool, error) {
	return u.Update.ShareContent(ctx, input)
}

// BookmarkContent increments the bookmark count for a content item
func (u UseCasesContentImpl) BookmarkContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return u.Update.BookmarkContent(ctx, userID, contentID)
}

// UnBookmarkContent decrements the bookmark count for a content item
func (u UseCasesContentImpl) UnBookmarkContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return u.Update.UnBookmarkContent(ctx, userID, contentID)
}