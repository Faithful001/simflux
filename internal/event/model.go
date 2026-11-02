package event

type Event struct {
	ID   int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(255);not null;unique" json:"name"`
}

func (Event) TableName() string {
	return "events"
}