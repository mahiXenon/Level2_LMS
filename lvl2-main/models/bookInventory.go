package models

type BookInventory struct {
	ISBN            string `json:"isbn" binding:"required" gorm:"primary_key"`
	LibraryId       uint   `json:"library_id" binding:"required" gorm:"primary_key"`
	Title           string `json:"title" binding:"required"`
	Author          string `json:"author" binding:"required"`
	Publisher       string `json:"publisher" binding:"required"`
	Version         string `json:"version" binding:"required"`
	TotalCopies     int    `json:"total_copies" binding:"required"`
	AvailableCopies int    `json:"available_copies" binding:"required"`
}
type InputBook struct {
	ISBN        string `json:"isbn" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Author      string `json:"author" binding:"required"`
	Publisher   string `json:"publisher" binding:"required"`
	Version     string `json:"version" binding:"required"`
	TotalCopies int    `json:"total_copies" binding:"required"`
}

type UpdateBookDetails struct {
	ISBN          string `json:"isbn" binding:"required"`
	ADD           int    `json:"add" default:"0"`
	DecreaseCount int    `json:"decrease_count" default:"0"`
}
