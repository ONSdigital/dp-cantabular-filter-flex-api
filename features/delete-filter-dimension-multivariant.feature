Feature: Filter Dimensions Private Endpoints Not Enabled

  Background:
    Given private endpoints are not enabled
    And the following version document with dataset id "cantabular-example-1", edition "2021" and version "1" is available from dp-dataset-api:
    """
    {
      "alerts": [],
      "collection_id": "dfb-38b11d6c4b69493a41028d10de503aabed3728828e17e64914832d91e1f493c6",
      "dimensions": [
        {
          "label": "City",
          "links": {
            "code_list": {},
            "options": {},
            "version": {}
          },
          "href": "http://api.localhost:23200/v1/code-lists/city",
          "id": "city",
          "name": "geography"
        },
        {
          "label": "Number of siblings (3 mappings)",
          "links": {
            "code_list": {},
            "options": {},
            "version": {}
          },
          "href": "http://api.localhost:23200/v1/code-lists/siblings",
          "id": "siblings_3",
          "name": "siblings"
        }
      ],
      "edition": "2021",
      "id": "c733977d-a2ca-4596-9cb1-08a6e724858b",
      "links": {
        "dataset": {
          "href": "http://dp-dataset-api:22000/datasets/cantabular-example-1",
          "id": "cantabular-example-1"
        },
        "dimensions": {},
        "edition": {
          "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021",
          "id": "2021"
        },
        "self": {
          "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1"
        }
      },
      "release_date": "2021-11-19T00:00:00.000Z",
      "state": "published",
      "usage_notes": [],
      "version": 1
    }
    """

    And I have these filters:
    """
    [
      {
        "filter_id": "94310d8d-72d6-492a-bc30-27584627edb1",
        "links": {
          "version": {
            "href": "http://mockhost:9999/datasets/cantabular-example-1/editions/2021/version/1",
            "id": "1"
          },
          "self": {
            "href": ":27100/filters/94310d8d-72d6-492a-bc30-27584627edb1"
          }
        },
        "instance_id": "c733977d-a2ca-4596-9cb1-08a6e724858b",
        "dimensions": [
          {
            "name": "siblings",
            "id": "siblings_3",
            "label": "Number of siblings (3 mappings)",
            "options": [
              "0-3",
              "4-7",
              "7+"
            ],
            "is_area_type": false
          },
          {
            "name": "geography",
            "id": "city",
            "label": "City",
            "options": [
              "Cardiff",
              "London",
              "Swansea"
            ],
            "is_area_type": true,
            "filter_by_parent":"country"
          },
          {
            "name": "accommodation_type",
            "id": "accommodation_type",
            "label": "accommodation type",
            "options": [],
            "is_area_type": false
          }
        ],
        "dataset": {
          "id": "cantabular-example-1",
          "edition": "2021",
          "version": 1
        },
        "published": true,
        "population_type": "Example",
        "type": "flexible"
      }
    ]
    """

  Scenario: Delete filter dimensions successfully
    When I DELETE "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/accommodation_type"
    Then the HTTP status code should be "204"

  Scenario: Delete area_type filter dimensions error
    When I DELETE "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography"
    Then the HTTP status code should be "409"
  
  Scenario: Delete filter dimensions with only 2 dimensions in filter error
    When I DELETE "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/accommodation_type"
    Then the HTTP status code should be "204"
    When I DELETE "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/siblings"
    Then the HTTP status code should be "409"