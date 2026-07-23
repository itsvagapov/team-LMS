package dto

type CreateLessonRequest struct {
	Title   string `json:"title" binding:"required,max=50"`
	Content string `json:"content" binding:"required,max=2000"`
}
