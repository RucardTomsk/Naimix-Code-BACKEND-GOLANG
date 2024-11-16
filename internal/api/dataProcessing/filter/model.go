package filter

import (
	"errors"
	"fmt"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/enum"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/helpers"
	"github.com/google/uuid"
	"strconv"
	"strings"
)

const (
	OperatorEquality      = "[eq]"  // =
	OperatorNotEquality   = "[ne]"  // !=
	OperatorLowerThan     = "[lt]"  // <
	OperatorLowerThanEq   = "[lte]" // <=
	OperatorGreaterThan   = "[gt]"  // >
	OperatorGreaterThanEq = "[gte]" // >=
	OperatorLike          = "[like]"
	OperatorIN            = "[in]"
)

type Options struct {
	IsToApply bool
	Fields    []Field
}

type Field struct {
	Name     string
	Value    string
	Operator string
	Type     enum.ValidateType
}

func (o *Options) AddField(name, value, operator string) {
	name = strings.TrimSpace(name)
	value = strings.TrimSpace(value)
	o.Fields = append(o.Fields, Field{
		Name:     name,
		Value:    value,
		Operator: operator,
	})
}

// ValidateField
//
// The method takes a map of field names and their corresponding data types as input and returns an error as output.
//
// This method is also a method attached to the Options struct. It is used to validate the fields in the struct based on the data types defined in the input map.
// The method first builds a string of all the field names in the struct and creates a map of the corresponding data types using the input map.
//
// The method then iterates over each field in the struct and checks if it is a valid filtering field.
// If the field is not a valid filtering field, the method returns an error. Otherwise, the method checks if the data type defined in the input map matches the data type of the field.
// If there is a match, the method updates the field data type and value as necessary.
// If there is no match, the method returns an error.
//
// If the field name is "id", the data type is set to enum.TYPE_STRING.
// If the data type is enum.TYPE_DATA, the method calls the parseUnixTimeStamp function to parse the field value and convert it into a formatted string.
// The data type is then set to enum.TYPE_DATA.
//
// The method returns a nil error if all fields are validated successfully.
// If there is an error during validation, the method returns an error containing a message describing the issue.
func (o *Options) ValidateField(filterRules map[string]map[string]enum.ValidateType) error {

	columnRules := make(map[string]struct {
		NameTable string
		Type      enum.ValidateType
	})

	for nameTable, rules := range filterRules {
		for nameColumn, typeColumn := range rules {
			columnRules[nameColumn] = struct {
				NameTable string
				Type      enum.ValidateType
			}{NameTable: nameTable, Type: typeColumn}
		}
	}

	for index, value := range o.Fields {
		switch columnRules[value.Name].Type {
		case enum.TYPE_INT:
			o.Fields[index].Name = columnRules[value.Name].NameTable + "." + value.Name
			o.Fields[index].Type = enum.TYPE_INT

		case enum.TYPE_STRING:
			o.Fields[index].Name = columnRules[value.Name].NameTable + "." + value.Name
			o.Fields[index].Type = enum.TYPE_STRING

		case enum.TYPE_UUID:
			_, err := uuid.Parse(value.Value)
			if err != nil {
				return fmt.Errorf("incorrect type field: %s", value.Name)
			}
			o.Fields[index].Name = columnRules[value.Name].NameTable + "." + value.Name
			o.Fields[index].Type = enum.TYPE_STRING

		case enum.TYPE_DATA:
			dateTime, err := helpers.ParseUnixTimeStampToString(value.Value)
			if err != nil {
				return err
			}
			o.Fields[index].Name = columnRules[value.Name].NameTable + "." + value.Name
			o.Fields[index].Value = dateTime
			o.Fields[index].Type = enum.TYPE_DATA

		case enum.TYPE_BOOL:
			o.Fields[index].Name = columnRules[value.Name].NameTable + "." + value.Name
			o.Fields[index].Type = enum.TYPE_BOOL

		default:
			return fmt.Errorf("incorrect filtering field: %s", value.Name)
		}
	}

	return nil
}

// CreateConditionsFilter
//
// The function takes a variable number of string arguments as input and returns a map of interface{} and an error as output.
//
// The function is a method attached to the Options struct and is used to create a set of filter conditions based on the fields in the struct.
// The method loops through each field in the struct and calls the conversionType function to convert the field value into the appropriate type for the database query.
// If there is an error during the conversion, the method returns a nil map and the error.
//
// The method then checks if any optional arguments have been passed in as part of the input and if so, adds the field name with the appropriate operator and condition to the filter map.
// If no optional arguments have been passed, the method uses the field name without any prefix. The map is returned as the first output value along with a nil error.
func (o *Options) CreateConditionsFilter() (map[string]interface{}, error) {
	mapConditionsFilter := make(map[string]interface{})

	for _, field := range o.Fields {
		addCondition, value, err := conversionType(field.Value, field.Type)
		if err != nil {
			return nil, err
		}
		if field.Operator == "ILIKE" {
			value = "%" + value.(string) + "%"
		}
		if field.Operator == "IN" {
			value = strings.Split(value.(string), ",")
		}
		mapConditionsFilter[fmt.Sprintf("%s %s ?%s", field.Name, field.Operator, addCondition)] = value
	}
	return mapConditionsFilter, nil
}

/*
This function is called conversionType and takes two arguments:

	value (string type)
	_type (enum type.ValidateType)

It is used to convert the value passed as the value argument to the corresponding data type defined in enum.ValidateType.

The function returns a tuple of three elements: string, interface{} and error.
The first element (string) is a string that can be added to an SQL query to perform data type conversion.
The second element (interface{}) is a converted data type value.
The third element (error) is an error that may occur during the conversion of a data type.

The function uses the switch statement to determine the data type specified in _type.
If _type is equal to enum.TYPE_DATA, the function returns the string "::timestamp" and the value as the converted value of the data type.

If _type is equal to enum.TYPE_INT, the function tries to convert the value to an integer data type using the strconv function.ParseInt().
If the conversion is successful, the function returns an empty string for the SQL statement and the converted value as the second element of the tuple.
If a conversion error has occurred, the function returns an empty string and the value as the second element of the tuple, as well as an error that occurred during the conversion as the third element of the tuple.

If _type is equal to enum.TYPE_STRING, the function returns an empty string for the SQL statement and the value as the converted value of the data type.

If _type does not match any of the defined data types, the function returns an empty string for the SQL statement, nil as the second element of the tuple, and the error "not type", which indicates that the data type is not defined.
*/
func conversionType(value string, _type enum.ValidateType) (string, interface{}, error) {
	switch _type {
	case enum.TYPE_DATA:
		return "::timestamp", value, nil

	case enum.TYPE_INT:
		valueInt, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return "", valueInt, nil
		} else {
			return "", nil, err
		}

	case enum.TYPE_STRING:
		return "", value, nil

	case enum.TYPE_BOOL:
		return "::boolean", value, nil

	default:
		return "", nil, errors.New("not type")
	}
}
