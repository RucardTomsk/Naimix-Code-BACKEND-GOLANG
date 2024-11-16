package base

import "github.com/google/uuid"

type Blame string

const (
	BlameUser     Blame = "User"
	BlamePostgres Blame = "Postgres"
	BlameS3       Blame = "S3"
	BlameServer   Blame = "Server"
	BlameUnknown  Blame = "Unknown"
	BlameMail     Blame = "Mail"
)

// ResponseOK is a base OK response from server.
type ResponseOK struct {
	Status string `json:"status" example:"OK"`
}

// ResponseOKWithID is a base OK response from server with additional ID in answer.
type ResponseOKWithID struct {
	Status string    `json:"status" example:"OK"`
	ID     uuid.UUID `json:"ID" example:"12345678-1234-1234-1234-000000000000"`
}

// ResponseFailure is a general error response from server.
type ResponseFailure struct {
	Status  string `json:"status" example:"Error"`
	Blame   Blame  `json:"blame" example:"Guilty System"`
	Message string `json:"message" example:"error occurred"`
}

type ResponseOKWithContent struct {
	Status  string `json:"status" example:"OK"`
	Content string `json:"content"`
}

type ResponseOKWithListContent struct {
	Status   string   `json:"status" example:"OK"`
	Contents []string `json:"content"`
}
