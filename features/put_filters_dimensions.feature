Feature: Updating a filter's dimensions

  Background:
    Given private endpoints are not enabled
    And the following version document with dataset id "cantabular-example-1", edition "2021" and version "1" is available from dp-dataset-api:
    """
    {
      "alerts": [],
      "collection_id": "dfb-38b11d6c4b69493a41028d10de503aabed3728828e17e64914832d91e1f493c6",
      "is_based_on":{"@type": "cantabular_flexible_table"},
      "dimensions": [
        {
          "name": "geography",
          "id": "city",
          "label": "City",
          "links": {
            "code_list": {},
            "options": {},
            "version": {}
          },
          "href": "http://api.localhost:23200/v1/code-lists/city"
        },
        {
          "name": "geography",
          "id": "country",
          "label": "Country",
          "links": {
            "code_list": {},
            "options": {},
            "version": {}
          },
          "href": "http://api.localhost:23200/v1/code-lists/country"
        },
        {
          "id": "siblings_3",
          "name": "siblings",
          "label": "Number of siblings (3 mappings)",
          "links": {
            "code_list": {},
            "options": {},
            "version": {}
          },
          "href": "http://api.localhost:23200/v1/code-lists/siblings"
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
    And I have this filter with an ETag of "city":
    """
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
      "events": null,
      "instance_id": "c733977d-a2ca-4596-9cb1-08a6e724858b",
      "dimensions": [
        {
          "name": "siblings",
          "id": "siblings_3",
          "label": "Number of siblings (3 mappings)",
          "options": [],
          "dimension_url": "http://dimension.url/siblings",
          "is_area_type": false
        },
        {
          "name": "geography",
          "id": "city",
          "label": "City",
          "options": [
            "London"
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
    """
    And Cantabular returns these dimensions for the dataset "Example" and search term "country":
    """
    {
      "dataset": {
        "variables": {
          "search": {
            "edges": [
              {
                "node": {
                  "name": "country",
                  "label": "Country",
                  "mapFrom": [
                    {
                      "edges": [
                        {
                          "node": {
                            "label": "City",
                            "name": "city"
                          }
                        }
                      ],
                      "totalCount": 1
                    }
                  ]
                }
              }
            ]
          }
        }
      }
    }
    """

  Scenario: Replacing a filter dimension (returns the dimension)
    When I PUT "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography"
    """
    {
      "name": "geography",
      "id": "country",
      "is_area_type": true
    }
    """
    Then I should receive the following JSON response:
    """
    {
      "name": "geography",
      "id": "country",
      "label": "Country",
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
          "id": "country"
        }
      }
    }
    """
    And the HTTP status code should be "200"

    Scenario: Replacing a non-goegraphy filter dimension (returns an error)
    When I PUT "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography"
    """
    {
      "name": "siblings",
      "id": "siblings_3",
      "is_area_type": false
    }
    """
    Then I should receive an errors array
    And the HTTP status code should be "404"

  Scenario: Replacing a filter dimension (updates the ETag)
    When I PUT "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography"
    """
    {
      "name": "geography",
      "id": "country",
      "is_area_type": true
    }
    """
    Then the ETag is a hash of the filter "94310d8d-72d6-492a-bc30-27584627edb1"

    And the HTTP status code should be "200"


  # It would be good to also validate the options/area type bool were saved correctly, however the endpoint to
  # retrieve a dimension hasn't yet been implemented.
  Scenario: Replacing a filter dimension (updates the filter)
    When I PUT "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography"
    """
    {
      "name": "geography",
      "id": "country",
      "is_area_type": true
    }
    """
    And I GET "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
    Then I should receive the following JSON response:
    """
    {
      "items": [
        {
          "name": "geography",
          "id": "country",
          "label": "Country",
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
              "id": "country"
            }
          }
        },
        {
          "name": "siblings",
          "id": "siblings_3",
          "label": "Number of siblings (3 mappings)",
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
              "id": "siblings_3"
            }
          }
        }
      ],
      "limit": 20,
      "offset": 0,
      "count": 2,
      "total_count": 2
    }
    """
    And a document in collection "filters" with key "filter_id" value "94310d8d-72d6-492a-bc30-27584627edb1" should match:
    """
    {
      "_id": "94310d8d-72d6-492a-bc30-27584627edb1",
      "filter_id": "94310d8d-72d6-492a-bc30-27584627edb1",
      "links": {
        "version": {
          "href": "http://mockhost:9999/datasets/cantabular-example-1/editions/2021/version/1",
          "id": "1"
        },
        "self": {
          "href": ":27100/filters/94310d8d-72d6-492a-bc30-27584627edb1"
        },
        "dimensions": {
          "href": ""
        }
      },
      "etag": "b93d5717c0bb574bf989b5db3527c2f7d0a7abac",
      "instance_id": "c733977d-a2ca-4596-9cb1-08a6e724858b",
      "dataset": {
        "id": "cantabular-example-1",
        "edition": "2021",
        "version": {
          "$numberInt":"1"
        }
      },
      "dimensions": [
        {
          "name": "siblings",
          "id": "siblings_3",
          "label": "Number of siblings (3 mappings)",
          "options": [],
          "is_area_type": false
        },
        {
          "name": "geography",
          "id": "country",
          "label": "Country",
          "options": [],
          "is_area_type": true
        }
      ],
      "population_type": "Example",
      "published": true,
      "type": "flexible",
      "unique_timestamp":{
        "$timestamp":{
          "i":1,
          "t":1.643200024e+09
        }
      },
      "last_updated":{
        "$date":{
          "$numberLong":"1643200024783"
        }
      }
    }
    """

  Scenario: An invalid JSON body (results in a 400 status code)
    When I PUT "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography"
    """
    {
      "name": "country
    }
    """
    Then I should receive an errors array
    And the HTTP status code should be "400"

  Scenario: PUT dimension but 'is_area_type' is missing
    When I PUT "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/siblings"
    """
    {
      "name": "siblings",
      "id": "siblings",
      "options": ["4-7", "7+"]
    }
    """
    And I should receive the following JSON response:
    """
    {"errors":["failed to parse request: invalid request: missing field: [is_area_type]"]}
    """
    Then the HTTP status code should be "400"

  Scenario: An invalid JSON body (doesn't update the filter)
    When I PUT "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography"
    """
    {
      "name": "country
    }
    """
    And I GET "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
    Then I should receive the following JSON response:
    """
    {
      "items": [
        {
          "name": "geography",
          "id": "city",
          "label": "City",
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
              "id": "city"
            }
          }
        },
        {
          "name": "siblings",
          "id": "siblings_3",
          "label": "Number of siblings (3 mappings)",
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
              "id": "siblings_3"
            }
          }
        }
      ],
      "limit": 20,
      "offset": 0,
      "count": 2,
      "total_count": 2
    }
    """

  Scenario: An If-Match header is provided and doesn't match (returns a 409 status code)
    When I set the "If-Match" header to "stale"
    And I PUT "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography"
    """
    {
      "name": "geography",
      "id": "country",
      "is_area_type": true
    }
    """
    Then I should receive an errors array
    And the HTTP status code should be "409"

  Scenario: An If-Match header is provided and doesn't match (doesn't update the filter)
    When I set the "If-Match" header to "stale"
    And I PUT "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography"
    """
    {
      "name": "geography",
      "id": "country",
      "is_area_type": true
    }
    """
    And I GET "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
    Then I should receive the following JSON response:
    """
    {
      "items": [
        {
          "name": "geography",
          "id": "city",
          "label": "City",
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
              "id": "city"
            }
          }
        },
        {
          "name": "siblings",
          "id": "siblings_3",
          "label": "Number of siblings (3 mappings)",
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
              "id": "siblings_3"
            }
          }
        }
      ],
      "limit": 20,
      "offset": 0,
      "count": 2,
      "total_count": 2
    }
    """

  Scenario: The filter doesn't exist in the database
    When I PUT "/filters/not-found/dimensions/geography"
    """
    {
      "name": "geography",
      "id": "country",
      "is_area_type": true
    }
    """
    Then I should receive an errors array
    And the HTTP status code should be "400"

  Scenario: The dimension doesn't exist in the database
    When I PUT "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/sex"
    """
    {
      "name": "geography",
      "id": "country",
      "is_area_type": true
    }
    """
    Then I should receive an errors array
    Then the HTTP status code should be "404"

  Scenario: The dimension doesn't exist in Cantabular
    When I PUT "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography"
    """
    {
      "name": "geography",
      "id": "fake",
      "is_area_type": false
    }
    """
    Then I should receive an errors array
    Then the HTTP status code should be "404"

  Scenario: Searching Cantabular results in an error
    When Cantabular responds with an error
    And I PUT "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/geography"
    """
    {
      "id": "country",
      "name": "geography",
      "is_area_type": true
    }
    """
    Then I should receive an errors array
    Then the HTTP status code should be "500"
