package converter

import (
	"github.com.Saveliy12/microservices2.git/inventory/internal/model"
	repositoryModel "github.com.Saveliy12/microservices2.git/inventory/internal/repository/model"
)

func RepoPartToModel(part repositoryModel.Part) model.Part {
	metadata := make(map[string]model.MetadataValue, len(part.Metadata))
	for key, value := range part.Metadata {
		metadata[key] = model.MetadataValue{
			StringValue: value.StringValue,
			Int64Value:  value.Int64Value,
			DoubleValue: value.DoubleValue,
			BoolValue:   value.BoolValue,
		}
	}

	tags := make([]string, len(part.Tags))
	copy(tags, part.Tags)

	return model.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      model.Category(part.Category),
		Dimensions: model.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		},
		Manufacturer: model.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		},
		Tags:      tags,
		Metadata:  metadata,
		CreatedAt: part.CreatedAt,
		UpdatedAt: part.UpdatedAt,
	}
}

func RepoPartsToModel(parts []repositoryModel.Part) []model.Part {
	result := make([]model.Part, 0, len(parts))
	for _, part := range parts {
		result = append(result, RepoPartToModel(part))
	}

	return result
}

func ModelPartToRepoPart(part model.Part) repositoryModel.Part {
	metadata := make(map[string]repositoryModel.MetadataValue, len(part.Metadata))
	for key, value := range part.Metadata {
		metadata[key] = repositoryModel.MetadataValue{
			StringValue: value.StringValue,
			Int64Value:  value.Int64Value,
			DoubleValue: value.DoubleValue,
			BoolValue:   value.BoolValue,
		}
	}

	tags := make([]string, len(part.Tags))
	copy(tags, part.Tags)

	return repositoryModel.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      repositoryModel.Category(part.Category),
		Dimensions: repositoryModel.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		},
		Manufacturer: repositoryModel.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		},
		Tags:      tags,
		Metadata:  metadata,
		CreatedAt: part.CreatedAt,
		UpdatedAt: part.UpdatedAt,
	}
}

func ModelFilterToRepoFilter(filter model.PartsFilter) repositoryModel.PartsFilter {
	uuids := make([]string, len(filter.UUIDs))
	copy(uuids, filter.UUIDs)

	names := make([]string, len(filter.Names))
	copy(names, filter.Names)

	manufacturerCountries := make([]string, len(filter.ManufacturerCountries))
	copy(manufacturerCountries, filter.ManufacturerCountries)

	tags := make([]string, len(filter.Tags))
	copy(tags, filter.Tags)

	categories := make([]repositoryModel.Category, 0, len(filter.Categories))
	for _, category := range filter.Categories {
		categories = append(categories, repositoryModel.Category(category))
	}

	return repositoryModel.PartsFilter{
		UUIDs:                 uuids,
		Names:                 names,
		Categories:            categories,
		ManufacturerCountries: manufacturerCountries,
		Tags:                  tags,
	}
}
