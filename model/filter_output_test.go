package model_test

import (
	"encoding/json"
	"testing"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDownloadsOrder(t *testing.T) {
	Convey("Given a Downloads struct", t, func() {
		d := model.Downloads{
			XLS: &model.FileInfo{
				HREF: "XLSX",
			},
			CSV: &model.FileInfo{
				HREF: "CSV",
			},
			CSVW: &model.FileInfo{
				HREF: "CSVW",
			},
			TXT: &model.FileInfo{
				HREF: "TXT",
			},
		}

		Convey("When marshalled", func() {
			b, err := json.Marshal(d)
			So(err, ShouldBeNil)

			Convey("The downloads should be in the expected order", func() {
				expected := `{"xls":{"href":"XLSX","size":"","public":"","private":"","skipped":false},"csv":{"href":"CSV","size":"","public":"","private":"","skipped":false},"txt":{"href":"TXT","size":"","public":"","private":"","skipped":false},"csvw":{"href":"CSVW","size":"","public":"","private":"","skipped":false}}`
				So(string(b), ShouldResemble, expected)
			})
		})
	})
}
