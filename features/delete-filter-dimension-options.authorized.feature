Feature: Delete Filter Dimension Options

  Background:
    Given private endpoints are not enabled
    And I have these filters:
    """
    [
      {
        "filter_id": "94310d8d-72d6-492a-bc30-27584627edb1",
        "links": {
          "version": {
            "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1",
            "id": "1"
          },
          "dimensions": {
            "href": ":27100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
          },
          "self": {
            "href": ":27100/filters/94310d8d-72d6-492a-bc30-27584627edb1"
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
        "filter_id": "83210d8d-72d6-492a-bc30-27584627abc2",
        "links": {
          "version": {
            "href": "http://mockhost:9999/datasets/cantabular-example-unpublished/editions/2021/version/1",
            "id": "1"
          },
          "dimensions": {
            "href": ":27100/filters/83210d8d-72d6-492a-bc30-27584627abc2/dimensions"
          },
          "self": {
            "href": ":27100/filters/83210d8d-72d6-492a-bc30-27584627abc2"
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
        "published": false,
        "population_type": "Example",
        "type": "flexible"
      }
    ]
    """
    Scenario: Delete options successfully
    When I DELETE "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/City/options"
    Then the HTTP status code should be "204"
    And a document in collection "filters" with key "filter_id" value "94310d8d-72d6-492a-bc30-27584627edb1" has empty "City" options

    Scenario: Delete Option Dimension Name not found
    When I DELETE "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/NOT-EXISTS/options"
    Then the HTTP status code should be "400"
    And I should receive the following JSON response:
    """
    {
      "errors": [
      "failed to delete options: failed to find dimension index: could not find dimension"
      ]
    }
    """
    Scenario: Delete Option Filter ID not found
    When I DELETE "/filters/NOT-EXISTS/dimensions/City/options"
    Then the HTTP status code should be "404"
    And I should receive the following JSON response:
    """
    {
      "errors": [
      "failed to delete option: filter not found"
      ]
    }
    """
