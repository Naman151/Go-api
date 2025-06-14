package storage

import "github.com/Naman151/Go-api/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudentsList() ([]types.Student, error)
	// UpdateStudentById(id int64) (types.Student, error)
	DeleteStudentById(id int64) error
}
