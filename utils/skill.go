package utils

import (
	"errors"
	"fmt"

	initializers "github.com/Zenithive/it-crm-backend/Initializers"
	"github.com/Zenithive/it-crm-backend/internal/graphql/generated"
	"github.com/Zenithive/it-crm-backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func ConvertSkills(modelSkills []models.Skill) []*generated.Skill {
	skills := make([]*generated.Skill, len(modelSkills))
	for i, skill := range modelSkills {
		skills[i] = &generated.Skill{
			ID:   skill.ID.String(),
			Name: skill.Name,
		}
	}
	return skills
}

func FetchSkills(skillIDs []string) ([]models.Skill, error) {
	var skills []models.Skill
	for _, skillIDStr := range skillIDs {
		skillID, err := uuid.Parse(skillIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid skill ID %s: %w", skillIDStr, err)
		}

		var skill models.Skill
		if err := initializers.DB.First(&skill, "id = ?", skillID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("skill with ID %s not found", skillIDStr)
			}
			return nil, fmt.Errorf("error retrieving skill: %w", err)
		}
		skills = append(skills, skill)
	}
	return skills, nil
}
