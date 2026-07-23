package dto

type CreateCourseRequest struct {
	Title       string `json:"title" binding:"required,max=50"`
	Description string `json:"description" binding:"max=200"`
}

type UpdateCourseRequest struct {
	Title       *string `json:"title" binding:"omitempty,max=50"`
	Description *string `json:"description" binding:"omitempty,max=200"`
}
