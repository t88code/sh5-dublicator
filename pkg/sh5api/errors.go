package sh5api

type ProcessingError struct {
	error
	*ErrorReply
}

type TransportError struct {
	error
	HttpCode int
}

type ErrorReply struct {
	ErrorCode  int    `json:"errorCode"`
	ErrMessage string `json:"errMessage"`
	Version    string `json:"Version"`
	UserName   string `json:"UserName"`
	ActionName string `json:"actionName"`
	ActionType string `json:"actionType"`
	ErrorInfo  struct {
		ModuleId   int    `json:"moduleId"`
		ModuleName string `json:"moduleName"`
		ErrorId    int    `json:"errorId"`
	} `json:"errorInfo"`
}
