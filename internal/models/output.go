package models

// UserOutput представляет структуру пользователя в JSON-ответе
type UserOutput struct {
	UserID   string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	AgeGroup string `json:"age_group"`
}

// UsersOutput представляет обертку для списка пользователей в JSON-ответе
type UsersOutput struct {
	Users []UserOutput `json:"users"`
}
