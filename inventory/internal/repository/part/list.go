package part

import (
	"context"

	"github.com.Saveliy12/microservices2.git/inventory/internal/model"
	repositoryConverter "github.com.Saveliy12/microservices2.git/inventory/internal/repository/converter"
	repositoryModel "github.com.Saveliy12/microservices2.git/inventory/internal/repository/model"
)

func (r *Repository) ListParts(_ context.Context, filter model.PartsFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoFilter := repositoryConverter.ModelFilterToRepoFilter(filter)

	parts := make([]repositoryModel.Part, 0, len(r.parts))
	for _, part := range r.parts {
		if !matchesFilter(part, repoFilter) {
			continue
		}

		parts = append(parts, part)
	}

	return repositoryConverter.RepoPartsToModel(parts), nil
}

func matchesFilter(part repositoryModel.Part, filter repositoryModel.PartsFilter) bool {
	if len(filter.UUIDs) > 0 && !containsString(filter.UUIDs, part.UUID) {
		return false
	}

	if len(filter.Names) > 0 && !containsString(filter.Names, part.Name) {
		return false
	}

	if len(filter.Categories) > 0 && !containsCategory(filter.Categories, part.Category) {
		return false
	}

	if len(filter.ManufacturerCountries) > 0 && !containsString(filter.ManufacturerCountries, part.Manufacturer.Country) {
		return false
	}

	if len(filter.Tags) > 0 && !hasAnyTag(filter.Tags, part.Tags) {
		return false
	}

	return true
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}

func containsCategory(values []repositoryModel.Category, target repositoryModel.Category) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}

func hasAnyTag(filterTags []string, partTags []string) bool {
	for _, tag := range partTags {
		if containsString(filterTags, tag) {
			return true
		}
	}

	return false
}
