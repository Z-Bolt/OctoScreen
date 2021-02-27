package octoprintApis

import (
	// "bytes"
	// "encoding/json"
	// "io"
)


// https://docs.octoprint.org/en/master/api/job.html
const JobApiUri = "/api/job"


var JobToolErrors = StatusMapping {
	409: "Printer is not operational or the current print job state does not match the preconditions for the command.",
}
