package enum

type ValidateType int

const (
	TYPE_INT ValidateType = iota
	TYPE_STRING
	TYPE_DATA
	TYPE_BOOL
	TYPE_UUID
)

func (s ValidateType) String() string {
	return [...]string{"int", "string", "datetime", "boolean", "uuid"}[s]
}
