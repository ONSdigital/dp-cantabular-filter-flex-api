package api

import (
	"testing"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_validateFilterOutput(t *testing.T) {
	type args struct {
		filterOutput model.FilterOutput
	}

	Convey("Given a set of blank filter outputs, test the validation logic", t, func() {
		blankInfo := model.FileInfo{
			HREF:    " ",
			Size:    " ",
			Public:  " ",
			Private: " ",
			Skipped: true,
		}
		partialblankInfo := model.FileInfo{
			HREF:    " ",
			Size:    " ",
			Public:  "test1 test ",
			Private: " ",
			Skipped: true,
		}
		tests := []struct {
			name string
			args args
		}{
			{
				name: "test1",
				args: args{filterOutput: model.FilterOutput{
					CSV:  &blankInfo,
					CSVW: new(model.FileInfo),
					TXT:  new(model.FileInfo),
					XLS:  &partialblankInfo,
				},
				},
			},
			{
				name: "test2",
				args: args{filterOutput: model.FilterOutput{
					CSV:  &blankInfo,
					CSVW: &blankInfo,
					TXT:  &partialblankInfo,
					XLS:  &partialblankInfo,
				},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if err := validateFilterOutput(tt.args.filterOutput); err == nil {

					t.Errorf("validateFilterOutput() error = %v", err)
				}
			})
		}
	})

	Convey("Given a set of valid filter outputs, test the validation logic", t, func() {
		validInfo := model.FileInfo{
			HREF:    "href string ",
			Size:    "size string ",
			Public:  "public string ",
			Private: "private string",
			Skipped: true,
		}

		tests := []struct {
			name string
			args args
		}{
			{
				name: "test1",
				args: args{filterOutput: model.FilterOutput{
					CSV:  &validInfo,
					CSVW: &validInfo,
					TXT:  &validInfo,
					XLS:  &validInfo,
				},
				},
			},
			{
				name: "test2",
				args: args{filterOutput: model.FilterOutput{
					CSV:  &validInfo,
					CSVW: &validInfo,
					TXT:  &validInfo,
					XLS:  &validInfo,
				},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if err := validateFilterOutput(tt.args.filterOutput); err != nil {

					t.Errorf("validateFilterOutput() error = %v", err)
				}
			})
		}
	})

}
