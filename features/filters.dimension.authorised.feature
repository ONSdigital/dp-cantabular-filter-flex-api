Feature: Filter Outputs Private Endpoints Enabled

  Background:
    Given private endpoints are enabled
    And the following version document with dataset id "cantabular-example-1", edition "2021" and version "1" is available from dp-dataset-api:
      """
      {
        "alerts": [],
        "collection_id": "dfb-38b11d6c4b69493a41028d10de503aabed3728828e17e64914832d91e1f493c6",
        "dimensions": [
          {
            "label": "City",
            "links": {
              "code_list": {},
              "options": {},
              "version": {}
            },
            "href": "http://api.localhost:23200/v1/code-lists/city",
            "id": "city",
            "name": "geography"
          },
          {
            "label": "Number of siblings (3 mappings)",
            "links": {
              "code_list": {},
              "options": {},
              "version": {}
            },
            "href": "http://api.localhost:23200/v1/code-lists/siblings_3",
            "id": "siblings_3",
            "name": "siblings"
          }
        ],
        "edition": "2021",
        "id": "c733977d-a2ca-4596-9cb1-08a6e724858b",
        "links": {
          "dataset": {
            "href": "http://dp-dataset-api:22000/datasets/cantabular-example-1",
            "id": "cantabular-example-1"
          },
          "dimensions": {},
          "edition": {
            "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021",
            "id": "2021"
          },
          "self": {
            "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1"
          }
        },
        "release_date": "2021-11-19T00:00:00.000Z",
        "state": "published",
        "usage_notes": [],
        "version": 1
      }
      """
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
            "href": ":27100/filters/94310d8d-72d6-492a-bc30-27584627edb1"
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
            "name": "geography",
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
      }
    ]
    """

  Scenario: Get filter dimensions successfully
    Given I am identified as "user@ons.gov.uk"
    And I am authorised
    When I GET "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
    Then I should receive the following JSON response:
    """
    {
      "items": [
        {
          "name": "geography",
          "links": {
            "filter": {
              "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1",
              "id": "94310d8d-72d6-492a-bc30-27584627edb1"
            },
            "options": {
              "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography/options"
            },
            "self": {
              "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography",
              "id": "geography"
            }
          }
        },
        {
          "name": "siblings",
          "links": {
            "filter": {
              "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1",
              "id": "94310d8d-72d6-492a-bc30-27584627edb1"
            },
            "options": {
              "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/siblings/options"
            },
            "self": {
              "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/siblings",
              "id": "siblings"
            }
          }
        }
      ],
      "count": 2,
      "offset": 0,
      "limit": 20,
      "total_count": 2
    }
    """
    And the HTTP status code should be "200"

  Scenario: Get a specific filter dimension successfully
    Given I am identified as "user@ons.gov.uk"
    And I am authorised
    When I GET "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography"
    Then I should receive the following JSON response:
    """
    {
      "name": "geography",
      "is_area_type": true,
      "links": {
        "filter": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1",
          "id": "94310d8d-72d6-492a-bc30-27584627edb1"
        },
        "options": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography/options"
        },
        "self": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography",
          "id": "geography"
        }
      }
    }
    """
    And the HTTP status code should be "200"

  Scenario: Get a specific filter dimension when the filter is not present
    Given I am identified as "user@ons.gov.uk"
    And I am authorised
    When I GET "/filters/00000000-0000-0000-0000-000000000000/dimensions/geography"
    Then the HTTP status code should be "400"

  Scenario: Get a specific filter dimension when the dimension is not present
    Given I am identified as "user@ons.gov.uk"
    And I am authorised
    When I GET "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/other"
    Then the HTTP status code should be "404"

  Scenario: Add a filter dimension successfully
    Given I am identified as "user@ons.gov.uk"
    And I am authorised

    When I POST "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
    """
    {
      "name": "siblings",
      "is_area_type": false,
      "options": ["4-7", "7+"]
    }
    """
    Then the HTTP status code should be "201"
    Then I should receive the following JSON response:
    """
    {
      "name": "siblings",
      "links": {
        "filter": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1",
          "id": "94310d8d-72d6-492a-bc30-27584627edb1"
        },
        "options": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/siblings/options"
        },
        "self": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/siblings",
          "id": "siblings"
        }
      }
    }
    """
    And an ETag is returned


  Scenario: Add a filter dimension with no options successfully
    Given I am identified as "user@ons.gov.uk"
    And I am authorised
    When I POST "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
    """
    {
      "name": "siblings",
      "is_area_type": false
    }
    """

    Then the HTTP status code should be "201"

    And I should receive the following JSON response:
    """
    {
      "name": "siblings",
      "links": {
        "filter": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1",
          "id": "94310d8d-72d6-492a-bc30-27584627edb1"
        },
        "options": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/siblings/options"
        },
        "self": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/siblings",
          "id": "siblings"
        }
      }
    }
    """
    And an ETag is returned

  Scenario: Add a filter dimension but request is malformed
    Given I am identified as "user@ons.gov.uk"
    And I am authorised
    When I POST "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
    """
    Not valid JSON
    """
    Then the HTTP status code should be "400"

  Scenario: Add a filter dimension but the filter is not found
    Given I am identified as "user@ons.gov.uk"
    And I am authorised
    When I POST "/filters/not-valid-filter-id/dimensions"
    """
    {
      "name": "siblings",
      "links": {
        "filter": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1",
          "id": "94310d8d-72d6-492a-bc30-27584627edb1"
        },
        "options": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/siblings/options"
        },
        "self": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/siblings",
          "id": "siblings"
        }
      }
    }
    """
    Then the HTTP status code should be "404"

  Scenario: Add a filter dimension but the filter's dimensions have been modified since I retrieved it
    Given I am identified as "user@ons.gov.uk"
    And I am authorised
    And I provide If-Match header "invalid-etag"
    When I POST "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
    """
    {
      "name": "siblings",
      "links": {
        "filter": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1",
          "id": "94310d8d-72d6-492a-bc30-27584627edb1"
        },
        "options": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/siblings/options"
        },
        "self": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/siblings",
          "id": "siblings"
        }
      }
    }
    """
    Then the HTTP status code should be "409"

  Scenario: Add a filter dimension but something goes wrong on the server
    Given I am identified as "user@ons.gov.uk"
    And I am authorised
    Given the client for the dataset API failed and is returning errors
    When I POST "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
    """
    {
      "name": "siblings",
      "links": {
        "filter": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1",
          "id": "94310d8d-72d6-492a-bc30-27584627edb1"
        },
        "options": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/siblings/options"
        },
        "self": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/siblings",
          "id": "siblings"
        }
      }
    }
    """
    Then the HTTP status code should be "500"

  Scenario: Get paginated filter dimensions successfully (limit)
    Given I am identified as "user@ons.gov.uk"
    And I am authorised
    When I GET "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions?limit=1"
    Then I should receive the following JSON response:
    """
    {
      "items": [
        {
          "name": "geography",
          "links": {
            "filter": {
              "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1",
              "id": "94310d8d-72d6-492a-bc30-27584627edb1"
            },
            "options": {
              "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography/options"
            },
            "self": {
              "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography",
              "id": "geography"
            }
          }
        }
      ],
      "count": 1,
      "offset": 0,
      "limit": 1,
      "total_count": 2
    }
    """
    And the HTTP status code should be "200"

  Scenario: Get paginated filter dimensions successfully (offset)
    Given I am identified as "user@ons.gov.uk"
    And I am authorised
    When I GET "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions?offset=1"
    Then I should receive the following JSON response:
    """
    {
      "items": [
        {
          "name": "siblings",
          "links": {
            "filter": {
              "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1",
              "id": "94310d8d-72d6-492a-bc30-27584627edb1"
            },
            "options": {
              "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/siblings/options"
            },
            "self": {
              "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/siblings",
              "id": "siblings"
            }
          }
        }
      ],
      "count": 1,
      "offset": 1,
      "limit": 20,
      "total_count": 2
    }
    """
    And the HTTP status code should be "200"

  Scenario: Get paginated filter dimensions successfully (0 limit)
    Given I am identified as "user@ons.gov.uk"
    And I am authorised
    When I GET "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions?limit=0"
    Then I should receive the following JSON response:
    """
    {
      "items": [],
      "count": 0,
      "offset": 0,
      "limit": 0,
      "total_count": 2
    }
    """
    And the HTTP status code should be "200"

  Scenario: Get filter dimensions unsuccessfully
    Given I am identified as "user@ons.gov.uk"
    And I am authorised
    When I GET "/filters/94310d8d-72d6-492a-03cb-27584627edb1/dimensions"
    Then I should receive the following JSON response:
    """
    {
      "errors": ["failed to get filter dimensions"]
    }
    """
    And the HTTP status code should be "404"

  Scenario: Get filter dimensions unsuccessfully (cannot parse offset)
    Given I am identified as "user@ons.gov.uk"
    And I am authorised
    When I GET "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions?offset=dog"
    Then I should receive the following JSON response:
    """
    {
      "errors": ["invalid parameter offset"]
    }
    """
    And the HTTP status code should be "400"

  Scenario: Get filter dimensions unsuccessfully (cannot parse limit)
    Given I am identified as "user@ons.gov.uk"
    And I am authorised
    When I GET "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions?limit=dog"
    Then I should receive the following JSON response:
    """
    {
      "errors": ["invalid parameter limit"]
    }
    """
    And the HTTP status code should be "400"

  Scenario: Get filter dimensions unsuccessfully (invalid limit, too large)
    Given I am identified as "user@ons.gov.uk"
    And I am authorised
    Given the maximum pagination limit is set to 100
    When I GET "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions?limit=101"
    Then I should receive the following JSON response:
    """
    {
      "errors": ["limit cannot be larger than 100"]
    }
    """
    And the HTTP status code should be "400"

  Scenario: Get filter dimensions unsuccessfully (invalid limit, less than 0)
    Given I am identified as "user@ons.gov.uk"
    And I am authorised
    When I GET "/filters/94310d8d-72d6-492a-03cb-27584627edb1/dimensions?limit=-1"
    Then I should receive the following JSON response:
    """
    {
      "errors": ["limit cannot be less than 0"]
    }
    """
    And the HTTP status code should be "400"

  Scenario: Get filter dimensions unsuccessfully (invalid offset, less than 0)
    Given I am identified as "user@ons.gov.uk"
    And I am authorised
    When I GET "/filters/94310d8d-72d6-492a-03cb-27584627edb1/dimensions?offset=-1"
    Then I should receive the following JSON response:
    """
    {
      "errors": ["offset cannot be less than 0"]
    }
    """
    And the HTTP status code should be "400"

