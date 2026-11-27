package interfaces

import (
	"context"
	"encoding/json"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

type ProfilesService interface {
	AddUserProfileSection(
		ctx context.Context,
		userID int64,
		sectionKey string,
		title string,
		description string,
	) (int64, error)
	GetUserProfileSections(ctx context.Context, userID int64) ([]models.DBProfileSection, error)
	AddUserProfileFact(
		ctx context.Context,
		userID int64,
		sectionID int64,
		factKey string,
		value json.RawMessage,
		sourceMessageID int64,
	) error
	GetUserProfileCommunicationFacts(
		ctx context.Context,
		userID int64,
	) ([]models.ProfileFact, error)
	UpdateUserProfileFact(
		ctx context.Context,
		userID int64,
		sectionID int64,
		factKey string,
		value *json.RawMessage,
		setValue bool,
	) error
	GetUserProfileFactsBySection(ctx context.Context, userID int64, sectionID int64) ([]models.ProfileFact, error)
}
