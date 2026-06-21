package enums

type HTTPStatus string

func (s HTTPStatus) String() string {
	return string(s)
}

const (
	StatusOk                  HTTPStatus = "OK"
	StatusCreated             HTTPStatus = "CREATED"
	StatusAccepted            HTTPStatus = "ACCEPTED"
	StatusNoContent           HTTPStatus = "NO_CONTENT"
	StatusMovedPermanently    HTTPStatus = "MOVED_PERMANENTLY"
	StatusFound               HTTPStatus = "FOUND"
	StatusNotModified         HTTPStatus = "NOT_MODIFIED"
	StatusTemporaryRedirect   HTTPStatus = "TEMPORARY_REDIRECT"
	StatusPermanentRedirect   HTTPStatus = "PERMANENT_REDIRECT"
	StatusBadRequest          HTTPStatus = "BAD_REQUEST"
	StatusUnauthorized        HTTPStatus = "UNAUTHORIZED"
	StatusForbidden           HTTPStatus = "FORBIDDEN"
	StatusNotFound            HTTPStatus = "NOT_FOUND"
	StatusMethodNotAllowed    HTTPStatus = "METHOD_NOT_ALLOWED"
	StatusConflict            HTTPStatus = "CONFLICT"
	StatusInternalServerError HTTPStatus = "INTERNAL_SERVER_ERROR"
	StatusNotImplemented      HTTPStatus = "NOT_IMPLEMENTED"
	StatusBadGateway          HTTPStatus = "BAD_GATEWAY"
	StatusServiceUnavailable  HTTPStatus = "SERVICE_UNAVAILABLE"
	StatusGatewayTimeout      HTTPStatus = "GATEWAY_TIMEOUT"
)
