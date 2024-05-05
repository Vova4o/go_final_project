package handlers

import (
	"errors"
	"testing"
	"time"
)

func TestCheckTask(t *testing.T) {
	tests := []struct {
		name string
		task task
		want error
	}{
		{
			name: "Valid task",
			task: task{
				Date:    "20240427", // tomorrow's date
				Title:   "Test Task",
				Comment: "This is a test task",
				Repeat:  "d 1",
			},
			want: nil,
		},
		{
			name: "Wrong next date",
			task: task{
				Date:    time.Now().AddDate(0, 0, -1).Format("20060102"), // yesterday's date
				Title:   "Test Task",
				Comment: "This is a test task",
				Repeat:  "d",
			},
			want: errors.New("неверное правило повторения"),
		},
		{
			name: "Wrong repeat rule",
			task: task{
				Date:    "20210101",
				Title:   "Test Task",
				Comment: "",
				Repeat:  "l",
			},
			want: errors.New("неверное правило повторения"),
		},
		{
			name: "No title",
			task: task{
				Date:    "20210101",
				Title:   "",
				Comment: "",
				Repeat:  "l",
			},
			want: errors.New("не указан заголовок задачи"),
		},
		{
			name: "Valid task",
			task: task{
				Date:    "20241001",
				Title:   "Test Task",
				Comment: "",
				Repeat:  "y",
			},
			want: nil,
		},
		{
			name: "wrong date format",
			task: task{
				Date:    "2024-10-01",
				Title:   "Test Task",
				Comment: "",
				Repeat:  "y",
			},
			want: errors.New("дата представлена в формате, отличном от 20060102"),
		},

		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkTask(tt.task)
			if err != nil && tt.want != nil {
				if err.Error() != tt.want.Error() {
					t.Errorf("checkTask() data = %v, error = %v, wantErr %v", got, err, tt.want)
				}
			} else if err != tt.want {
				t.Errorf("checkTask() data = %v, error = %v, wantErr %v", got, err, tt.want)
			}
		})
	}
}

type testCase struct {
	name     string
	task     task
	wantTask *task
	wantErr  error
}

func TestCheckTask2(t *testing.T) {
	tests := []testCase{
		{
			name: "wrong date format",
			task: task{
				Date:    "2024-10-01",
				Title:   "Test Task",
				Comment: "",
				Repeat:  "y",
			},
			wantTask: &task{
				Date:    "2024-10-01", // Expected date
				Title:   "Test Task",
				Comment: "",
				Repeat:  "y",
			},
			wantErr: errors.New("дата представлена в формате, отличном от 20060102"),
		},
		{
			name: "test case 2",
			task: task{
				Date:    "20240129",
				Title:   "",
				Comment: "",
				Repeat:  "",
			},
			wantTask: &task{
				Date:    "20240129",
				Title:   "", // Expected title
				Comment: "",
				Repeat:  "",
			},
			wantErr: errors.New("не указан заголовок задачи"),
			// Add your expected task and error here
		},
		{
			name: "test case 3",
			task: task{
				Date:    "20240192",
				Title:   "Qwerty",
				Comment: "",
				Repeat:  "",
			},
			wantTask: &task{
				Date:    "20240192",
				Title:   "Qwerty", // Expected title
				Comment: "",
				Repeat:  "",
			},
			wantErr: errors.New("дата представлена в формате, отличном от 20060102"),
			// Add your expected task and error here
		},
		{
			name: "test case 4",
			task: task{
				Date:    "28.01.2024",
				Title:   "Заголовок",
				Comment: "",
				Repeat:  "",
			},
			wantTask: &task{
				Date:    "28.01.2024", // Expected date
				Title:   "Заголовок",
				Comment: "",
				Repeat:  "",
			},
			wantErr: errors.New("дата представлена в формате, отличном от 20060102"),

			// Add your expected task and error here
		},
		{
			name: "test case 5",
			task: task{
				Date:    "20240112",
				Title:   "Заголовок",
				Comment: "",
				Repeat:  "w",
			},
			wantTask: &task{
				Date:    "20240112", // Expected date
				Title:   "Заголовок",
				Comment: "",
				Repeat:  "w",
			},
			wantErr: errors.New("неверное правило повторения"),
			// Add your expected task and error here
		},
		{
			name: "test case 6",
			task: task{
				Date:    "20240212",
				Title:   "Заголовок",
				Comment: "",
				Repeat:  "ooops",
			},
			wantTask: &task{
				Date:    "20240212", // Expected date
				Title:   "Заголовок",
				Comment: "",
				Repeat:  "ooops",
			},
			wantErr: errors.New("неверное правило повторения"),
			// Add your expected task and error here
		},
		{
			name: "test case 7",
			task: task{
				Date:    "today",
				Title:   "Заголовок",
				Comment: "",
				Repeat:  "",
			},
			wantTask: &task{
				Date:    "today", // Expected date
				Title:   "Заголовок",
				Comment: "",
				Repeat:  "",
			},
			wantErr: errors.New("дата представлена в формате, отличном от 20060102"),
			// Add your expected task and error here
		},
		{
			name: "test case 8",
			task: task{
				Date:    "20231225",
				Title:   "Заголовок",
				Comment: "",
				Repeat:  "d 12",
			},
			wantTask: &task{
				Date:    "20240505", // Expected date
				Title:   "Заголовок",
				Comment: "",
				Repeat:  "d 12",
			},
			wantErr: nil,
			// Add your expected task and error here
		},
		{
			name: "test case 9",
			task: task{
				Date:    "20250426",
				Title:   "Заголовок",
				Comment: "",
				Repeat:  "d 1",
			},
			wantTask: &task{
				Date:    "20250426", // Expected date
				Title:   "Заголовок",
				Comment: "",
				Repeat:  "d 1",
			},
			wantErr: nil,
			// Add your expected task and error here
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkTask(tt.task)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("checkTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantTask != nil && got.Date != tt.wantTask.Date {
				t.Errorf("checkTask() = %v, want %v", got, tt.wantTask)
			}
		})
	}
}
