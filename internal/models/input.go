package models

// UserInput представляет структуру одного пользователя
// во входящем XML-документе
type UserInput struct {
	UserID string `xml:"id,attr"`
	Name   string `xml:"name"`
	Email  string `xml:"email"`
	Age    int    `xml:"age"`
}

// UsersInput представляет список пользователей в XML-документе
type UsersInput struct {
	Users []UserInput `xml:"user"`
}
