package publictypes

import "net/http"

type ControllerType func(w http.ResponseWriter, r *http.Request)
