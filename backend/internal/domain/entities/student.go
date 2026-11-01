package entities

type Student struct {
	ID       int64  `json:"id"`
	CourseID int64  `json:"courseId"`
	FullName string `json:"fullName"`
}
