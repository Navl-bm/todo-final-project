// Пакет database используется для взаимодействия с базой данных приложения

package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Navl-bm/todo-final-project/config"
	"github.com/Navl-bm/todo-final-project/errors"
	"github.com/Navl-bm/todo-final-project/models"
	_ "modernc.org/sqlite"
)

// Константы для определения типа фильтрации задач
const (
	SearchTypeNil            = 0
	SearchTypeTitleOrComment = 1
	SearchTypeDate           = 2
)

// DB метод для взаимодействия с базой данных
type DB struct {
	db  *sql.DB
	cfg *config.Config
}

// NewDataBase создает новый экземпляр базы данных
func NewDataBase(cfg *config.Config) DB {
	return DB{db: nil, cfg: cfg}
}

// Init выполняет подключение к базе данных
// В случае остутствия файла создает его и таблицу для хранения данных
func (db *DB) Init() error {
	_, err := os.Stat(db.cfg.DBName)
	var install bool
	if err != nil {
		install = true
	}

	db.db, err = sql.Open("sqlite", db.cfg.DBName)
	if err != nil {
		return err
	}

	if install {
		_, err := db.db.Exec(`
		CREATE TABLE scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(128) NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat VARCHAR(128) NOT NULL DEFAULT ""
		)
		`)

		if err != nil {
			return err
		}

		_, err = db.db.Exec(`CREATE INDEX idx_scheduler_date ON scheduler (date)`)
		if err != nil {
			return err
		}
	}

	return nil
}

// AddTask выполняет добавление новой задачи с переданными параметрами
func (db *DB) AddTask(task *models.Task) (int64, error) {
	var id int64

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := db.db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

// GetTask получает список задач по заданному фильтру
// задачи отсортированы по возрастанию даты, начиная с ближайшей
func (db *DB) GetTasks(limit int, search string, searchType int) (*models.TasksList, error) {
	var query string

	switch searchType {
	case SearchTypeNil:
		query = `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit `
	case SearchTypeTitleOrComment:
		query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit`
	case SearchTypeDate:
		query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date LIMIT :limit `
	}

	rows, err := db.db.Query(query,
		sql.Named("limit", limit),
		sql.Named("search", fmt.Sprintf("%c%s%c", '%', search, '%')),
		sql.Named("date", search),
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var taskList []*models.Task

	for rows.Next() {
		task := models.Task{}
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		taskList = append(taskList, &task)
	}

	if taskList == nil {
		taskList = []*models.Task{}
	}

	return &models.TasksList{Tasks: taskList}, nil
}

// GetTask возвращает задачу по заданному id
func (db *DB) GetTask(id int) (*models.Task, error) {

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id`

	row := db.db.QueryRow(query,
		sql.Named("id", id),
	)
	task := models.Task{}

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &task, nil
}

// UpdateTask обновляет данные задачи по заданному id
func (db *DB) UpdateTask(task *models.Task) error {

	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`

	res, err := db.db.Exec(query,
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)

	if err != nil {
		log.Println(err)
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.ErrTaskNotFound
	}

	return nil
}

// DeleteTask удаляет задачу по заданному id
func (db *DB) DeleteTask(id int) error {

	query := `DELETE FROM scheduler WHERE id = :id`

	res, err := db.db.Exec(query,
		sql.Named("id", id),
	)

	if err != nil {
		log.Println(err)
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.ErrTaskNotFound
	}

	return nil
}
