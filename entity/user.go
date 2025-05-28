package entity

type User struct {
	ID    uint64 `gorm:"primaryKey,autoIncrement,type:bigint" json:"id"`
	Name  string `gorm:"type:varchar(255)" json:"name"`
	Email string `gorm:"type:varchar(255)" json:"email"`
	Age   uint64 `gorm:"type:bigint" json:"age"`
}
