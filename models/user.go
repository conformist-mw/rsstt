package models

type SeenItem struct {
	ID     uint `gorm:"primarykey"`
	ItemID uint `gorm:"uniqueIndex" json:"item_id"`
}
