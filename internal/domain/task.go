package domain

type Task struct {
	ID     string `gorm:"primary_key" json:"ID"`
	Name   string `gorm:"not null" json:"Name"`
	Status bool   `gorm:"not null" json:"Status"`
}

type TaskUpdate struct {
	Name string `json:"name"`
}

type TaskCreate struct {
	Name string `json:"name"`
}
