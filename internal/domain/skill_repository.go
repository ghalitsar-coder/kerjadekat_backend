package domain

import "context"

// SkillRepository reads skill master data.
type SkillRepository interface {
	ListCategories(ctx context.Context) ([]SkillCategory, error)
	FindCategoryByID(ctx context.Context, id int) (*SkillCategory, error)
}
