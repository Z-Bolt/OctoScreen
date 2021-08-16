package dataModels

type FilamentManagerSelection struct {
	// Unknown?
	ClientId string `json: "client_id"`

	// Currently selected spool
	Spool FilamentManagerSpool `json: "spool"`

	// Tool the spool is selected on
	Tool int `json: "tool"`
}

type FilamentManagerSelections struct {
	Selections []*FilamentManagerSelection `json: "selections"`
}