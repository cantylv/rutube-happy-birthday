package entity

type ErrorDetail struct {
	Error string `json:"error" valid:"-"`
}

type ResponseDetail struct {
	Detail string `json:"detail" valid:"-"`
}

type JwtTokenHeader struct {
	Exp string `json:"exp" valid:"-"`
}

type JwtTokenPayload struct {
	Id string `json:"id" valid:"-"`
}
