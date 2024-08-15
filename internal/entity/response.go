package entity

type ErrorDetail struct {
	Error string `json:"error" valid:"-"`
}

type ResponseDetail struct {
	Detail string `json:"detail" valid:"-"`
}
