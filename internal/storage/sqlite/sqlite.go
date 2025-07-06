package sqlite

import (
	"database/sql"
	 "fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sharmaprinceji/student-api/internal/config"
	 "github.com/sharmaprinceji/student-api/internal/types"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	
	_,err=db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		age INTEGER,
		email TEXT UNIQUE,
		city Text
	)`)

	if err != nil {
		return nil, err
	}	

	return &Sqlite{
		Db:db,
	},nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int, city string) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students(name, email, age, city) VALUES(?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name, email, age, city)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id = ? LIMIT 2")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.ID, &student.Name, &student.Age, &student.Email, &student.City)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %s",fmt.Sprint(id)) // No student found with the given ID
		}

		return types.Student{}, fmt.Errorf("failed to query student: %w", err) // Other error occurred while querying
	}

	return student, nil
}

func (s *Sqlite) GetAllStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student
		err = rows.Scan(&student.ID, &student.Name, &student.Age, &student.Email, &student.City)
		if err != nil {
			return nil, fmt.Errorf("failed to scan student: %w", err)
		}
		students = append(students, student)
	}

	return students, nil
}

func (s *Sqlite) GetStudentsPaginated(page int, limit int) ([]types.Student, int, error) {
	offset := (page - 1) * limit

	// Get paginated results
	rows, err := s.Db.Query("SELECT id, name, age, email, city FROM students LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var students []types.Student
	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.ID, &student.Name, &student.Age, &student.Email, &student.City)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan student: %w", err)
		}
		students = append(students, student)
	}

	// Get total count
	var totalCount int
	err = s.Db.QueryRow("SELECT COUNT(*) FROM students").Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count students: %w", err)
	}

	return students, totalCount, nil
}

func (s *Sqlite) UpdateStudentById(id int64, name string, email string, age int, city string) (int64, error) {
	stmt, err := s.Db.Prepare("UPDATE students SET name = ?, email = ?, age = ?, city = ? WHERE id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name, email, age, city, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (s *Sqlite) DeleteStudent(id int64) (int64, error) {
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}