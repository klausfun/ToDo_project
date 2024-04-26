package todo

type User struct {
	Id       int    `json:"-" db:"id"`                   // (2)
	Name     string `json:"name" binding:"required"`     // (1)
	Username string `json:"username" binding:"required"` // (1)
	Password string `json:"password" binding:"required"` // (1)
}

// (1)
// при регистрации нам необходимо получать от пользователя поля: Name, Username, Password
// теги binding:"required" - валидируют наличие данных полей в теле запроса и являются реализацией фреймворка gin

// (2)
// для того чтобы метод Get работал у структуры необходимо прописать теги db с названием поля из базы
