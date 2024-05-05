package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
	"github.com/vova4o/go_final_project/internal/config"
	"github.com/vova4o/go_final_project/internal/models"
)

// InitDB создаёт базу данных, если она не существует, и возвращает объект sql.DB для работы с ней.
func InitDB() (*sql.DB, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	file := config.DBPath()

	dbFile := filepath.Join(currentDir, file)
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	log.Printf("Database file: %s\n", dbFile)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if install {
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat TEXT(128)
		);`)
		if err != nil {
			return nil, err
		}

		_, err = db.Exec(`CREATE INDEX IF NOT EXISTS indexdate ON scheduler (date)`)
		if err != nil {
			log.Println("Не создан индекс", err)
		}
	}

	return db, nil
}

// CloseDB закрывает соединение с базой данных.
func CloseDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

// AddTask добавляет задачу в базу данных. Возвращает идентификатор задачи.
// исходные данные: дата, заголовок, комментарий, правило повторения.
func AddTask(db *sql.DB, date string, title string, comment string, repeat string) (int64, error) {
	result, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)", date, title, comment, repeat)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// FindTask ищет задачу по идентификатору ID. Возвращает задачу или ошибку.
func FindTask(db *sql.DB, search string) (models.DBTask, error) {
	task := models.DBTask{}
	if search == "" {
		return task, errors.New("не указан id задачи")
	}

	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	rows, err := db.Query(query, search)
	if err != nil {
		return task, err
	}
	defer rows.Close()

	if !rows.Next() {
		return task, errors.New("задача не найдена")
	}

	err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, err
	}

	return task, nil
}

// UpdateTask обновляет задачу в базе данных. Возвращает ошибку.
func UpdateTask(db *sql.DB, task models.DBTask) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	_, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return errors.New("задача не найдена")
	}

	return nil
}

// Tasks возвращает список задач из базы данных. Возвращает список задач или ошибку.
func Tasks(db *sql.DB, offset int) ([]models.DBTask, error) {
	query := fmt.Sprintf("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 10 OFFSET %d", offset)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.DBTask
	for rows.Next() {
		var t models.DBTask
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func SearchTasks(db *sql.DB, search string) ([]models.DBTask, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ?"
	rows, err := db.Query(query, "%"+search+"%", "%"+search+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.DBTask
	for rows.Next() {
		var t models.DBTask
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func TasksByDate(db *sql.DB, date string) ([]models.DBTask, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ?"
	rows, err := db.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.DBTask
	for rows.Next() {
		var t models.DBTask
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

// DoneTask помечает задачу как выполненную. Возвращает ошибку. Если задача повторяющаяся, то создаёт новую задачу на следующую дату.
func DoneTask(db *sql.DB, id string) error {
	var taskWeDeleting models.DBTask
	err := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&taskWeDeleting.ID, &taskWeDeleting.Date, &taskWeDeleting.Title, &taskWeDeleting.Comment, &taskWeDeleting.Repeat)
	if err != nil {
		return errors.New("задача не найдена")
	}

	if taskWeDeleting.Repeat == "" {
		_, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			return errors.New("задача не найдена")
		}
	} else {
		taskWeDeleting.Date, err = NextDate(time.Now(), taskWeDeleting.Date, taskWeDeleting.Repeat)
		if err != nil {
			return err
		}
		err = UpdateTask(db, taskWeDeleting)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteTask удаляет задачу из базы данных. Возвращает ошибку.
func DeleteTask(db *sql.DB, id string) error {
	var exists bool
	err := db.QueryRow("SELECT exists(SELECT 1 FROM scheduler WHERE id=?)", id).Scan(&exists)
	if err != nil || !exists {
		return fmt.Errorf("задача не найдена")
	}

	_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

// NextDate returns the next date of the task
// с параметрами:
// now — время от которого ищется ближайшая дата;
// date — исходное время в формате 20060102, от которого начинается отсчёт повторений;
// repeat — правило повторения в описанном выше формате.
func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", nil
	}

	t, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	switch repeat[0] {
	case 'y':
		return addYear(t, now)
	case 'd':
		return addDays(t, now, repeat)
	case 'w':
		// err := errors.New("не реализовано")
		// return date, err
		return addWeeks(t, now, repeat)
	case 'm':
		// err := errors.New("не реализовано")
		// return date, err
		return addMonths(t, now, repeat)
	}
	return "", nil
}

func addYear(t time.Time, now time.Time) (string, error) {
	for {
		t = t.AddDate(1, 0, 0)
		if t.After(now) {
			break
		}
	}
	return t.Format("20060102"), nil
}

func addDays(t time.Time, now time.Time, repeat string) (string, error) {
	rep := strings.Split(repeat, " ")
	if len(rep) != 2 {
		err := errors.New("не указан интервал в днях")
		return "", err
	}
	daysNumber, err := strconv.Atoi(rep[1])
	if err != nil {
		return "", err
	}
	if daysNumber > 400 || daysNumber < 1 {
		err := errors.New("превышен максимально допустимый интервал")
		return "", err
	} else {
		for {
			t = t.AddDate(0, 0, daysNumber)
			if t.After(now) {
				break
			}
		}
		return t.Format("20060102"), nil
	}
}

func addWeeks(t time.Time, now time.Time, repeat string) (string, error) {
	var weekDays []int
	if len(repeat) < 3 {
		return "", errors.New("не указан интервал в днях недели")
	}

	weekDaysStr := strings.Split(repeat[2:], ",")
	if len(weekDaysStr) == 0 {
		return "", errors.New("не указан интервал в днях недели")
	}

	for _, day := range weekDaysStr {
		dayNumber, err := strconv.Atoi(day)
		if err != nil {
			return "", err
		}
		if dayNumber < 1 || dayNumber > 7 {
			return "", fmt.Errorf("недопустимое значение %d дня недели", dayNumber)
		}
		weekDays = append(weekDays, dayNumber)
	}

	for i, day := range weekDays {
		if day == 7 {
			weekDays[i] = 0
		}
	}

	sort.Ints(weekDays)
	var nextWeekDay int
	for _, wd := range weekDays {
		if wd >= int(t.Weekday()) { // Check if wd is greater than or equal to the current weekday
			nextWeekDay = wd
			break
		}
	}
	if nextWeekDay == 0 { // If no future weekday was found in this week, take the first day of the next week
		nextWeekDay = weekDays[0]
	}
	for {
		t = t.AddDate(0, 0, 1)
		if t.After(now) && int(t.Weekday()) == nextWeekDay {
			return t.Format("20060102"), nil
		}
	}
}

func addMonths(t time.Time, now time.Time, repeat string) (string, error) {
	var listOfDays time.Time
	var err error
	if len(repeat) < 3 {
		return "", errors.New("не указан интервал в месяцах")
	}

	repSlice := strings.Split(repeat, " ")

	if len(repSlice) == 2 {
		listOfDays, err = getNextDate(now, t.Format("20060102"), repeat)
		if err != nil {
			return "", err
		}
	} else if len(repSlice) == 3 {
		listOfDays = getNextMonthDate(now, t.Format("20060102"), repeat)
	} else {
		return "", errors.New("неверный формат повторения")
	}

	return listOfDays.Format("20060102"), nil
}

func getNextMonthDate(now time.Time, target string, rule string) time.Time {
	// Parse the rule
	ruleParts := strings.Split(rule, " ")
	daysPart := strings.Split(ruleParts[1], ",")
	monthsPart := strings.Split(ruleParts[2], ",")

	// Convert days and months to integers
	days := make([]int, len(daysPart))
	for i, day := range daysPart {
		days[i], _ = strconv.Atoi(day)
	}
	months := make([]int, len(monthsPart))
	for i, month := range monthsPart {
		months[i], _ = strconv.Atoi(month)
	}

	// Initialize targetTime and nearestDate
	targetTime, _ := time.Parse("20060102", target)
	var nearestDate time.Time

	// Loop until we find a date after today
	for {
		for _, day := range days {
			for _, month := range months {
				var date time.Time
				if day < 0 {
					// If day is negative, calculate the date from the end of the current month
					endOfMonth := time.Date(targetTime.Year(), time.Month(month+1), 0, 0, 0, 0, 0, time.UTC)
					date = endOfMonth.AddDate(0, 0, day+1)
				} else {
					date = time.Date(targetTime.Year(), time.Month(month), day, 0, 0, 0, 0, time.UTC)
				}
				// If the date is before now, add a year
				if date.Before(now) {
					date = date.AddDate(1, 0, 0)
				}
				// If this is the first date we've found, or it's earlier than the current nearest date, update nearestDate
				if (nearestDate.IsZero() || date.Before(nearestDate)) && date.After(now) {
					nearestDate = date
				}
			}
		}
		// If nearestDate is not zero, it means we have found a date after today, so break the loop
		if !nearestDate.IsZero() {
			break
		}
		// If no date after today is found in this year, increment the year
		targetTime = targetTime.AddDate(1, 0, 0)
	}

	return nearestDate
}

func getNextDate(now time.Time, target string, rule string) (time.Time, error) {
	// Parse target date
	targetTime, _ := time.Parse("20060102", target)

	// Split rule into "m" and the rest
	ruleParts := strings.SplitN(rule, " ", 2)
	if len(ruleParts) != 2 {
		return time.Time{}, fmt.Errorf("invalid rule format")
	}

	// Check if the rest of the rule contains a comma
	var monthDays []int
	if strings.Contains(ruleParts[1], ",") {
		daysParts := strings.Split(ruleParts[1], ",")
		monthDays = make([]int, len(daysParts))
		for i, part := range daysParts {
			day, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil {
				return time.Time{}, fmt.Errorf("invalid day in rule: %v", err)
			}
			if day < -2 || day > 31 {
				return time.Time{}, fmt.Errorf("не правильно указан формат повтора")
			}
			monthDays[i] = day
		}
	} else {
		day, err := strconv.Atoi(strings.TrimSpace(ruleParts[1]))
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid day in rule: %v", err)
		}
		if day < -2 || day > 31 {
			return time.Time{}, fmt.Errorf("не правильно указан формат повтора")
		}
		monthDays = []int{day}
	}
	// Find the nearest date
	var nearestDate time.Time
	for {
		for _, day := range monthDays {
			var date time.Time
			// Define endOfMonth here
			endOfMonth := time.Date(targetTime.Year(), targetTime.Month()+1, 0, 0, 0, 0, 0, time.UTC)
			if day < 0 {
				// If day is negative, calculate the date from the end of the current month
				date = endOfMonth.AddDate(0, 0, day+1)
			} else {
				if day > endOfMonth.Day() {
					targetTime = targetTime.AddDate(0, 1, 0)
					endOfMonth = time.Date(targetTime.Year(), targetTime.Month()+1, 0, 0, 0, 0, 0, time.UTC)
				}
				date = time.Date(targetTime.Year(), targetTime.Month(), day, 0, 0, 0, 0, time.UTC)
			}
			// Now endOfMonth is accessible here
			if date.Before(now) || date.After(endOfMonth) {
				date = date.AddDate(0, 1, 0)
			}
			// If this is the first date we've found, or it's earlier than the current nearest date, update nearestDate
			if (nearestDate.IsZero() || date.Before(nearestDate)) && date.After(now) {
				nearestDate = date
				// break
			}
		}
		// If nearestDate is not zero, it means we have found a date after today, so break the loop
		if !nearestDate.IsZero() {
			break
		}
		// If no date after today is found in this month, increment the month
		targetTime = targetTime.AddDate(0, 1, 0)
	}

	return nearestDate, nil
}
