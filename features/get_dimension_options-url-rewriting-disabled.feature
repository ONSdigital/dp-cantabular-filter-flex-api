Feature: Get Filter Dimension Options Private Endpoints
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
    Scenario: Filter Dimension Option Found
      And I set the "X-Forwarded-Proto" header to "https"
      And I set the "X-Forwarded-Host" header to "api.example.com"
      And I set the "X-Forwarded-Path-Prefix" header to "v1"
      And URL rewriting is enabled
      When I GET "/filters/83210d8d-72d6-492a-bc30-27584627abc2/dimensions/City/options?limit=10&offset=0"
      Then the HTTP status code should be "200"
      Then I should receive the following JSON response:
      """
      {
      "items": [
        {
          "option": "Cardiff",
          "self": {
            "href": "http://localhost:22100/filters/83210d8d-72d6-492a-bc30-27584627abc2/dimensions/City/options",
            "id": "Cardiff"
          },
          "filter": {
            "href": "http://localhost:22100/filters/83210d8d-72d6-492a-bc30-27584627abc2",
            "id": "83210d8d-72d6-492a-bc30-27584627abc2"
          },
          "Dimension": {
            "href": "http://localhost:22100/filters/83210d8d-72d6-492a-bc30-27584627abc2/dimensions/City",
            "id": "City"
          }
        },
        {
          "option": "London",
          "self": {
            "href": "http://localhost:22100/filters/83210d8d-72d6-492a-bc30-27584627abc2/dimensions/City/options",
            "id": "London"
          },
          "filter": {
            "href": "http://localhost:22100/filters/83210d8d-72d6-492a-bc30-27584627abc2",
            "id": "83210d8d-72d6-492a-bc30-27584627abc2"
          },
          "Dimension": {
            "href": "http://localhost:22100/filters/83210d8d-72d6-492a-bc30-27584627abc2/dimensions/City",
            "id": "City"
          }
        },
        {
          "option": "Swansea",
          "self": {
            "href": "http://localhost:22100/filters/83210d8d-72d6-492a-bc30-27584627abc2/dimensions/City/options",
            "id": "Swansea"
          },
          "filter": {
            "href": "http://localhost:22100/filters/83210d8d-72d6-492a-bc30-27584627abc2",
            "id": "83210d8d-72d6-492a-bc30-27584627abc2"
          },
          "Dimension": {
            "href": "http://localhost:22100/filters/83210d8d-72d6-492a-bc30-27584627abc2/dimensions/City",
            "id": "City"
          }
        }
      ],
      "limit": 10,
      "offset": 0,
      "count": 3,
      "total_count": 3
      }
      """
    Scenario: Filter Dimension Zero Page Limit
      In the case of zero page limit, a reasonable page limit is introduced.

      When I GET "/filters/83210d8d-72d6-492a-bc30-27584627abc2/dimensions/City/options?limit=0&offset=0"
      Then the HTTP status code should be "200"
      And I should receive the following JSON response:
      """
      {
      "items": [
        {
          "option": "Cardiff",
          "self": {
            "href": "http://localhost:22100/filters/83210d8d-72d6-492a-bc30-27584627abc2/dimensions/City/options",
            "id": "Cardiff"
          },
          "filter": {
            "href": "http://localhost:22100/filters/83210d8d-72d6-492a-bc30-27584627abc2",
            "id": "83210d8d-72d6-492a-bc30-27584627abc2"
          },
          "Dimension": {
            "href": "http://localhost:22100/filters/83210d8d-72d6-492a-bc30-27584627abc2/dimensions/City",
            "id": "City"
          }
        },
        {
          "option": "London",
          "self": {
            "href": "http://localhost:22100/filters/83210d8d-72d6-492a-bc30-27584627abc2/dimensions/City/options",
            "id": "London"
          },
          "filter": {
            "href": "http://localhost:22100/filters/83210d8d-72d6-492a-bc30-27584627abc2",
            "id": "83210d8d-72d6-492a-bc30-27584627abc2"
          },
          "Dimension": {
            "href": "http://localhost:22100/filters/83210d8d-72d6-492a-bc30-27584627abc2/dimensions/City",
            "id": "City"
          }
        },
        {
          "option": "Swansea",
          "self": {
            "href": "http://localhost:22100/filters/83210d8d-72d6-492a-bc30-27584627abc2/dimensions/City/options",
            "id": "Swansea"
          },
          "filter": {
            "href": "http://localhost:22100/filters/83210d8d-72d6-492a-bc30-27584627abc2",
            "id": "83210d8d-72d6-492a-bc30-27584627abc2"
          },
          "Dimension": {
            "href": "http://localhost:22100/filters/83210d8d-72d6-492a-bc30-27584627abc2/dimensions/City",
            "id": "City"
          }
        }
      ],
      "limit": 20,
      "offset": 0,
      "count": 3,
      "total_count": 3
      }
      """
   Scenario: Filter Not Found
     When I GET "/filters/notExists/dimensions/someDimension/options?page_limit=0&offset=0"
     Then the HTTP status code should be "404"
     And I should receive the following JSON response:
     """
     {
     "errors": [
         "failed to get filter dimension option"
     ]
     }
     """
   Scenario: Dimension Not Found
      When I GET "/filters/83210d8d-72d6-492a-bc30-27584627abc2/dimensions/not-exist/options"
      Then the HTTP status code should be "400"
      And I should receive the following JSON response:
      """
      {
      "errors": [
         "failed to get filter dimension option"
      ]
      }
      """
   Scenario: Mongo Client Errors
     Given Mongo datastore is failing
     When I GET "/filters/test/dimensions/test/options"
     Then the HTTP status code should be "500"
     And I should receive the following JSON response:
     """
     {
     "errors": [
         "failed to get filter dimension option"
     ]
     }
     """
