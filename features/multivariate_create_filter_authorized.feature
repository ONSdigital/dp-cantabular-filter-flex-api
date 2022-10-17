 Feature: Filters Private Endpoints Enabled

  Background:
    Given private endpoints are enabled
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

    And Cantabular returns dimensions for the dataset "dummy_data_households" for the following search terms:
      """
      {
      "responses": {
       "ladcd":     {
      "dataset": {
        "variables": {
            "edges": [
              {
                "node": {
                  "name": "LADCD",
                  "label": "Local Authority code",
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
    },
       "hh_deprivation_health":     {
      "dataset": {
        "variables": {
            "edges": [
              {
                "node": {
                  "name": "hh_deprivation_health",
                  "label": "Household deprived in the health and disability dimension (3 categories)",
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
    },
       "hh_deprivation":     {
      "dataset": {
        "variables": {
            "edges": [
              {
                "node": {
                  "name": "hh_deprivation",
                  "label": "Household deprivation (6 categories)",
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
      }
      """
  Scenario: Creating a new multivariate filter journey when authorized
    Given I am identified as "user@ons.gov.uk"

    And I am authorised

    When I POST "/filters"
    """
    {
      "dataset":{
        "id":      "cantabular-example-1",
        "edition": "2021",
        "version": 1
      },
      "population_type": "dummy_data_households",
      "dimensions": [
        {
          "name": "ladcd",
          "is_area_type": true,
          "filter_by_parent": ""
        },
        {
          "name": "hh_deprivation_health",
          "is_area_type": false,
          "filter_by_parent": ""
        },
        {
          "name": "hh_deprivation",
          "is_area_type": false,
          "filter_by_parent": ""
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
          "href": "http://mockhost:9999/datasets/cantabular-example-1/editions/2021/version/1",
          "id": "1"
        },
        "self": {
          "href": ":27100/filters/94310d8d-72d6-492a-bc30-27584627edb1"
        },
        "dimensions": {
          "href": ":27100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
        }
      },
      "instance_id": "c733977d-a2ca-4596-9cb1-08a6e724858b",
      "dataset": {
        "id":      "cantabular-example-1",
        "edition": "2021",
        "version": 1
      },
      "population_type": "dummy_data_households",
      "published":       true,
      "type":            "multivariate"
    }
    """

    And the HTTP status code should be "201"

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
          "href": ":27100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
        }
      },
      "etag": "e02ac5f3b1472258c821ff704f91e39957c19938",
      "instance_id": "c733977d-a2ca-4596-9cb1-08a6e724858b",
      "dimensions": [
        {
          "name": "LADCD",
          "label": "Local Authority code",
          "id": "LADCD",
          "is_area_type": true,
          "options": null
        },
        {
          "name": "hh_deprivation_health",
          "label": "Household deprived in the health and disability dimension (3 categories)",
          "id": "hh_deprivation_health",
          "is_area_type": false,
          "options": null
        },
        {
          "name": "hh_deprivation",
          "id": "hh_deprivation",
          "label": "Household deprivation (6 categories)",
          "is_area_type": false,
          "options": null
        }
      ],
      "dataset": {
        "id": "cantabular-example-1",
        "edition": "2021",
        "version": {
          "$numberInt": "1"
        }
      },
      "published": true,
      "population_type": "dummy_data_households",
      "type": "multivariate",
      "unique_timestamp": {
        "$timestamp": {
          "i": 1,
          "t": 1643200024.0
        }
      },
      "last_updated": {
        "$date": {
          "$numberLong": "1643200024783"
        }
      }
    }
    """
  Scenario: Creating a new multivariate with no dims when authorized
    Given I am identified as "user@ons.gov.uk"

    And I am authorised

    When I POST "/filters"
    """
    {
      "dataset":{
        "id":      "cantabular-example-1",
        "edition": "2021",
        "version": 1
      },
      "population_type": "dummy_data_households",
      "dimensions": []
    }
    """

    Then the HTTP status code should be "400"
    Then I should receive the following JSON response:
    """
    {
      "errors": [
          "failed to parse request: invalid request: missing/invalid field: 'dimensions' must contain at least 2 values"
      ]
    }
    """
  Scenario: Creating a new multivariate with a bad dim
    Given I am identified as "user@ons.gov.uk"

    And I am authorised

    When I POST "/filters"
    """
    {
      "dataset":{
        "id":      "cantabular-example-1",
        "edition": "2021",
        "version": 1
      },
      "population_type": "dummy_data_households",
      "dimensions": [
        {
          "name": "DOES NOT EXIST",
          "is_area_type": true,
          "filter_by_parent": ""
        },
        {
          "name": "hh_deprivation_health",
          "is_area_type": false,
          "filter_by_parent": ""
        },
        {
          "name": "hh_deprivation",
          "is_area_type": false,
          "filter_by_parent": ""
        }
      ]
    }
    """

    Then I should receive the following JSON response:
    """
    {
      "errors": [
          "failed to find dimension: DOES NOT EXIST"
      ]
    }
    """

    And the HTTP status code should be "404"

