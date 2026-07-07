package model

type SourceService string

const (
	ServiceUser      SourceService = "user-service"
	ServiceCourse    SourceService = "course-service"
	ServiceHomework  SourceService = "homework-service"
	ServiceAnalytics SourceService = "analytics-service"
)