Feature: Filter Dimensions Private Endpoints Are Enabled

  Background:
    Given private endpoints are enabled
    And I am authorised
    And I am identified as "user@ons.gov.uk"
    And the following version document with dataset id "cantabular-example-1", edition "2021" and version "1" is available from dp-dataset-api:
    """
    {
      "alerts": [],
      "collection_id": "dfb-38b11d6c4b69493a41028d10de503aabed3728828e17e64914832d91e1f493c6",
        "is_based_on":{"@type": "cantabular_multivariate_table"},
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
          "href": "http://api.localhost:23200/v1/code-lists/siblings",
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
            "id": "siblings_3",
            "label": "Number of siblings (3 mappings)",
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
            "label": "City",
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
        "population_type": "dummy_data_households",
        "type": "multivariate"
      }
    ]
    """

    And Cantabular returns these dimensions for the dataset "dummy_data_households" and search term "hh_deprivation":
    """
    {
      "dataset": {
        "variables": {
            "edges": [
              {
                "node": {
                  "name": "hh_deprivation",
                  "label": "TEST_LABEL",
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
    """

  And Cantabular returns these categorisations for the dataset "dummy_data_households" and search term "hh_deprivation":
  """
  {
  "dataset": {
    "variables": {
      "search": {
        "edges": [
          {
            "node": {
              "categories": {
                "edges": [
                  {
                    "label": "hh_a",
                    "name": "hh_a"
                  },
                  {
                    "label": "hh_b",
                    "name": "hh_b"
                  }
                ]
              },
              "name": "hh_deprivation",
              "label": "hh_deprivation"
            }
          }
        ]
      }
    }
  }
  }
  """

  And the metadata api returns this response:
  """
  {
  "variable": "hh_deprivation"
  }
  """
  Scenario: Add a multivariate filter dimension with no options successfully

    When I POST "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
    """
    {
      "name": "hh_deprivation",
      "is_area_type": false,
      "filter_by_parent": ""
    }
    """
    Then the HTTP status code should be "201"
    And I should receive the following JSON response:
    """
    {
      "id": "hh_a",
      "name": "hh_a",
      "label": "TEST-LABEL",
      "links": {
        "filter": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1",
          "id": "94310d8d-72d6-492a-bc30-27584627edb1"
        },
        "options": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/hh_deprivation/options"
        },
        "self": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/hh_deprivation",
          "id": "hh_deprivation"
        }
      }
    }
    """

    Scenario: Add a multivariate dimension that does not exist
      When I POST "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
      """
      {
        "name":             "DOESNOTEXIST",
        "is_area_type":     false,
        "filter_by_parent": ""
      }
      """

      Then I should receive the following JSON response:
      """
      {
        "errors": [
            "failed to find dimension: DOESNOTEXIST"
        ]
      }
      """
      And the HTTP status code should be "404"

    Scenario: Add a multivariate dimension but metadata errors
    Given the metadata API returns an error
    When I POST "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
    """
    {
    "name": "hh_deprivation",
    "is_area_type": false,
    "filter_by_parent": ""
    }
    """
    Then the HTTP status code should be "500"
    And I should receive the following JSON response:
    """
    {
        "errors": [
            "failed to add dimension: internal server error"
        ]
    }
    """
