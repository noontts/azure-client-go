package http

// PaymentResponse represents the structure of the payment API response
// matching the provided JSON example.
type PaymentResponse struct {
	StringExample string      `json:"string_example"`
	IntExample    int         `json:"int_example"`
	FloatExample  float64     `json:"float_example"`
	BooleanTrue   bool        `json:"boolean_true"`
	BooleanFalse  bool        `json:"boolean_false"`
	NullExample   interface{} `json:"null_example"`
	ArrayExample  []any       `json:"array_example"`
	ObjectExample struct {
		NestedString string `json:"nested_string"`
		NestedNumber int    `json:"nested_number"`
		NestedArray  []int  `json:"nested_array"`
	} `json:"object_example"`
}
