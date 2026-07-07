package model

type UserActivityEventType string

const (
	EventUserRegistered UserActivityEventType = "user.registered"
	EventUserLoggedIn   UserActivityEventType = "user.logged_in"
	EventRoleChanged    UserActivityEventType = "user.role_changed"

	EventStudentEnrolled    UserActivityEventType = "student.enrolled"
	EventSubmissionCreated  UserActivityEventType = "submission.created"
	EventSubmissionReviewed UserActivityEventType = "submission.reviewed"
)