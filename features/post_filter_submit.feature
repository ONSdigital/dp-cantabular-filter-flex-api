@post
Feature: Post Filter Private Endpoints Not Enabled

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
          "self": {
            "href": ":27100/filters/94310d8d-72d6-492a-bc30-27584627edb1"
          }
        },
        "events": null,
        "instance_id": "TEST-INSTANCE-ID",
        "dimensions": [
          {
            "name": "Number of siblings (3 mappings)",
            "options": [
              "0-3",
              "4-7",
              "7+"
            ],
            "dimension_url": "http://dimension.url/siblings",
            "is_area_type": false
          },
          {
            "name": "City",
            "options": [
              "Cardiff",
              "London",
              "Swansea"
            ],
            "dimension_url": "http://dimension.url/city",
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
          "self": {
            "href": ":27100/filters/83210d8d-72d6-492a-bc30-27584627abc2"
          }
        },
        "events": null,
        "instance_id": "c733977d-a2ca-4596-9cb1-08a6e724858b",
        "dimensions": [
          {
            "name": "Number of siblings (3 mappings)",
            "options": [
              "0-3",
              "4-7",
              "7+"
            ],
            "dimension_url": "http://dimension.url/siblings",
            "is_area_type": false
          },
          {
            "name": "City",
            "options": [
              "Cardiff",
              "London",
              "Swansea"
            ],
            "dimension_url": "http://dimension.url/city",
            "is_area_type": true
          }
        ],
        "dataset": {
          "id": "cantabular-example-unpublished",
          "edition": "2021",
          "version": 1
        },
        "published": false,
        "population_type": "Example",
        "type": "flexible"
      }
    ]
    """
  Scenario: POST Filter Not Found
    When I POST "/filters/cannot-find/submit"
    """
    """
    Then the HTTP status code should be "500"
  Scenario: POST filter successfully
    When I POST "/filters/TEST-FILTER-ID/submit"
    """
    """
    Then I should receive the following time ignored JSON response:
    """
        {

              "instance_id":"TEST-INSTANCE-ID",
              "filter_id":"TEST-FILTER-ID",

              "dimension_list_url": "",
"events":[{
                      "timestamp": "2016-07-17T08:38:25.316Z",
                      "name": "cantabular-export-start"
              }],
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
            "href": ":27100/filters/94310d8d-72d6-492a-bc30-27584627edb1"
          }
        },
        "population_type": "Example",
        "dimensions": [
          {
            "name": "Number of siblings (3 mappings)",
            "options": [
              "0-3",
              "4-7",
              "7+"
            ],
            "dimension_url": "http://dimension.url/siblings",
            "is_area_type": false
          },
          {
            "name": "City",
            "options": [
              "Cardiff",
              "London",
              "Swansea"
            ],
            "dimension_url": "http://dimension.url/city",
            "is_area_type": true
          }
        ]
        }
    """
    And the HTTP status code should be "202"
    And one event with the following fields are in the produced kafka topic catabular-export-start:
      | InstanceID        | DatasetID            | Edition          | Version          | FilterID       |
      | TEST-INSTANCE-ID  | cantabular-example-1 | 2021             | 1                | TEST-FILTER-ID |
