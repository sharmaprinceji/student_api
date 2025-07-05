package storage

 import "github.com/sharmaprinceji/student-api/internal/types"

type Storage interface {
	CreateStudent(name string,email string,age int,city string) (int64,error)
	GetStudentById(id int64) (types.Student, error)
	GetAllStudents() ([]types.Student, error)
	UpdateStudentById(id int64, name string, email string, age int, city string) (int64, error)
	DeleteStudent(id int64) (int64, error)
}
