Feature: Filter Dimensions Private Endpoints Not Enabled

  Background:
    Given private endpoints are not enabled
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
            "href": ":27100/flex/filters/94310d8d-72d6-492a-bc30-27584627edb1"
          }
        },
        "events": null,
        "unique_timestamp": "2022-01-26T12:27:04.783936865Z",
        "last_updated": "2022-01-26T12:27:04.783936865Z",
        "etag": "6e627be2c355c7cebe5de4af6a0c6d75c9523fb3",
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


  Scenario: Get filter dimensions successfully
    When I GET "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
    Then I should receive the following JSON response:
    """
    {
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
    And the HTTP status code should be "200"

  Scenario: Get filter dimensions unsuccessfully
    When I GET "/filters/94310d8d-72d6-492a-03cb-27584627edb1/dimensions"
    Then I should receive the following JSON response:
    """
    {
      "errors": ["failed to get filter dimensions"]
    }
    """
    And the HTTP status code should be "404"
