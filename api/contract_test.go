package api

import (
	"github.com/google/uuid"
)

func testUUID() uuid.UUID{
	const idstr := "307c53db-4495-436f-8f2e-8435deb8144e"
	id, err := uuid.Parse(idstr)
	So(err, ShouldBeNil)

	return id
}

func TestCreateFiltersRequestValid(t *testing.T){
	Convey("Given a valid createFilterRequest request object", t, func() {
		req := createFilterRequest{
			InstanceID: 
			DatasetID      string            `bson:"dataset_id"      json:"dataset_id"`
			Edition        string            `bson:"edition"         json:"edition"`
			Version        int               `bson:"version"         json:"version"`
			CantabularBlob string            `bson:"cantabular_blob" json:"cantabular_blob"`
			Dimensions     []model.Dimension `bson:"dimensions"      json:"dimensions"`
		}

		Convey("When ParseRequest(body, request) is called", func() {
			err := api.ParseRequest(bytes.NewReader(b), &req)
			So(err, ShouldBeNil)

			Convey("The request object should be populated with the expected values ", func() {
				expected := fooRequest{
					Foo:   "I am foo",
					Bar:   2,
					IsValid: true,
				}
				So(req, ShouldResemble, expected)
			})
		})
	})
}

