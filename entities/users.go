package entities

type Users struct {
	UserId    int64
	Username  string
	Password  string
	FirstName string
	LastName  string
	UserType  string
	IsEnabled string
}

type UsersLogin struct {
	UserId    int64  `json:"userId"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	UserType  string `json:"userType"`
}
