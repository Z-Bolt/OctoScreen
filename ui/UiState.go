package ui

type UiState struct {
	placeholder string
}

func (this UiState) String() string {
	return this.placeholder
}

var (
	Uninitialized	= UiState{""}
	connecting		= UiState{"connecting"}
	idle			= UiState{"idle"}
	printing		= UiState{"printing"}
)
