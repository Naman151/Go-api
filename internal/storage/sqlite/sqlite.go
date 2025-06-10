package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/Naman151/Go-api/internal/config"
	"github.com/Naman151/Go-api/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

var CreateTable = `CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		age INTEGER
		)`

var InsertStudent = `INSERT INTO students (name, email, age) VALUES (?, ?, ?)`
var SelectStudentById = `SELECT id, name, email, age FROM students where id = ? LIMIT 1`
var SelectStudents = `SELECT id, name, email, age FROM students`

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(CreateTable)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare(InsertStudent)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare(SelectStudentById)
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("No Student with Id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}
	return student, nil
}

func (s *Sqlite) GetStudentsList() ([]types.Student, error) {
	stmt, err := s.Db.Prepare(SelectStudents)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var students []types.Student

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var student types.Student

		rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}
	return students, nil
}

// func (s *Sqlite) UpdateStudentById() ([]types.Student, error) {
// 	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer stmt.Close()

// 	var students []types.Student

// 	rows, err := stmt.Query()
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer rows.Close()

// 	for rows.Next() {
// 		var student types.Student

// 		rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
// 		if err != nil {
// 			return nil, err
// 		}

// 		students = append(students, student)
// 	}
// 	return students, nil
// }

func (s *Sqlite) DeleteStudentById(id int64) error {
	stmt, err := s.Db.Prepare("DELETE * FROM students where id = ? LIMIT 1")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
