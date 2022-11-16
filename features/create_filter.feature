Feature: Filters Private Endpoints Not Enabled

  Background:
    Given private endpoints are not enabled

    And the following dataset document with dataset id "cantabular-example-1" is available from dp-dataset-api:
      """
      {
        "id":"cantabular-example-1",
        "links":{
          "self":{
            "href":"http://hostname/datasets/cantabular-flexible-table-component-test",
            "id":"cantabular-flexible-table-component-test"
          }
        },
        "state":"published",
        "title":"cantabular-example-1",
        "release_date": "2021-11-19T00:00:00.000Z",
        "type":"cantabular_flexible_table",
        "is_based_on":{
          "@type":"cantabular_flexible_table",
          "@id":"Example"
        }
      }
      """

    And the following version document with dataset id "cantabular-example-1", edition "2021" and version "1" is available from dp-dataset-api:
      """
      {
        "alerts": [],
        "collection_id": "dfb-38b11d6c4b69493a41028d10de503aabed3728828e17e64914832d91e1f493c6",
        "is_based_on":{"@type": "cantabular_flexible_table"},
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
            "name": "city"
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
            "name": "siblings_3"
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

    And the following version document with dataset id "cantabular_table_example", edition "2021" and version "1" is available from dp-dataset-api:
    """
    {
        "alerts": [],
        "collection_id": "dfb-38b11d6c4b69493a41028d10de503aabed3728828e17e64914832d91e1f493c6",
        "is_based_on":{"@type": "cantabular_table"},
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
            "name": "city"
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
            "name": "siblings_3"
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

    And the following version document with dataset id "cantabular-example-unpublished", edition "2021" and version "1" is available from dp-dataset-api:
      """
      {
        "alerts": [],
        "collection_id": "dfb-38b11d6c4b69493a41028d10de503aabed3728828e17e64914832d91e1f493c6",
        "is_based_on":{"@type": "flexible"},
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
            "name": "city"
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
            "name": "siblings_3"
          }
        ],
        "edition": "2021",
        "id": "c733977d-a2ca-4596-9cb1-08a6e724858b",
        "links": {
          "dataset": {
            "href": "http://dp-dataset-api:22000/datasets/cantabular-example-unpublished",
            "id": "cantabular-example-1"
          },
          "dimensions": {},
          "edition": {
            "href": "http://localhost:22000/datasets/cantabular-example-unpublished/editions/2021",
            "id": "2021"
          },
          "self": {
            "href": "http://localhost:22000/datasets/cantabular-example-unpublished/editions/2021/versions/1"
          }
        },
        "release_date": "2021-11-19T00:00:00.000Z",
        "state": "associated",
        "usage_notes": [],
        "version": 1
      }
      """
    And Cantabular returns dimensions for the dataset "dummy_data_households" for the following search terms:
      """
      {
        "responses": {
          "city": {
            "dataset": {
              "variables": {
                "edges": [
                  {
                    "node": {
                      "name": "city",
                      "label": "City",
                      "mapFrom": [
                        {
                          "edges": [
                            {
                              "node": {
                                "label": "OA",
                                "name": "oa"
                              }
                            }
                          ]
                        }
                      ],
                      "totalCount": 1
                    }
                  }
                ]
              }
            }
          },
          "region": {
            "dataset": {
              "variables": {
                "edges": [
                  {
                    "node": {
                      "name": "region",
                      "label": "Region",
                      "mapFrom": [
                        {
                          "edges": [
                            {
                              "node": {
                                "label": "OA",
                                "name": "oa"
                              }
                            }
                          ]
                        }
                      ],
                      "totalCount": 1
                    }
                  }
                ]
              }
            }
          }
        }
      }
      """

  Scenario: Creating a new filter happy

    When I POST "/filters"
    """
    {
      "dataset":{
          "id":      "cantabular-example-1",
          "edition": "2021",
          "version": 1,
          "lowest_geography": "lowest-geography"
      },
      "population_type": "Example",
      "dimensions": [
        {
          "name": "siblings_3",
          "options": [
            "0-3",
            "4-7",
            "7+"
          ],
          "is_area_type": false
        },{
          "name": "city",
          "options": [
            "Cardiff",
            "London",
            "Swansea"
          ],
          "is_area_type": true
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
        "id": "cantabular-example-1",
        "edition": "2021",
        "version": 1,
        "lowest_geography": "lowest-geography",
        "release_date": "2021-11-19T00:00:00.000Z",
        "title": "cantabular-example-1"
      },
      "population_type": "Example",
      "published": true,
      "type": "flexible"
    }
    """

    Scenario: Creating a new filter happy with non-default geography

    When I POST "/filters"
    """
    {
      "dataset":{
          "id":      "cantabular-example-1",
          "edition": "2021",
          "version": 1,
          "lowest_geography": "lowest-geography"
      },
      "population_type": "Example",
      "dimensions": [
        {
          "name": "siblings_3",
          "options": [
            "0-3",
            "4-7",
            "7+"
          ],
          "is_area_type": false
        },{
          "name": "region",
          "is_area_type": true,
          "options": [
            "0"
          ]
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
        "id": "cantabular-example-1",
        "edition": "2021",
        "version": 1,
        "lowest_geography": "lowest-geography",
        "release_date": "2021-11-19T00:00:00.000Z",
        "title": "cantabular-example-1"
      },
      "population_type": "Example",
      "published": true,
      "type": "flexible"
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
      "etag": "54d59f1041887a77575a660098dc33ecfeaab05f",
      "instance_id": "c733977d-a2ca-4596-9cb1-08a6e724858b",
      "dataset": {
        "id": "cantabular-example-1",
        "edition": "2021",
        "version": {
          "$numberInt":"1"
        },
        "lowest_geography": "lowest-geography",
        "release_date": "2021-11-19T00:00:00.000Z",
        "title": "cantabular-example-1"
      },
      "dimensions": [
        {
          "name": "region",
          "id": "region",
          "label": "Region",
          "is_area_type": true,
          "options": [
            "0"
          ]
        },
        {
          "name": "siblings_3",
          "id": "siblings_3",
          "label": "Number of siblings (3 mappings)",
          "options": [
            "0-3",
            "4-7",
            "7+"
          ],
          "is_area_type": false
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

  Scenario: Creating a new filter unauthenticated on unpublished version

    When I POST "/filters"
    """
    {
      "dataset":{
          "id":      "cantabular-example-unpublished",
          "edition": "2021",
          "version": 1,
          "lowest_geography": "lowest-geography"
      },
      "population_type": "Example",
      "dimensions": [
        {
          "name": "siblings_3",
          "options": [
            "0-3",
            "4-7",
            "7+"
          ],
          "is_area_type": false
        },{
          "name": "city",
          "options": [
            "Cardiff",
            "London",
            "Swansea"
          ],
          "is_area_type": true
        }
      ]
    }
    """

    Then I should receive the following JSON response:
    """
    {
      "errors": [
        "dataset not found"
      ]
    }
    """

    And the HTTP status code should be "404"

  Scenario: Creating a new filter bad request body

    When I POST "/filters"
    """
    """

    Then I should receive the following JSON response:
    """
    {
      "errors": [
        "badly formed request body: unexpected end of JSON input"
      ]
    }
    """

    And the HTTP status code should be "400"

  Scenario: Creating a new single dimension filter

    When I POST "/filters"
    """
    {
      "dataset":{
        "id":      "cantabular-example-1",
        "edition": "2021",
        "version": 1,
        "lowest_geography": "lowest-geography"
      },
      "population_type": "Example",
      "dimensions": [
        {
          "id": "siblings_3",
          "name": "siblings_3",
          "options": [
            "0-3",
            "4-7",
            "7+"
          ],
          "is_area_type": false
        }
      ]
    }
    """

    Then I should receive an errors array

    And the HTTP status code should be "400"

  Scenario: Creating a new filter (invalid request, passing label)

    When I POST "/filters"
    """
    {
      "dataset":{
        "id":      "cantabular-example-1",
        "edition": "2021",
        "version": 1,
        "lowest_geography": "lowest-geography"
      },
      "population_type": "Example",
      "dimensions": [
        {
          "label": "Number of Siblings (3 mappings)",
          "name": "siblings_3",
          "options": [
            "0-3",
            "4-7",
            "7+"
          ],
          "is_area_type": false
        },{
          "label": "City",
          "name": "city",
          "options": [
            "Cardiff",
            "London",
            "Swansea"
          ],
          "is_area_type": true
        }
      ]
    }
    """

    Then I should receive an errors array

    And the HTTP status code should be "400"

#  Scenario: Creating a new filter but 'is_area_type' missing from dimension
#    When I POST "/filters"
#    """
#    {
#      "dataset":{
#          "id":      "cantabular-example-1",
#          "edition": "2021",
#          "version": 1
#      },
#      "population_type": "Example",
#      "dimensions": [
#        {
#          "name": "siblings_3",
#          "options": [
#            "0-3",
#            "4-7",
#            "7+"
#          ]
#        },{
#          "name": "city",
#          "options": [
#            "Cardiff",
#            "London",
#            "Swansea"
#          ],
#          "is_area_type": true
#        }
#      ]
#    }
#    """
#
#    Then I should receive the following JSON response:
#    """
#    {"errors":["missing field: ['is_area_type']"]}
#    """
#
#    And the HTTP status code should be "404"

  Scenario: Creating a new filter but multiple geography dimensions selected
    When I POST "/filters"
    """
    {
      "dataset":{
          "id":      "cantabular-example-1",
          "edition": "2021",
          "version": 1,
          "lowest_geography": "lowest-geography"
      },
      "population_type": "Example",
      "dimensions": [
        {
          "name": "region",
          "options": [
            "South East",
            "North West",
            "South"
          ],
          "is_area_type": true
        },{
          "name": "city",
          "options": [
            "Cardiff",
            "London",
            "Swansea"
          ],
          "is_area_type": true
        }
      ]
    }
    """

    Then I should receive the following JSON response:
    """
    {"errors":["failed to validate dimensions: multiple geography dimensions not permitted"]}
    """

    And the HTTP status code should be "400"

    Scenario: Creating a new invalid request (duplicate dimensions)

    When I POST "/filters"
    """
    {
      "dataset":{
          "id":      "cantabular-example-1",
          "edition": "2021",
          "version": 1,
          "lowest_geography": "lowest-geography"
       },
      "population_type": "Example",
      "dimensions": [
        {
          "name": "siblings_3",
          "options": [
            "0-3",
            "4-7",
            "7+"
          ],
          "is_area_type": false
        },
        {
          "name": "siblings_3",
          "options": [
            "0-3",
            "4-7",
            "7+"
          ],
          "is_area_type": false
        }
      ]
    }
    """

    Then I should receive the following JSON response:
    """
    {
      "errors": [
        "duplicate dimension chosen: siblings_3"
      ]
    }
    """

    And the HTTP status code should be "400"

  Scenario: Do not create a filter based on cantabular blob
    When I POST "/filters"
    """
      {
      "dataset":{
          "id":      "cantabular_table_example",
          "edition": "2021",
          "version": 1,
          "lowest_geography": "lowest-geography"
      },
      "population_type": "Example",
      "dimensions": [
        {
          "name": "siblings_3",
          "options": [
            "0-3",
            "4-7",
            "7+"
          ],
          "is_area_type": false
        },{
          "name": "city",
          "options": [
            "Cardiff",
            "London",
            "Swansea"
          ],
          "is_area_type": true
        }
      ]
    }
    """
    Then the HTTP status code should be "400"
    And I should receive the following JSON response:
    """
    {
      "errors": [
        "dataset is of invalid type"
      ]
    }
    """
