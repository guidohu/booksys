package handlers

type Status struct {
	OK  bool   `json:"ok"`
	Msg string `json:"msg"`
}

type DataResponse any

type Response struct {
	Status
	Data DataResponse `json:"data,omitempty"`
}

func FailureResponse(msg string) Response {
	return Response{
		Status: Status{
			OK:  false,
			Msg: msg,
		},
	}
}
