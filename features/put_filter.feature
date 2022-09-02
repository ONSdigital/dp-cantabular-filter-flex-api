Feature: Put Filter Private Endpoints Not Enabled

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

  Scenario: PUT filter successfully
    When I PUT "/filters/94310d8d-72d6-492a-bc30-27584627edb1"
    """
    """
    Then I should receive the following JSON response:
    """
    {
      "events": [
        {
          "timestamp": "2016-07-17T08:38:25.316+000",
          "name": "cantabular-export-start"
        }
      ],
      "dataset": {
        "id": "string",
        "edition": "string",
        "version": 0
      },
      "population_type": "string"
    }
    """
    And the HTTP status code should be "200"
