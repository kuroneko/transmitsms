package transmitsms

const (
	ErrApiOverLimit = 429
)

type ApiError struct {
	Message      string
	HttpCode     int
	ResponseBody string
}

func (err *ApiError) Error() string {
	return err.Message
}
