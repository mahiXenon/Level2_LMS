package models

import "time"

type RequestInput struct {
	ISBN        string `json:"isbn" binding:"required"`
	LibraryId   uint   `json:"library_id" binding:"required"`
	RequestType string `json:"request_type" binding:"required"`
}

type RequestEvent struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ISBN        string    `json:"isbn" binding:"required"`
	UserId      uint      `json:"user_id" binding:"required"`
	LibraryId   uint      `json:"library_id" binding:"required"`
	RequestDate time.Time `json:"request_date"`
	ApproveDate time.Time `json:"approve_date"`
	ApproverId  uint      `json:"approver_id"`
	RequestType string    `json:"request_type" default:"borrow"`
}

type ListRequest struct {
	ID          uint      `json:"id"`
	ISBN        string    `json:"isbn"`
	UserId      uint      `json:"user_id"`
	RequestType string    `json:"request_type"`
	RequestDate time.Time `json:"request_date"`
}
