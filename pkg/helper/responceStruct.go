package helper

type SingleRespStruct struct {
	Error string `json:"error"`
}

// to responce with avoid unwanted details
type UserRespStrcut struct {
	ID        uint   `json:"id" copier:"must"`
	FirstName string `json:"first_name" copier:"must"`
	LastName  string `json:"last_name" copier:"must"`
	Age       uint   `json:"age" copier:"must"`
	Email     string `json:"email" copier:"must"`
	Phone     string `json:"phone" copier:"must"`
}
