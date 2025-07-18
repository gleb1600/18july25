package storage

type taskStatus string

const (
	TaskStatusCreated                 taskStatus = "Task Created"
	TaskStatusCompletedSuccessfully   taskStatus = "Task Completed Successfully"
	TaskStatusCompletedUnsuccessfully taskStatus = "Task Completed Unsuccessfully"
)

type task struct {
	ID     string
	Status taskStatus
}
