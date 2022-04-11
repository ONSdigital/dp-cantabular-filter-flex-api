Feature: Submit Filter Private Endpoints Not Enabled

  Background:
    Given private endpoints are not enabled
    And I have these filters:
    """
    [
      {
        "filter_id": "TEST-FILTER-ID",
        "links": {
          "version": {
            "href": "http://mockhost:9999/datasets/cantabular-example-1/editions/2021/version/1",
            "id": "1"
          },
          "dimensions": {
            "href": ":27100/filters/TEST-FILTER-ID/dimensions"
          },
          "self": {
            "href": ":27100/filters/TEST-FILTER-ID"
          }
        },
        "instance_id": "TEST-INSTANCE-ID",
        "dimensions": [
          {
            "name": "Number of siblings (3 mappings)",
            "options": [
              "0-3",
              "4-7",
              "7+"
            ],
            "is_area_type": false
          },
          {
            "name": "City",
            "options": [
              "Cardiff",
              "London",
              "Swansea"
            ],
            "is_area_type": true
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
      },
      {
        "filter_id": "test-case-2",
        "links": {
          "version": {
            "href": "http://mockhost:9999/datasets/cantabular-example-unpublished/editions/2021/version/1",
            "id": "1"
          },
          "dimensions": {
            "href": ":27100/filters/test-case-2/dimensions"
          },
          "self": {
            "href": ":27100/filters/test-case-2"
          }
        },
        "instance_id": "c733977d-a2ca-4596-9cb1-08a6e724858b",
        "dimensions": [
          {
            "name": "Number of siblings (3 mappings)",
            "options": [
              "0-3",
              "4-7",
              "7+"
            ],
            "is_area_type": false
          },
          {
            "name": "City",
            "options": [
              "Cardiff",
              "London",
              "Swansea"
            ],
            "is_area_type": true
          }
        ],
        "dataset": {
          "id": "cantabular-example-unpublished",
          "edition": "2021",
          "version": 1
        },
        "population_type": "Example",
        "type": "flexible"
      }
    ]
    """
  Scenario: Submit Filter Not Found
    When I POST "/filters/cannot-find/submit"
    """
    """
    Then the HTTP status code should be "404"

  Scenario: Submit filter successfully
    When I POST "/filters/TEST-FILTER-ID/submit"
    """
    """

    Then I should receive the following JSON response:
    """
    {
      "instance_id":"TEST-INSTANCE-ID",
      "filter_id":  "TEST-FILTER-ID",
      "events":[
        {
          "timestamp": "2022-01-26T12:27:04Z",
          "name": "cantabular-export-start"
        }
      ],
      "dataset":{
        "id":"cantabular-example-1",
        "edition":"2021",
        "version": 1
      },
      "links": {
        "version": {
          "href": "http://mockhost:9999/datasets/cantabular-example-1/editions/2021/version/1",
          "id": "1"
        },
        "self": {
          "href": ":27100/filters/TEST-FILTER-ID"
        },
        "dimensions": {
          "href": ":27100/filters/TEST-FILTER-ID/dimensions"
        }
      },
      "population_type": "Example"
    }

    """
    And the HTTP status code should be "202"

    And one event with the following fields are in the produced kafka topic cantabular-export-start:
      | InstanceID        | DatasetID            | Edition          | Version          | FilterID       |
      | TEST-INSTANCE-ID  | cantabular-example-1 | 2021             | 1                | TEST-FILTER-ID |
