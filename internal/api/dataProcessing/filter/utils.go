package filter

import (
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/base"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/enum"
)

func GetFilterRules(
	baseStruct any,
	objectName string,
	customRules map[string]map[string]enum.ValidateType) map[string]map[string]enum.ValidateType {

	rulesMap := map[string]map[string]enum.ValidateType{}

	switch baseStruct.(type) {
	case base.EntityWithIdKey:
		rulesMap[objectName] = map[string]enum.ValidateType{
			"id":         enum.TYPE_UUID,
			"created_at": enum.TYPE_DATA,
			"updated_at": enum.TYPE_DATA,
		}

	case base.EntityWithIdKeyUniqueIndex:
		rulesMap[objectName] = map[string]enum.ValidateType{
			"id":         enum.TYPE_UUID,
			"created_at": enum.TYPE_DATA,
			"updated_at": enum.TYPE_DATA,
		}

	case base.ArchivableEntityWithIdKey:
		rulesMap[objectName] = map[string]enum.ValidateType{
			"id":          enum.TYPE_UUID,
			"created_at":  enum.TYPE_DATA,
			"updated_at":  enum.TYPE_DATA,
			"archived_at": enum.TYPE_DATA,
		}

	case base.EntityWithIntegerKey:
		rulesMap[objectName] = map[string]enum.ValidateType{
			"id":          enum.TYPE_INT,
			"created_at":  enum.TYPE_DATA,
			"updated_at":  enum.TYPE_DATA,
			"archived_at": enum.TYPE_DATA,
		}
	}

	for key, value := range customRules {
		if key == objectName {
			for nameColumn, typeColumn := range value {
				rulesMap[objectName][nameColumn] = typeColumn
			}
		} else {
			rulesMap[key] = value
		}
	}

	return rulesMap
}
