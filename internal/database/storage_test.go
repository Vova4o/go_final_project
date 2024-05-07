package database

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/vova4o/go_final_project/internal/models"
	"github.com/vova4o/go_final_project/internal/nextdate"
)

func TestTasks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := &Storage{Db: db}

	rows := sqlmock.NewRows([]string{"id", "date", "title", "comment", "repeat"}).
		AddRow("1", "20240131", "Заголовок задачи", "", "").
		AddRow("2", "20240131", "Фитнес", "", "d 3")

	mock.ExpectQuery("^SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 10 OFFSET (.+)$").WillReturnRows(rows)

	tasks, err := s.Tasks(0)
	if err != nil {
		t.Errorf("error was not expected while getting tasks: %s", err)
	}

	want := []models.DBTask{
		{ID: "1", Date: "20240131", Title: "Заголовок задачи", Comment: "", Repeat: ""},
		{ID: "2", Date: "20240131", Title: "Фитнес", Comment: "", Repeat: "d 3"},
	}

	assert.Equal(t, want, tasks, "The two arrays should be the same.")
}

/*
{"20240125", "w 1,2,3", "20240129"},
{"20240126", "w 7", "20240128"},
{"20230126", "w 4,5", "20240201"},
{"20230226", "w 8,4,5", ""},
*/
func TestAddWeeks(t *testing.T) {
	tests := []struct {
		name    string
		t       time.Time
		repeat  string
		want    string
		wantErr bool
	}{
		{
			name:    "add w 1,2,3 mo,tu,we",
			t:       time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
			repeat:  "w 1,2,3",
			want:    "20240129",
			wantErr: false,
		},
		{
			name:    "add w 7",
			t:       time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC),
			repeat:  "w 7",
			want:    "20240128",
			wantErr: false,
		},
		{
			name:    "add w 4,5",
			t:       time.Date(2023, 1, 26, 0, 0, 0, 0, time.UTC),
			repeat:  "w 4,5",
			want:    "20240201",
			wantErr: false,
		},
		{
			name:    "add w 8,4,5",
			t:       time.Date(2023, 2, 26, 0, 0, 0, 0, time.UTC),
			repeat:  "w 8,4,5",
			want:    "",
			wantErr: true,
		},
		{
			name:    "no w",
			t:       time.Date(2023, 2, 26, 0, 0, 0, 0, time.UTC),
			repeat:  "w",
			want:    "",
			wantErr: true,
		},
		{
			name:    "add w 5",
			t:       time.Date(2023, 1, 26, 0, 0, 0, 0, time.UTC),
			repeat:  "w 5",
			want:    "20240202",
			wantErr: false,
		},
		// add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nowStr := "20240126"
			now, _ := time.Parse("20060102", nowStr)
			got, err := nextdate.AddWeeks(tt.t, now, tt.repeat)
			if (err != nil) != tt.wantErr {
				t.Errorf("addWeeks(%v, %v) returned error: %v, wantErr: %v", tt.t, tt.repeat, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("addWeeks(%v, %v) = %v, want %v", tt.t, tt.repeat, got, tt.want)
			}
		})
	}
}

func TestAddMonth(t *testing.T) {
	tests := []struct {
		name    string
		t       time.Time
		repeat  string
		want    string
		wantErr bool
	}{
		/*
			{"20231106", "m 13", "20240213"},
				{"20240120", "m 40,11,19", ""},
				{"20240116", "m 16,5", "20240205"},
				{"20240126", "m 25,26,7", "20240207"},
				{"20240409", "m 31", "20240531"},
				{"20240329", "m 10,17 12,8,1", "20240810"},
				{"20230311", "m 07,19 05,6", "20240507"},
				{"20230311", "m 1 1,2", "20240201"},
				{"20240222", "m -2,-3", ""},
				{"20240326", "m -1,-2", "20240330"},
		*/
		{
			name:    "add m 31",
			t:       time.Date(2024, 4, 9, 0, 0, 0, 0, time.UTC),
			repeat:  "m 31",
			want:    "20240531",
			wantErr: false,
		},
		{
			name:    "add m 07,19 05,6",
			t:       time.Date(2023, 3, 11, 0, 0, 0, 0, time.UTC),
			repeat:  "m 07,19 05,6",
			want:    "20240507",
			wantErr: false,
		},
		{
			name:    "add m 10,17 12,8,1",
			t:       time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC),
			repeat:  "m 10,17 12,8,1",
			want:    "20240810",
			wantErr: false,
		},
		{
			name:    "add m 1 1,2",
			t:       time.Date(2023, 3, 11, 0, 0, 0, 0, time.UTC),
			repeat:  "m 1 1,2",
			want:    "20240201",
			wantErr: false,
		},
		{
			name:    "add m 13",
			t:       time.Date(2023, 11, 6, 0, 0, 0, 0, time.UTC),
			repeat:  "m 13",
			want:    "20240213",
			wantErr: false,
		},
		{
			name:    "add m -1,-2",
			t:       time.Date(2024, 3, 26, 0, 0, 0, 0, time.UTC),
			repeat:  "m -1,-2",
			want:    "20240330",
			wantErr: false,
		},
		{
			name:    "add m -1,18",
			t:       time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			repeat:  "m -1,18",
			want:    "20240218",
			wantErr: false,
		},
		{
			name:    "add m -1",
			t:       time.Date(2024, 1, 27, 0, 0, 0, 0, time.UTC),
			repeat:  "m -1",
			want:    "20240131",
			wantErr: false,
		},
		{
			name:    "add m -2",
			t:       time.Date(2024, 2, 22, 0, 0, 0, 0, time.UTC),
			repeat:  "m -2",
			want:    "20240228",
			wantErr: false,
		},
		{
			name:    "add m 10,17 12,8,1",
			t:       time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC),
			repeat:  "m 10,17 12,8,1",
			want:    "20240810",
			wantErr: false,
		},
		{
			name:    "add m 07,19 05,6",
			t:       time.Date(2023, 3, 11, 0, 0, 0, 0, time.UTC),
			repeat:  "m 07,19 05,6",
			want:    "20240507",
			wantErr: false,
		},
		{
			name:    "add m 40,11,19",
			t:       time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			repeat:  "m 40,11,19",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nowStr := "20240126"
			now, _ := time.Parse("20060102", nowStr)
			got, err := nextdate.AddMonths(tt.t, now, tt.repeat)
			if (err != nil) != tt.wantErr {
				t.Errorf("addMonth(%v, %v) returned error: %v, wantErr: %v", tt.t, tt.repeat, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("addMonth(%v, %v) = %v, want %v", tt.t, tt.repeat, got, tt.want)
			}
		})
	}
}

func TestSearchTasksByDate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	s := &Storage{Db: db}
	// Mock the query
	rows := sqlmock.NewRows([]string{"id", "date", "title", "comment", "repeat"}).
		AddRow("1", "20220101", "Test Title", "Test Comment", "Test Repeat")
	mock.ExpectQuery("^SELECT (.+) FROM scheduler WHERE date = ?").WithArgs("20220101").WillReturnRows(rows)

	// Test SearchTasksByDate function
	tasks, err := s.TasksByDate("20220101")
	if err != nil {
		t.Fatalf("Failed to search tasks by date: %v", err)
	}

	// Check if the returned tasks are correct
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
	if tasks[0].ID != "1" || tasks[0].Date != "20220101" || tasks[0].Title != "Test Title" || tasks[0].Comment != "Test Comment" || tasks[0].Repeat != "Test Repeat" {
		t.Errorf("Returned task does not match expected task")
	}
}
