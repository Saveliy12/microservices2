package converter

import (
	"github.com.Saveliy12/microservices2.git/inventory/internal/model"
	inventoryv1 "github.com.Saveliy12/microservices2.git/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func PartsToProto(parts []model.Part) []*inventoryv1.Part {
	result := make([]*inventoryv1.Part, 0, len(parts))
	for _, part := range parts {
		result = append(result, PartToProto(part))
	}

	return result
}

func FilterFromProto(filter *inventoryv1.PartsFilter) model.PartsFilter {
	if filter == nil {
		return model.PartsFilter{}
	}

	return model.PartsFilter{
		UUIDs:                 stringWrappersToStrings(filter.GetUuids()),
		Names:                 stringWrappersToStrings(filter.GetNames()),
		Categories:            categoriesFromProto(filter.GetCategories()),
		ManufacturerCountries: stringWrappersToStrings(filter.GetManufacturerCountries()),
		Tags:                  stringWrappersToStrings(filter.GetTags()),
	}
}

func PartToProto(part model.Part) *inventoryv1.Part {
	metadata := make(map[string]*inventoryv1.Value, len(part.Metadata))
	for key, value := range part.Metadata {
		metadata[key] = metadataValueToProto(value)
	}

	createdAt := (*timestamppb.Timestamp)(nil)
	if !part.CreatedAt.IsZero() {
		createdAt = timestamppb.New(part.CreatedAt)
	}

	updatedAt := (*timestamppb.Timestamp)(nil)
	if !part.UpdatedAt.IsZero() {
		updatedAt = timestamppb.New(part.UpdatedAt)
	}

	tags := make([]string, len(part.Tags))
	copy(tags, part.Tags)

	return &inventoryv1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      categoryToProto(part.Category),
		Dimensions: &inventoryv1.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		},
		Manufacturer: &inventoryv1.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		},
		Tags:      tags,
		Metadata:  metadata,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func categoryToProto(category model.Category) inventoryv1.Category {
	switch category {
	case model.CategoryEngine:
		return inventoryv1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventoryv1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return inventoryv1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventoryv1.Category_CATEGORY_WING
	default:
		return inventoryv1.Category_CATEGORY_UNSPECIFIED
	}
}

func categoryFromProto(category inventoryv1.Category) model.Category {
	switch category {
	case inventoryv1.Category_CATEGORY_ENGINE:
		return model.CategoryEngine
	case inventoryv1.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case inventoryv1.Category_CATEGORY_PORTHOLE:
		return model.CategoryPorthole
	case inventoryv1.Category_CATEGORY_WING:
		return model.CategoryWing
	default:
		return model.CategoryUnspecified
	}
}

func categoriesFromProto(categories []inventoryv1.Category) []model.Category {
	result := make([]model.Category, 0, len(categories))
	for _, category := range categories {
		result = append(result, categoryFromProto(category))
	}

	return result
}

func stringWrappersToStrings(values []*wrapperspb.StringValue) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, value.GetValue())
	}

	return result
}

func metadataValueToProto(value model.MetadataValue) *inventoryv1.Value {
	protoValue := &inventoryv1.Value{}

	if value.StringValue != nil {
		protoValue.StringValue = *value.StringValue
	}

	if value.Int64Value != nil {
		protoValue.Int64Value = *value.Int64Value
	}

	if value.DoubleValue != nil {
		protoValue.DoubleValue = *value.DoubleValue
	}

	if value.BoolValue != nil {
		protoValue.BoolValue = *value.BoolValue
	}

	return protoValue
}
