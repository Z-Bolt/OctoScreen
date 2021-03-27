package octoprintApis

import (
	// "bytes"
	// "encoding/json"
	// "io"

	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


const ConnectionApiUri = "/api/connection"


var ConnectionErrors = StatusMapping {
	400: "The selected port or baudrate for a connect command are not part of the available option",
}
