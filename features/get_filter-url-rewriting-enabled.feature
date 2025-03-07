Feature: Get Filter Private Endpoints Not Enabled with URL rewriting enabled

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
            "href": "http://localhost:27100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
          },
          "self": {
            "href": "http://localhost:27100/filters/94310d8d-72d6-492a-bc30-27584627edb1"
          }
        },
        "instance_id": "c733977d-a2ca-4596-9cb1-08a6e724858b",
        "dimensions": [
          {
            "name": "silbings",
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
          "version": 1,
          "lowest_geography": "lowest-geography",
          "release_date": "2021-11-19T00:00:00.000Z",
          "title": "cantabular-example-1"
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
            "href": "http://localhost:27100/filters/83210d8d-72d6-492a-bc30-27584627abc2/dimensions"
          },
          "self": {
            "href": "http://localhost:27100/filters/83210d8d-72d6-492a-bc30-27584627abc2"
          }
        },
        "instance_id": "c733977d-a2ca-4596-9cb1-08a6e724858b",
        "dimensions": [
          {
            "name": "siblings",
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
          "version": 1,
          "release_date": "2021-11-19T00:00:00.000Z",
          "title": "cantabular-example-unpublished"
        },
        "published": false,
        "population_type": "Example",
        "type": "flexible"
      }
    ]
    """

  Scenario: Get filter successfully
    And I set the "X-Forwarded-API-Host" header to "api.example.com"
    And I set the "X-Forwarded-Path-Prefix" header to "v1"
    And URL rewriting is enabled
    When I GET "/filters/94310d8d-72d6-492a-bc30-27584627edb1"
    
    Then I should receive the following JSON response:
    """
    {
      "filter_id": "94310d8d-72d6-492a-bc30-27584627edb1",
      "instance_id": "c733977d-a2ca-4596-9cb1-08a6e724858b",
      "links": {
        "version": {
          "href": "https://api.example.com/v1/datasets/cantabular-example-1/editions/2021/versions/1",
          "id": "1"
        },
        "self": {
          "href": "https://api.example.com/v1/filters/94310d8d-72d6-492a-bc30-27584627edb1"
        },
        "dimensions": {
          "href": "https://api.example.com/v1/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
        }
      },
      "dataset": {
        "id": "cantabular-example-1",
        "edition": "2021",
        "version": 1,
        "lowest_geography": "lowest-geography",
        "release_date": "2021-11-19T00:00:00.000Z",
        "title": "cantabular-example-1"
      },
      "published": true,
      "population_type": "Example",
      "type": "flexible",
      "custom": false
    }
    """
    And the HTTP status code should be "200"
