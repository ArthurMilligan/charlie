package services

import (
	"context"
	"encoding/json"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
)

type Profiles struct {
	ProfilesStorage interfaces.ProfilesStorage
}

func NewProfileService(p interfaces.ProfilesStorage) *Profiles {
	return &Profiles{
		ProfilesStorage: p,
	}
}

func (p *Profiles) AddUserProfileSection(
	ctx context.Context,
	userID int64,
	sectionKey string,
	title string,
	description string,
) (int64, error) {
	if title == "" || userID == 0 || sectionKey == "" || description == "" {
		return 0, ErrInvalidArgs
	}

	var section = models.ProfileSection{
		UserID:      userID,
		SectionKey:  sectionKey,
		Title:       title,
		Description: description,
	}

	sectionID, err := p.ProfilesStorage.AddUserProfileSection(ctx, section)

	if err != nil {
		return 0, err
	}

	return sectionID, nil
}

func (p *Profiles) GetUserProfileCommunicationFacts(
	ctx context.Context,
	userID int64,
) ([]models.ProfileFact, error) {
	if userID == 0 {
		return []models.ProfileFact{}, ErrInvalidArgs
	}

	section, err := p.ProfilesStorage.GetUserProfileCommunicationSection(ctx, userID)

	if err != nil {
		return []models.ProfileFact{}, err
	}

	if section == nil {
		err := p.ProfilesStorage.AddUserProfileCommunicationSection(ctx, userID)
		if err != nil {
			return []models.ProfileFact{}, err
		}
		return []models.ProfileFact{}, nil
	}

	facts, err := p.ProfilesStorage.GetUserProfileFactsBySection(ctx, userID, section.ID)

	if err != nil {
		return []models.ProfileFact{}, err
	}

	return facts, nil
}

func (p *Profiles) GetUserProfileSections(ctx context.Context, userID int64) ([]models.DBProfileSection, error) {
	if userID == 0 {
		return nil, ErrInvalidArgs
	}

	sections, err := p.ProfilesStorage.GetUserProfileSections(ctx, userID)

	if err != nil {
		return nil, err
	}

	return sections, nil
}

func (p *Profiles) AddUserProfileFact(
	ctx context.Context,
	userID int64,
	sectionID int64,
	factKey string,
	value json.RawMessage,
	sourceMessageID int64,
) error {
	if sectionID == 0 || userID == 0 || factKey == "" || sourceMessageID == 0 {
		return ErrInvalidArgs
	}

	var fact = models.ProfileFact{
		UserID:          userID,
		SectionID:       sectionID,
		FactKey:         factKey,
		Value:           value,
		SourceMessageID: sourceMessageID,
	}

	err := p.ProfilesStorage.AddUserProfileFact(ctx, fact)

	if err != nil {
		return err
	}

	return nil
}

func (p *Profiles) UpdateUserProfileFact(
	ctx context.Context,
	userID int64,
	sectionID int64,
	factKey string,
	value *json.RawMessage,
	setValue bool,
) error {
	if sectionID == 0 || userID == 0 || factKey == "" {
		return ErrInvalidArgs
	}

	var factPatch = models.ProfileFactPatch{
		UserID:    userID,
		SectionID: sectionID,
		FactKey:   factKey,
		Value:     value,
		SetValue:  setValue,
	}

	_, err := p.ProfilesStorage.UpdateUserProfileFact(ctx, factPatch)

	if err != nil {
		return err
	}

	return nil
}

func (p *Profiles) GetUserProfileFactsBySection(ctx context.Context, userID int64, sectionID int64) ([]models.ProfileFact, error) {
	if userID == 0 || sectionID == 0 {
		return nil, ErrInvalidArgs
	}

	facts, err := p.ProfilesStorage.GetUserProfileFactsBySection(ctx, userID, sectionID)

	if err != nil {
		return nil, err
	}

	return facts, nil
}
