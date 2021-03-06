package types

const (
	ValidateTag = "validate"

	StringMinLenKey string = "min_len"
	StringMaxLenKey string = "max_len"

	NumberMinKey string = "min"
	NumberMaxKey string = "max"

	ArrayMinItemsKey string = "min_items"
	ArrayMaxItemsKey string = "max_items"
	ArrayItemKey     string = "item"

	PointerNullableKey string = "nullable"
	PointerNotNullKey  string = "not_null"

	InterfaceFuncKey string = "func"

	StructFuncKey string = "func"

	MapMinItemsKey string = "min_items"
	MapMaxItemsKey string = "max_items"
	MapKeyKey      string = "key"
	MapValueKey    string = "value"
)
