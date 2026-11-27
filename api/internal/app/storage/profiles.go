package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfilesStorage struct {
	db *pgxpool.Pool
}

func NewProfilesStorage(db *pgxpool.Pool) *ProfilesStorage {
	return &ProfilesStorage{db: db}
}

func (p *ProfilesStorage) AddUserProfileSection(ctx context.Context, section models.ProfileSection) (int64, error) {
	query := `
		INSERT INTO
		profile_sections (user_id, section_key, title, description)
		VALUES ($1, $2, $3, $4) 
		RETURNING id
	`
	var sectionID int64

	err := p.db.QueryRow(ctx, query, section.UserID, section.SectionKey, section.Title, section.Description).Scan(&sectionID)

	if err != nil {
		return 0, err
	}

	return sectionID, nil
}

func (p *ProfilesStorage) AddUserProfileCommunicationSection(ctx context.Context, userID int64) error {
	query := `
		INSERT INTO
		profile_sections (user_id, section_key, title, description)
		VALUES (
			$1,
			'communication',
			'Общение', 
			'Хранит имя или предпочитаемое обращение, форму «ты» или «вы», а также пожелания к тону, длине и формату ответов. Используй эти настройки при формулировке каждого ответа.'
		)
		ON CONFLICT (user_id, section_key) DO NOTHING; 
	`

	_, err := p.db.Exec(ctx, query, userID)
	return err
}

func (p *ProfilesStorage) GetUserProfileCommunicationSection(ctx context.Context, userID int64) (*models.DBProfileSection, error) {
	query := `
		SELECT
			id,
			section_key,
			title,
			description
		FROM
			profile_sections
		WHERE 
			user_id = $1 AND section_key = 'communication'
	`
	var section models.DBProfileSection
	err := p.db.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&section.ID,
		&section.SectionKey,
		&section.Title,
		&section.Description,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &section, nil
}

func (p *ProfilesStorage) GetUserProfileSections(ctx context.Context, userID int64) ([]models.DBProfileSection, error) {
	query := `
		SELECT
			id,
			section_key,
			title,
			description
		FROM
			profile_sections
		WHERE 
			user_id = $1
	`

	rows, err := p.db.Query(ctx, query, userID)

	if err != nil {
		fmt.Print(err)

		return nil, err
	}

	defer rows.Close()

	sections := make([]models.DBProfileSection, 0)

	for rows.Next() {
		var section models.DBProfileSection

		err := rows.Scan(
			&section.ID,
			&section.SectionKey,
			&section.Title,
			&section.Description,
		)
		if err != nil {
			fmt.Print(err)

			return nil, err
		}

		sections = append(sections, section)
	}

	if err := rows.Err(); err != nil {
		fmt.Print(err)
		return nil, err
	}

	return sections, nil
}

func (p *ProfilesStorage) AddUserProfileFact(ctx context.Context, fact models.ProfileFact) error {
	query := `
		INSERT INTO profile_facts (
			user_id,
			section_id,
			fact_key,
			value,
			source_message_id
		)
		VALUES ($1, $2, $3, $4, $5) 
	`

	_, err := p.db.Exec(ctx, query,
		fact.UserID, fact.SectionID, fact.FactKey, fact.Value, fact.SourceMessageID,
	)
	return err
}

func (p *ProfilesStorage) UpdateUserProfileFact(ctx context.Context, fact models.ProfileFactPatch) (uint32, error) {
	query := `
	UPDATE profile_facts
	SET 
		value = CASE WHEN @set_value THEN @value ELSE value END,
		updated_at = now ()
	WHERE user_id = @user_id AND section_id = @section_id AND fact_key = @fact_key
	RETURNING section_id;
	`
	args := pgx.StrictNamedArgs{
		"user_id":    fact.UserID,
		"set_value":  fact.SetValue,
		"value":      fact.Value,
		"section_id": fact.SectionID,
		"fact_key":   fact.FactKey,
	}

	var updatedID uint32

	err := p.db.QueryRow(ctx, query, args).Scan(&updatedID)

	if err != nil {
		return 0, err
	}

	return updatedID, nil
}

func (p *ProfilesStorage) GetUserProfileFactsBySection(ctx context.Context, userID int64, sectionID int64) ([]models.ProfileFact, error) {
	query := `
		SELECT
			fact_key,
			value,
			source_message_id
		FROM
			profile_facts
		WHERE 
			user_id = $1 AND section_id=$2
	`
	rows, err := p.db.Query(ctx, query, userID, sectionID)

	if err != nil {
		fmt.Print(err)

		return nil, err
	}

	defer rows.Close()

	facts := make([]models.ProfileFact, 0)

	for rows.Next() {
		var fact models.ProfileFact

		err := rows.Scan(
			&fact.FactKey,
			&fact.Value,
			&fact.SourceMessageID,
		)

		if err != nil {
			fmt.Print(err)

			return nil, err
		}

		facts = append(facts, fact)
	}

	if err := rows.Err(); err != nil {
		fmt.Print(err)
		return nil, err
	}

	return facts, nil
}
