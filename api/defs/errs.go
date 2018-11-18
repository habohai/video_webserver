package defs

// Err 错误封装
type Err struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}

// ErrResponse 错误返回结构
type ErrResponse struct {
	HTTPSC int
	Error  Err
}

var (
	// ErrorRequestBodyParseFailed Request body is not correct.
	ErrorRequestBodyParseFailed = ErrResponse{
		HTTPSC: 400,
		Error: Err{
			Error:     "Request body is not correct.",
			ErrorCode: "001",
		},
	}

	// ErrorNotAuthUser User authentication failed.
	ErrorNotAuthUser = ErrResponse{
		HTTPSC: 401,
		Error: Err{
			Error:     "User authentication failed.",
			ErrorCode: "002",
		},
	}

	// ErrorDBError DB ops failed.
	ErrorDBError = ErrResponse{
		HTTPSC: 500,
		Error: Err{
			Error:     "DB ops failed.",
			ErrorCode: "003",
		},
	}

	// ErrorInternalFaults Internal service failed.
	ErrorInternalFaults = ErrResponse{
		HTTPSC: 500,
		Error: Err{
			Error:     "Internal service failed.",
			ErrorCode: "004",
		},
	}
)
