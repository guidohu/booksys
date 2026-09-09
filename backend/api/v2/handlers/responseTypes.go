package handlers

// Status is the envelope every API response carries. It reports whether the
// call succeeded and a human readable message.
type Status struct {
	OK  bool   `json:"ok"`
	Msg string `json:"msg"`
}

// DataResponse is the payload of a successful response. It is whatever the
// individual handler chose to return.
type DataResponse any

// Response is the body of every API response: a Status plus an optional
// payload.
type Response struct {
	Status
	Data DataResponse `json:"data,omitempty"`
}

// FailureResponse builds a Response that reports a failure with the given
// message.
func FailureResponse(msg string) Response {
	return Response{
		Status: Status{
			OK:  false,
			Msg: msg,
		},
	}
}
