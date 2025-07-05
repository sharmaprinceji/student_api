package storage

 import "github.com/sharmaprinceji/student-api/internal/types"

type Storage interface {
	CreateStudent(name string,email string,age int,city string) (int64,error)
	GetStudentById(id int64) (types.Student, error)
}
