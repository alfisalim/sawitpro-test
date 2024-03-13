// This file contains types that are used in the repository layer.
package repository

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type Profile struct {
	UserId    int64  `gorm:"column:user_id;PRIMARY_KEY;AUTO_INCREMENT"`
	FullName  string `gorm:"column:full_name"`
	Password  string `gorm:"column:password"`
	Phone     string `gorm:"column:phone"`
	Status    int64  `gorm:"column:status"`
	CreatedAt string `gorm:"column:created_at"`
	UpdatedAt string `gorm:"column:updated_at"`
}

func (Profile) TableName() string {
	return "users"
}

type LoginModel struct {
	LoginId  int64  `gorm:"column:login_id;PRIMARY_KEY;AUTO_INCREMENT"`
	UserId   int64  `gorm:"column:user_id"`
	Ip       string `gorm:"column:ip"`
	Token    string `gorm:"column:token"`
	Expires  string `gorm:"expires"`
	Requests int64  `gorm:"requests"`
}

func (LoginModel) TableName() string {
	return "login"
}
