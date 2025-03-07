package models

import "time"

type IssueBook struct {
	ISBN               string    `json:"isbn" binding:"required"`
	LibraryId          uint      `json:"library_id" binding:"required"`
	UserID             uint      `json:"user_id" binding:"required"`
	IssueApproverId    uint      `json:"issue_approver_id" binding:"required"`
	IssueStatus        string    `json:"issue_status" default:"pending"`
	IssueDate          time.Time `json:"issue_date"`
	ExpectedReturnDate time.Time `json:"expected_return_date"`
	ReturnDate         time.Time `json:"return_date"`
	ReturnApproverId   uint      `json:"return_approver_id"`
}

type ReaderRequest struct {
	ID uint `json:"id" binding:"required"`
}
