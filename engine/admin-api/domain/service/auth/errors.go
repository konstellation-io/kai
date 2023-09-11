package auth

type UnauthorizedError struct {
	Product string
	Action  Action

	Err error
}

func (e UnauthorizedError) Error() string {
	msg := "you don't have authorization to " + e.Action.String()

	if e.Product != "" {
		msg += " in product " + e.Product
	}

	if e.Err != nil {
		msg += ": " + e.Err.Error()
	}

	return msg
}
