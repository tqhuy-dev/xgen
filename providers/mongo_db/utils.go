package mongo_db

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ObjectIDToString converts a MongoDB ObjectID to a hex string
// Returns the hex string representation of the ObjectID
func ObjectIDToString(objectID bson.ObjectID) string {
	return objectID.Hex()
}

// StringToObjectID converts a hex string to a MongoDB ObjectID
// Returns the ObjectID and an error if the string is not a valid ObjectID hex string
func StringToObjectID(str string) (bson.ObjectID, error) {
	objectID, err := bson.ObjectIDFromHex(str)
	if err != nil {
		return bson.ObjectID{}, fmt.Errorf("invalid ObjectID hex string: %w", err)
	}
	return objectID, nil
}

// AnyToObjectID converts any type to a MongoDB ObjectID
// Supports the following input types:
//   - bson.ObjectID: returns as-is
//   - string: converts hex string to ObjectID
//   - []byte: converts byte slice to ObjectID (must be 12 bytes)
//   - [12]byte: converts byte array to ObjectID
//
// Returns the ObjectID and an error if conversion fails
func AnyToObjectID(value any) (bson.ObjectID, error) {
	if value == nil {
		return bson.ObjectID{}, fmt.Errorf("cannot convert nil to ObjectID")
	}

	switch v := value.(type) {
	case bson.ObjectID:
		// Already an ObjectID, return as-is
		return v, nil

	case string:
		// Try to parse as hex string
		return StringToObjectID(v)

	case []byte:
		// Try to create ObjectID from byte slice
		if len(v) != 12 {
			return bson.ObjectID{}, fmt.Errorf("byte slice must be exactly 12 bytes, got %d", len(v))
		}
		var arr [12]byte
		copy(arr[:], v)
		return bson.ObjectID(arr), nil

	case [12]byte:
		// Convert byte array to ObjectID
		return bson.ObjectID(v), nil

	default:
		return bson.ObjectID{}, fmt.Errorf("unsupported type for ObjectID conversion: %T", value)
	}
}

// MustStringToObjectID converts a hex string to a MongoDB ObjectID
// Panics if the string is not a valid ObjectID hex string
func MustStringToObjectID(str string) bson.ObjectID {
	objectID, err := StringToObjectID(str)
	if err != nil {
		panic(err)
	}
	return objectID
}

// MustAnyToObjectID converts any type to a MongoDB ObjectID
// Panics if the conversion fails
func MustAnyToObjectID(value any) bson.ObjectID {
	objectID, err := AnyToObjectID(value)
	if err != nil {
		panic(err)
	}
	return objectID
}

// IsValidObjectIDHex checks if a string is a valid ObjectID hex string
func IsValidObjectIDHex(str string) bool {
	_, err := bson.ObjectIDFromHex(str)
	return err == nil
}
