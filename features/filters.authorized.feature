Feature: Filters Private Endpoints Enabled

  Background:
    Given private endpoints are enabled

  Scenario: Creating a new filter journey when authorized
    Given I am identified as "user@ons.gov.uk"
    
    And I am authorised
    
    When I POST "/filters"
    """
    {
      "instance_id":     "054aa093-1c31-46dd-9472-14ff0b86abce",
      "dataset_id":      "c7b634c9-b4e9-4e7a-a0b8-d255d38db200",
      "edition":         "2021",
      "cantabular_blob": "Example",
      "version":         1,
      "dimensions": [
        {
          "name": "Number Of Siblings (3 categories)",
          "options": [
            "0-3",
            "4-7",
            "7+"
          ],
          "dimension_url": "http://dimension.url/siblings",
          "is_area_type": false
        },{
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

    Then I should receive the following JSON response:
    """
    {
      "filter_id": "94310d8d-72d6-492a-bc30-27584627edb1",
      "links": {
        "version": {
          "href": "localhost:8082/datasets/c7b634c9-b4e9-4e7a-a0b8-d255d38db200/editions/2021/version/1",
          "id": "054aa093-1c31-46dd-9472-14ff0b86abce"
        },
        "self": {
          "href": ":27100/flex/filters/94310d8d-72d6-492a-bc30-27584627edb1"
        }
      },
      "events": null,
      "unique_timestamp": "2022-01-26T12:27:04.783936865Z",
      "last_updated": "2022-01-26T12:27:04.783936865Z",
      "etag": "",
      "instance_id": "054aa093-1c31-46dd-9472-14ff0b86abce",
      "dimensions": [
        {
          "name": "Number Of Siblings (3 categories)",
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
        "id": "c7b634c9-b4e9-4e7a-a0b8-d255d38db200",
        "edition": "2021",
        "version": 1
      },
      "published": true,
      "cantabular_blob": ""
    }
    """

    And the HTTP status code should be "201"

    And the document in the database for id "94310d8d-72d6-492a-bc30-27584627edb1" should be:
    """
    {
      "filter_id": "94310d8d-72d6-492a-bc30-27584627edb1",
      "links": {
        "version": {
          "href": "localhost:8082/datasets/c7b634c9-b4e9-4e7a-a0b8-d255d38db200/editions/2021/version/1",
          "id": "054aa093-1c31-46dd-9472-14ff0b86abce"
        },
        "self": {
          "href": ":27100/flex/filters/94310d8d-72d6-492a-bc30-27584627edb1"
        }
      },
      "events": null,
      "unique_timestamp": "2022-01-26T12:27:04.783936865Z",
      "last_updated": "2022-01-26T12:27:04.783936865Z",
      "etag": "",
      "instance_id": "054aa093-1c31-46dd-9472-14ff0b86abce",
      "dimensions": [
        {
          "name": "Number Of Siblings (3 categories)",
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
        "id": "c7b634c9-b4e9-4e7a-a0b8-d255d38db200",
        "edition": "2021",
        "version": 1
      },
      "published": true,
      "cantabular_blob": ""
    }
    """

  Scenario: Creating a new filter journey when not authorized
    Given I am not identified
    
    And I am not authorised
    
    When I POST "/filters"
    """
    {"foo":"bar"}
    """

    Then the HTTP status code should be "401"
