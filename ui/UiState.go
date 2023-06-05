package ui

type UiState struct {
	placeholder string
}

func (this UiState) String() string {
	return this.placeholder
}

var (
	Uninitialized	= UiState{""}
	Connecting		= UiState{"connecting"}
	Idle			= UiState{"idle"}
	Printing		= UiState{"printing"}
)
