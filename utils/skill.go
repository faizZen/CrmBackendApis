package utils

import (
	"fmt"

	initializers "github.com/Zenithive/it-crm-backend/Initializers"
	"github.com/Zenithive/it-crm-backend/internal/graphql/generated"
	"github.com/Zenithive/it-crm-backend/models"
	"github.com/google/uuid"
)

func ConvertSkills(modelSkills []models.Skill) []*generated.Skill {
	skills := make([]*generated.Skill, len(modelSkills))
	for i, skill := range modelSkills {
		skills[i] = &generated.Skill{
			SkillID: skill.ID.String(),
			Name:    skill.Name,
		}
	}
	return skills
}

// func ConvertSkills(modelSkills []models.Skill) []*generated.Skill {

// 	skills := make([]*generated.Skill, len(modelSkills))
// 	for i, skill := range modelSkills {
// 		skills[i] = &generated.Skill{
// 			ID:   skill.ID.String(),
// 			Name: skill.Name,
// 		}
// 	}
// 	return skills
// }

// func FetchSkills(skillIDs []string) ([]models.Skill, error) {
// 	var skills []models.Skill
// 	for _, skillIDStr := range skillIDs {
// 		skillID, err := uuid.Parse(skillIDStr)
// 		if err != nil {
// 			return nil, fmt.Errorf("invalid skill ID %s: %w", skillIDStr, err)
// 		}

// 		var skill models.Skill
// 		if err := initializers.DB.First(&skill, "id = ?", skillID).Error; err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				return nil, fmt.Errorf("skill with ID %s not found", skillIDStr)
// 			}
// 			return nil, fmt.Errorf("error retrieving skill: %w", err)
// 		}
// 		skills = append(skills, skill)
// 	}
// 	return skills, nil

func FetchSkills(skillIDs []uint) ([]models.Skill, error) {
	var skills []models.Skill
	if err := initializers.DB.Find(&skills, "id IN ?", skillIDs).Error; err != nil {
		return nil, fmt.Errorf("error retrieving skills: %w", err)
	}
	return skills, nil
}
func ConvertSkill(s models.Skill) *generated.Skill {
	return &generated.Skill{
		SkillID:     s.ID.String(),
		Name:        s.Name,
		Description: *s.Description,
	}
}

// ConvertResourceSkills converts a slice of models.ResourceSkill into a slice of generated.ResourceSkill.
func ConvertResourceSkills(resourceSkills []models.ResourceSkill) []*generated.ResourceSkill {
	var result []*generated.ResourceSkill
	for _, rs := range resourceSkills {
		// Ensure that the Skill is loaded
		if rs.Skill.ID == uuid.Nil {
			// If not preloaded, you might choose to skip or handle the error.
			continue
		}

		generatedRS := &generated.ResourceSkill{
			Skill:           ConvertSkill(rs.Skill), // Convert the full skill data
			ExperienceYears: rs.ExperienceYears,
		}
		result = append(result, generatedRS)
	}
	return result
}
