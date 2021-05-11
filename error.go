package planet

import "fmt"

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("[%s] - %s", e.Code, e.Message)
}

var (
	ErrUnknown = Error{
		Code:    "PNT-001",
		Message: "There was an unknown error processing your request.",
	}
	ErrDecodeRequest = Error{
		Code:    "PNT-002",
		Message: "There was unable to decode your request.",
	}
	ErrPlanetIDEmpty = Error{
		Code:    "PNT-003",
		Message: "The PlanetID must not be empty.",
	}
	ErrPlanetAlreadyProcessed = Error{
		Code:    "PNT-004",
		Message: "The planet has already been processed.",
	}
)
