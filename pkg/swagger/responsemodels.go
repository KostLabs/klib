package swagger

type StatusOKResponse struct {
	Status struct {
		Code    int    `json:"code" example:"200"`
		Enum    string `json:"enum" example:"OK"`
		Message string `json:"message" example:"request successful"`
	} `json:"status"`
}

type StatusCreatedResponse struct {
	Status struct {
		Code    int    `json:"code" example:"201"`
		Enum    string `json:"enum" example:"CREATED"`
		Message string `json:"message" example:"resource created successfully"`
	} `json:"status"`
}

type StatusAcceptedResponse struct {
	Status struct {
		Code    int    `json:"code" example:"202"`
		Enum    string `json:"enum" example:"ACCEPTED"`
		Message string `json:"message" example:"request accepted for processing"`
	} `json:"status"`
}

type StatusNoContentResponse struct {
	Status struct {
		Code    int    `json:"code" example:"204"`
		Enum    string `json:"enum" example:"NO_CONTENT"`
		Message string `json:"message" example:"request successful, no content to return"`
	} `json:"status"`
}

// Location comes as header in the response, not in the body
type StatusMovedPermanentlyResponse struct {
	Status struct {
		Code    int    `json:"code" example:"301"`
		Enum    string `json:"enum" example:"MOVED_PERMANENTLY"`
		Message string `json:"message" example:"resource moved permanently"`
	} `json:"status"`
}

type StatusFoundResponse struct {
	Status struct {
		Code    int    `json:"code" example:"302"`
		Enum    string `json:"enum" example:"FOUND"`
		Message string `json:"message" example:"resource found at a different location"`
	} `json:"status"`
}

type StatusNotModifiedResponse struct {
	Status struct {
		Code    int    `json:"code" example:"304"`
		Enum    string `json:"enum" example:"NOT_MODIFIED"`
		Message string `json:"message" example:"resource not modified since last request"`
	} `json:"status"`
}

type StatusTemporaryRedirectResponse struct {
	Status struct {
		Code    int    `json:"code" example:"307"`
		Enum    string `json:"enum" example:"TEMPORARY_REDIRECT"`
		Message string `json:"message" example:"resource temporarily redirected to a different location"`
	} `json:"status"`
}

type StatusPermanentRedirectResponse struct {
	Status struct {
		Code    int    `json:"code" example:"308"`
		Enum    string `json:"enum" example:"PERMANENT_REDIRECT"`
		Message string `json:"message" example:"resource permanently redirected to a different location"`
	} `json:"status"`
}

type StatusBadRequestResponse struct {
	Status struct {
		Code  int    `json:"code" example:"400"`
		Enum  string `json:"enum" example:"BAD_REQUEST"`
		Error string `json:"error" example:"invalid request parameters"`
	} `json:"status"`
}

type StatusUnauthorizedResponse struct {
	Status struct {
		Code  int    `json:"code" example:"401"`
		Enum  string `json:"enum" example:"UNAUTHORIZED"`
		Error string `json:"error" example:"authentication required"`
	} `json:"status"`
}

type StatusForbiddenResponse struct {
	Status struct {
		Code  int    `json:"code" example:"403"`
		Enum  string `json:"enum" example:"FORBIDDEN"`
		Error string `json:"error" example:"access to the resource is forbidden"`
	} `json:"status"`
}

type StatusNotFoundResponse struct {
	Status struct {
		Code  int    `json:"code" example:"404"`
		Enum  string `json:"enum" example:"NOT_FOUND"`
		Error string `json:"error" example:"resource not found"`
	} `json:"status"`
}

type StatusMethodNotAllowedResponse struct {
	Status struct {
		Code  int    `json:"code" example:"405"`
		Enum  string `json:"enum" example:"METHOD_NOT_ALLOWED"`
		Error string `json:"error" example:"HTTP method not allowed for the resource"`
	} `json:"status"`
}

type StatusConflictResponse struct {
	Status struct {
		Code  int    `json:"code" example:"409"`
		Enum  string `json:"enum" example:"CONFLICT"`
		Error string `json:"error" example:"resource conflict occurred"`
	} `json:"status"`
}

type StatusInternalServerErrorResponse struct {
	Status struct {
		Code  int    `json:"code" example:"500"`
		Enum  string `json:"enum" example:"INTERNAL_SERVER_ERROR"`
		Error string `json:"error" example:"internal server error occurred"`
	} `json:"status"`
}

type StatusNotImplementedResponse struct {
	Status struct {
		Code  int    `json:"code" example:"501"`
		Enum  string `json:"enum" example:"NOT_IMPLEMENTED"`
		Error string `json:"error" example:"requested functionality not implemented"`
	} `json:"status"`
}

type StatusBadGatewayResponse struct {
	Status struct {
		Code  int    `json:"code" example:"502"`
		Enum  string `json:"enum" example:"BAD_GATEWAY"`
		Error string `json:"error" example:"upstream service unavailable"`
	} `json:"status"`
}

type StatusServiceUnavailableResponse struct {
	Status struct {
		Code  int    `json:"code" example:"503"`
		Enum  string `json:"enum" example:"SERVICE_UNAVAILABLE"`
		Error string `json:"error" example:"service temporarily unavailable"`
	} `json:"status"`
}

type StatusGatewayTimeoutResponse struct {
	Status struct {
		Code  int    `json:"code" example:"504"`
		Enum  string `json:"enum" example:"GATEWAY_TIMEOUT"`
		Error string `json:"error" example:"upstream service timed out"`
	} `json:"status"`
}
