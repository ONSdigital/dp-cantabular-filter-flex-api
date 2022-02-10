Feature: Filters Private Endpoints Not Enabled

  Background:
    Given private endpoints are not enabled

  Scenario: Creating a new filter happy

    When I POST "/filters"
    """
    {
      "dataset":{
          "id":      "c7b634c9-b4e9-4e7a-a0b8-d255d38db200",
          "edition": "2021",
          "version": 1
      },
      "population_type": "Example",
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
          "href": "localhost:8082/datasets/c7b634c9-b4e9-4e7a-a0b8-d255d38db200/editions/2021/version/1"
        },
        "self": {
          "href": ":27100/flex/filters/94310d8d-72d6-492a-bc30-27584627edb1"
        }
      },
      "events": null,
      "unique_timestamp": "2022-01-26T12:27:04.783936865Z",
      "last_updated": "2022-01-26T12:27:04.783936865Z",
      "etag": "defad5e5f9bf579bf894b3fbee3c95cba4321264",
      "instance_id": "00000000-0000-0000-0000-000000000000",
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
      "population_type": ""
    }
    """

    And the HTTP status code should be "201"

Scenario: Creating a new filter bad request body

    When I POST "/filters"
    """
    {
      "ins
    """

    Then I should receive the following JSON response:
    """
    {
      "errors": [
        "badly formed request body: unexpected end of JSON input"
      ]
    }
    """

    And the HTTP status code should be "400"

  Scenario: Creating a new invalid request

    When I POST "/filters"
    """
    {
      "dataset":{
          "id":      "c7b634c9-b4e9-4e7a-a0b8-d255d38db200",
          "edition": "2021",
          "version": 1
       },
      "population_type": "Example",
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
        }
      ]
    }
    """

    Then I should receive the following JSON response:
    """
    {
      "errors": [
        "failed to parse request: invalid request: missing/invalid field: 'dimensions' must contain at least 2 values"
      ]
    }
    """

    And the HTTP status code should be "400"
