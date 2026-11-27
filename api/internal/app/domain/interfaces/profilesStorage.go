package interfaces

import (
	"context"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

type ProfilesStorage interface {
	AddUserProfileSection(ctx context.Context, section models.ProfileSection) (int64, error)
	AddUserProfileCommunicationSection(ctx context.Context, userID int64) error
	GetUserProfileSections(ctx context.Context, userID int64) ([]models.DBProfileSection, error)
	GetUserProfileCommunicationSection(ctx context.Context, userID int64) (*models.DBProfileSection, error)
	AddUserProfileFact(ctx context.Context, fact models.ProfileFact) error
	UpdateUserProfileFact(ctx context.Context, fact models.ProfileFactPatch) (uint32, error)
	GetUserProfileFactsBySection(ctx context.Context, userID int64, sectionID int64) ([]models.ProfileFact, error)
}
