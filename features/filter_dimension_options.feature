Feature: Get Filter Private Endpoints Not Enabled

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
          "dimensions": {
            "href": ":27100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
          },
          "self": {
            "href": ":27100/filters/94310d8d-72d6-492a-bc30-27584627edb1"
          }
        },
        "instance_id": "c733977d-a2ca-4596-9cb1-08a6e724858b",
        "dimensions": [
          {
            "name": "silbings",
            "id": "siblings_3",
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
            "options": [
              "0",
              "1"
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
      }
    ]
    """

  Scenario: Add filter dimension option successfully
    When I POST "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography/options/2"
    """
    """

    Then I should receive the following JSON response:
    """
    {
      "option": "2",
      "Links": {
        "self": {
          "href": ":27100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography/options/2"
        },
        "filter": {
          "href": ":27100/filters/94310d8d-72d6-492a-bc30-27584627edb1"
        },
        "dimension": {
          "href": ":27100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography",
          "id": "city"
        }
      }
    }
    """

    And the HTTP status code should be "201"

  Scenario: Add filter dimension option but dimension not in filter
    When I POST "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/foo/options/2"
    """
    """

    Then I should receive the following JSON response:
    """
    {
      "errors": ["failed to add dimension option: dimension not found in filter"]
    }
    """

    And the HTTP status code should be "400"

  Scenario: Add filter dimension option but option already in filter
    When I POST "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography/options/1"
    """
    """

    Then I should receive the following JSON response:
    """
    {
      "errors": ["failed to add dimension option: option already added to dimension"]
    }
    """

    And the HTTP status code should be "400"

  Scenario: Add filter dimension option but Cantabular is failing
    Given Cantabular responds with an error

    When I POST "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography/options/2"
    """
    """

    Then I should receive the following JSON response:
    """
    {
      "errors": ["Internal Server Error"]
    }
    """

    And the HTTP status code should be "500"