Feature: Get Dataset Observations with URL rewriting

  Background:
    Given private endpoints are not enabled

    And the Cantabular service is a mocked extended Cantabular server

    And the following dataset document with dataset id "cantabular-flexible-table-test" is available from dp-dataset-api:
    """
    {
      "id":"cantabular-flexible-table-test",
      "links":{
        "self":{
          "href":"http://hostname/datasets/cantabular-flexible-table-test",
          "id":"cantabular-flexible-table-test"
        }
      },
      "state":"published",
      "title":"cantabular-flexible-table-test",
      "type":"cantabular_flexible_table",
      "is_based_on":{
        "@type":"cantabular_flexible_table",
        "@id":"Example"
      }
    }
    """

    And the following dataset document with dataset id "cantabular-multivariate-table-test" is available from dp-dataset-api:
    """
    {
      "id":"cantabular-multivariate-table-test",
      "links":{
        "self":{
          "href":"http://hostname/datasets/cantabular-multivariate-table-test",
          "id":"cantabular-multivariate-table-test"
        }
      },
      "state":"published",
      "title":"cantabular-multivariate-table-test",
      "type":"cantabular_multivariate_table",
      "is_based_on":{
        "@type":"cantabular_multivariate_table",
        "@id":"Example"
      }
    }
    """

    And the following version document with dataset id "cantabular-flexible-table-test", edition "latest" and version "1" is available from dp-dataset-api:
    """
    {
      "edition":"latest",
      "dimensions":[
        {
          "id":"city",
          "name":"city",
          "links":{
            "code_list":{
              "href":"http://hostname/code-lists/city",
              "id":"city"
            }
          },
          "description":"",
          "label":"City",
          "variable":"city"
        },
        {
          "id":"sex",
          "name":"sex",
          "links":{
            "code_list":{
              "href":"http://hostname/code-lists/sex",
              "id":"sex"
            }
          },
          "description":"",
          "label":"Sex",
          "variable":"sex"
        },
        {
          "id":"siblings_3",
          "name":"siblings_3",
          "links":{
            "code_list":{
              "href":"http://hostname/code-lists/siblings_3",
              "id":"siblings_3"
            }
          },
          "description":"",
          "label":"Number of siblings (3 mappings)",
          "variable":"siblings_3"
        }
      ],
      "id":"cantabular-flexible-table-testUUID",
      "links":{
        "dataset":{
          "href":"http://hostname/datasets/cantabular-flexible-table-test",
          "id":"cantabular-flexible-table-test"
        },
        "self":{
          "href":"http://hostname/datasets/cantabular-flexible-table-test/editions/latest/versions/1",
          "id":"1"
        },
        "version":{
          "href":"http://hostname/datasets/cantabular-flexible-table-test/editions/latest/versions/1",
          "id":"1"
        }
      },
      "state":"published",
      "version":1,
      "is_based_on":{
        "@type":"cantabular_flexible_table",
        "@id":"Example"
      }
    }
    """

    And the following version document with dataset id "cantabular-multivariate-table-test", edition "latest" and version "1" is available from dp-dataset-api:
    """
    {
      "edition":"latest",
      "dimensions":[
        {
          "id":"city",
          "name":"city",
          "links":{
            "code_list":{
              "href":"http://hostname/code-lists/city",
              "id":"city"
            }
          },
          "description":"",
          "label":"City",
          "variable":"city"
        },
        {
          "id":"sex",
          "name":"sex",
          "links":{
            "code_list":{
              "href":"http://hostname/code-lists/sex",
              "id":"sex"
            }
          },
          "description":"",
          "label":"Sex",
          "variable":"sex"
        }
      ],
      "id":"cantabular-multivariate-table-testUUID",
      "links":{
        "dataset":{
          "href":"http://hostname/datasets/cantabular-multivariate-table-test",
          "id":"cantabular-flexible-table-test"
        },
        "self":{
          "href":"http://hostname/datasets/cantabular-multivariate-table-test/editions/latest/versions/1",
          "id":"1"
        },
        "version":{
          "href":"http://hostname/datasets/cantabular-multivariate-table-test/editions/latest/versions/1",
          "id":"1"
        }
      },
      "state":"published",
      "version":1,
      "is_based_on":{
        "@type":"cantabular_multivariate_table",
        "@id":"Example"
      }
    }
    """

    And the following dimensions document for dataset id "cantabular-multivariate-table-test", edition "latest" and version "1" is available from dp-dataset-api:
    """
    {
      "items":[
        {
          "id":"city",
          "name":"city",
          "links":{
            "code_list":{
              "href":"http://hostname/code-lists/city",
              "id":"city"
            }
          },
          "description":"",
          "label":"City",
          "variable":"city"
        },
        {
          "id":"sex",
          "name":"sex",
          "links":{
            "code_list":{
              "href":"http://hostname/code-lists/sex",
              "id":"sex"
            }
          },
          "description":"",
          "label":"Sex",
          "variable":"sex"
        }
      ]
    }
    """

    And the following dimensions document for dataset id "cantabular-flexible-table-test", edition "latest" and version "1" is available from dp-dataset-api:
    """
    {
      "items":[
        {
          "id":"city",
          "name":"city",
          "links":{
            "code_list":{
              "href":"http://hostname/code-lists/city",
              "id":"city"
            }
          },
          "description":"",
          "label":"City",
          "variable":"city"
        },
        {
          "id":"sex",
          "name":"sex",
          "links":{
            "code_list":{
              "href":"http://hostname/code-lists/sex",
              "id":"sex"
            }
          },
          "description":"",
          "label":"Sex",
          "variable":"sex"
        },
        {
          "id":"siblings_3",
          "name":"siblings_3",
          "links":{
            "code_list":{
              "href":"http://hostname/code-lists/siblings_3",
              "id":"siblings_3"
            }
          },
          "description":"",
          "label":"Number of siblings (3 mappings)",
          "variable":"siblings_3"
        }
      ]
    }
    """

    And the following options document for dataset id "cantabular-flexible-table-test", edition "latest", version "1" and dimension "city" is available from dp-dataset-api:
    """
    {
      "items":[
        {
          "dimension":"city",
          "label":"London",
          "links":{
            "versions":{
              "href":"http://hostname/datasets/cantabular-flexible-table-test/editions/latest/versions/1",
              "id":"1"
            },
            "code_list":{
              "href":"http://hostname/code-lists/city",
              "id":"city"
            },
            "code":{
              "href":"http://hostname/code-lists/city/codes/0",
              "id":"0"
            }
          },
          "option":"0"
        },
        {
          "dimension":"city",
          "label":"Liverpool",
          "links":{
            "versions":{
              "href":"http://hostname/datasets/cantabular-flexible-table-test/editions/latest/versions/1",
              "id":"1"
            },
            "code_list":{
              "href":"http://hostname/code-lists/city",
              "id":"city"
            },
            "code":{
              "href":"http://hostname/code-lists/city/codes/1",
              "id":"1"
            }
          },
          "option":"1"
        },
        {
          "dimension":"city",
          "label":"Belfast",
          "links":{
            "versions":{
              "href":"http://hostname/datasets/cantabular-flexible-table-test/editions/latest/versions/1",
              "id":"1"
            },
            "code_list":{
              "href":"http://hostname/code-lists/city",
              "id":"city"
            },
            "code":{
              "href":"http://hostname/code-lists/city/codes/2",
              "id":"2"
            }
          },
          "option":"2"
        }
      ],
      "count":3,
      "offset":0,
      "limit":0,
      "total_count":3
    }
    """

    And the following options document for dataset id "cantabular-flexible-table-test", edition "latest", version "1" and dimension "sex" is available from dp-dataset-api:
    """
    {
      "items":[
        {
          "dimension":"sex",
          "label":"Male",
          "links":{
            "versions":{
              "href":"http://hostname/datasets/cantabular-flexible-table-test/editions/latest/versions/1",
              "id":"1"
            },
            "code_list":{
              "href":"http://hostname/code-lists/sex",
              "id":"sex"
            },
            "code":{
              "href":"http://hostname/code-lists/sex/codes/0",
              "id":"0"
            }
          },
          "option":"0"
        },
        {
          "dimension":"sex",
          "label":"Female",
          "links":{
            "versions":{
              "href":"http://hostname/datasets/cantabular-flexible-table-test/editions/latest/versions/1",
              "id":"1"
            },
            "code_list":{
              "href":"http://hostname/code-lists/sex",
              "id":"sex"
            },
            "code":{
              "href":"http://hostname/code-lists/sex/codes/1",
              "id":"1"
            }
          },
          "option":"1"
        }
      ],
      "count":2,
      "offset":0,
      "limit":0,
      "total_count":2
    }
    """

    And the following options document for dataset id "cantabular-flexible-table-test", edition "latest", version "1" and dimension "siblings_3" is available from dp-dataset-api:
    """
    {
      "items":[
        {
          "dimension":"siblings_3",
          "label":"No siblings",
          "links":{
            "versions":{
              "href":"http://hostname/datasets/cantabular-flexible-table-test/editions/latest/versions/1",
              "id":"1"
            },
            "code_list":{
              "href":"http://hostname/code-lists/siblings_3",
              "id":"siblings_3"
            },
            "code":{
              "href":"http://hostname/code-lists/siblings_3/codes/0",
              "id":"0"
            }
          },
          "option":"0"
        },
        {
          "dimension":"siblings_3",
          "label":"1 or 2 siblings",
          "links":{
            "versions":{
              "href":"http://hostname/datasets/cantabular-flexible-table-test/editions/latest/versions/1",
              "id":"1"
            },
            "code_list":{
              "href":"http://hostname/code-lists/siblings_3",
              "id":"siblings_3"
            },
            "code":{
              "href":"http://hostname/code-lists/siblings_3/codes/1-2",
              "id":"1-2"
            }
          },
          "option":"1-2"
        },
        {
          "dimension":"siblings_3",
          "label":"3 or more siblings",
          "links":{
            "versions":{
              "href":"http://hostname/datasets/cantabular-flexible-table-test/editions/latest/versions/1",
              "id":"1"
            },
            "code_list":{
              "href":"http://hostname/code-lists/siblings_3",
              "id":"siblings_3"
            },
            "code":{
              "href":"http://hostname/code-lists/siblings_3/codes/3+",
              "id":"3+"
            }
          },
          "option":"3+"
        }
      ],
      "count":3,
      "offset":0,
      "limit":0,
      "total_count":3
    }
    """

    And the following options document for dataset id "cantabular-multivariate-table-test", edition "latest", version "1" and dimension "city" is available from dp-dataset-api:
    """
    {
      "items":[
        {
          "dimension":"city",
          "label":"London",
          "links":{
            "versions":{
              "href":"http://hostname/datasets/cantabular-multivariate-table-test/editions/latest/versions/1",
              "id":"1"
            },
            "code_list":{
              "href":"http://hostname/code-lists/city",
              "id":"city"
            },
            "code":{
              "href":"http://hostname/code-lists/city/codes/0",
              "id":"0"
            }
          },
          "option":"0"
        }
      ],
      "count":1,
      "offset":0,
      "limit":0,
      "total_count":1
    }
    """

    And the following options document for dataset id "cantabular-multivariate-table-test", edition "latest", version "1" and dimension "sex" is available from dp-dataset-api:
    """
    {
      "items":[
        {
          "dimension":"sex",
          "label":"Male",
          "links":{
            "versions":{
              "href":"http://hostname/datasets/cantabular-multivariate-table-test/editions/latest/versions/1",
              "id":"1"
            },
            "code_list":{
              "href":"http://hostname/code-lists/sex",
              "id":"sex"
            },
            "code":{
              "href":"http://hostname/code-lists/sex/codes/0",
              "id":"0"
            }
          },
          "option":"0"
        }
      ],
      "count":1,
      "offset":0,
      "limit":0,
      "total_count":1
    }
    """

    And Cantabular returns these geography dimensions for the given request:
    """
    request:
    {
      "query":"query($dataset: String!, $limit: Int!, $offset: Int) {
        dataset(name: $dataset) {
          variables(rule: true, skip: $offset, first: $limit) {
            totalCount
            edges {
              node {
                name
                description
                meta{
                  ONS_Variable{
                    Geography_Hierarchy_Order
                  }
                }
                mapFrom {
                  edges {
                    node {
                      description
                      label
                      name
                    }
                  }
                }
                label
                categories{
                  totalCount
                }
              }
            }
          }
        }
      }",
      "variables": {"base":false,"category":"","dataset":"Example","filters":null,"limit":100,"offset":0,"rule":false,"text":"","variables":null}
    }
    response:
    {
      "data": {
        "dataset": {
          "variables": {
            "edges": [
              {
                "node": {
                  "categories": {"totalCount": 3},
                  "description":"",
                  "label": "City",
                  "mapFrom": [],
                  "name": "city"
                }
              },
              {
                "node": {
                  "categories": {"totalCount": 2},
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
                      ]
                    }
                  ],
                  "name": "country"
                }
              }
            ],
            "totalCount": 2
          }
        }
      }
    }
    """

    And the maximum rows allowed to be returned is 100

   Scenario: Get the dataset observations
    Given Cantabular returns this static dataset for the given request:
    """
    request:
    {
      "query":"query($dataset: String!, $variables: [String!]!, $filters: [Filter!]) {
        dataset(name: $dataset) {
          table(variables: $variables, filters: $filters) {
            rules {
              passed{
                count
              }
              evaluated
              {
                count
              }
              blocked {
                count
              }
            }
            dimensions {
              count
              variable { name label }
              categories { code label }
            }
            values
            error
          }
        }
      }",
      "variables": {"base":false,"category":"","dataset":"Example","filters":null,"limit":20,"offset":0,"rule":false,"text":"","variables":["city", "sex", "siblings_3"]}
    }
    response:
    {
      "data": {
        "dataset": {
          "table": {
            "dimensions": [
              {
                "categories": [
                  {"code": "0", "label": "London"},
                  {"code": "1","label": "Liverpool"},
                  {"code": "2","label": "Belfast" }
                ],
                "count": 3,
                "variable": {"label": "City","name": "city"}
              },
              {
                "categories": [
                  {"code": "0","label": "Male"},
                  {"code": "1","label": "Female"}
                ],
                "count": 2,
                "variable": {"label": "Sex","name": "sex"}
              },
              {
                "categories": [
                  {"code": "0","label": "No siblings"},
                  {"code": "1-2","label": "1 or 2 siblings"},
                  {"code": "3+","label": "3 or more siblings"
                  }
                ],
                "count": 3,
                "variable": {"label": "Number of siblings (3 mappings)", "name": "siblings_3"}
              }
            ],
            "error": null,
            "values": [1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 2]
          }
        }
      }
    }
    """

    And I set the "X-Forwarded-Proto" header to "https"
    And I set the "X-Forwarded-Host" header to "api.example.com"
    And I set the "X-Forwarded-Path-Prefix" header to "v1"
    And URL rewriting is enabled
    When I GET "/datasets/cantabular-flexible-table-test/editions/latest/versions/1/census-observations"

    Then the HTTP status code should be "200"

    And the getDatasetObservationResult result is:
    """
    {
    "observations": [
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "London",
                    "option_id": "0"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Male",
                    "option_id": "0"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "No siblings",
                    "option_id": "0"
                }
            ],
            "observation": 1
        },
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "London",
                    "option_id": "0"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Male",
                    "option_id": "0"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "1 or 2 siblings",
                    "option_id": "1-2"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "London",
                    "option_id": "0"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Male",
                    "option_id": "0"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "3 or more siblings",
                    "option_id": "3+"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "London",
                    "option_id": "0"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Female",
                    "option_id": "1"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "No siblings",
                    "option_id": "0"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "London",
                    "option_id": "0"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Female",
                    "option_id": "1"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "1 or 2 siblings",
                    "option_id": "1-2"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "London",
                    "option_id": "0"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Female",
                    "option_id": "1"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "3 or more siblings",
                    "option_id": "3+"
                }
            ],
            "observation": 1
        },
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "Liverpool",
                    "option_id": "1"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Male",
                    "option_id": "0"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "No siblings",
                    "option_id": "0"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "Liverpool",
                    "option_id": "1"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Male",
                    "option_id": "0"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "1 or 2 siblings",
                    "option_id": "1-2"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "Liverpool",
                    "option_id": "1"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Male",
                    "option_id": "0"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "3 or more siblings",
                    "option_id": "3+"
                }
            ],
            "observation": 1
        },
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "Liverpool",
                    "option_id": "1"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Female",
                    "option_id": "1"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "No siblings",
                    "option_id": "0"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "Liverpool",
                    "option_id": "1"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Female",
                    "option_id": "1"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "1 or 2 siblings",
                    "option_id": "1-2"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "Liverpool",
                    "option_id": "1"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Female",
                    "option_id": "1"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "3 or more siblings",
                    "option_id": "3+"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "Belfast",
                    "option_id": "2"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Male",
                    "option_id": "0"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "No siblings",
                    "option_id": "0"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "Belfast",
                    "option_id": "2"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Male",
                    "option_id": "0"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "1 or 2 siblings",
                    "option_id": "1-2"
                }
            ],
            "observation": 1
        },
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "Belfast",
                    "option_id": "2"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Male",
                    "option_id": "0"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "3 or more siblings",
                    "option_id": "3+"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "Belfast",
                    "option_id": "2"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Female",
                    "option_id": "1"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "No siblings",
                    "option_id": "0"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "Belfast",
                    "option_id": "2"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Female",
                    "option_id": "1"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "1 or 2 siblings",
                    "option_id": "1-2"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "Belfast",
                    "option_id": "2"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Female",
                    "option_id": "1"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "3 or more siblings",
                    "option_id": "3+"
                }
            ],
            "observation": 2
        }
    ],
    "links": {
        "dataset_metadata": {
            "href": "https://api.example.com/v1/datasets/cantabular-flexible-table-test/editions/latest/versions/1/metadata"
        },
        "self": {
            "href": "https://api.example.com/v1/datasets/cantabular-flexible-table-test",
            "id": "cantabular-flexible-table-test"
        },
        "version": {
            "href": "https://api.example.com/v1/datasets/cantabular-flexible-table-test/editions/latest/versions/1",
            "id": "1"
        }
    },
    "total_observations": 18
    }
    """

   Scenario: Get the dataset observations for a specific area-type
    Given Cantabular returns this static dataset for the given request:
    """
    request:
    {
      "query":"query($dataset: String!, $variables: [String!]!, $filters: [Filter!]) {
        dataset(name: $dataset) {
          table(variables: $variables, filters: $filters) {
            rules {
              passed{
                count
              }
              evaluated
              {
                count
              }
              blocked {
                count
              }
            }
            dimensions {
              count
              variable { name label }
              categories { code label }
            }
            values
            error
          }
        }
      }",
      "variables": {"base":false,"category":"","dataset":"Example","filters":null,"limit":20,"offset":0,"rule":false,"text":"","variables":["country", "sex", "siblings_3"]}
    }
    response:
     {
      "data": {
        "dataset": {
          "table": {
            "dimensions": [
             {
            "categories": [
              {
                "code": "E",
                "label": "England"
              },
              {
                "code": "N",
                "label": "Northern Ireland"
              }
            ],
            "count": 2,
            "variable": {
              "label": "Country",
              "name": "country"
            }
          },
              {
                "categories": [
                  {"code": "0","label": "Male"},
                  {"code": "1","label": "Female"}
                ],
                "count": 2,
                "variable": {"label": "Sex","name": "sex"}
              },
              {
                "categories": [
                  {"code": "0","label": "No siblings"},
                  {"code": "1-2","label": "1 or 2 siblings"},
                  {"code": "3+","label": "3 or more siblings"
                  }
                ],
                "count": 3,
                "variable": {"label": "Number of siblings (3 mappings)", "name": "siblings_3"}
              }
            ],
            "error": null,
             "values": [
          1,
          0,
          1,
          0,
          0,
          1,
          0,
          1,
          0,
          0,
          0,
          2
        ]
          }
        }
      }
    }
    """

    When I GET "/datasets/cantabular-flexible-table-test/editions/latest/versions/1/census-observations?area-type=country"

    Then the HTTP status code should be "200"

    And the getDatasetObservationResult result is:
    """
    {
    "observations": [
        {
            "dimensions": [
                {
                    "dimension": "Country",
                    "dimension_id": "country",
                    "option": "England",
                    "option_id": "E"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Male",
                    "option_id": "0"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "No siblings",
                    "option_id": "0"
                }
            ],
            "observation": 1
        },
        {
            "dimensions": [
                {
                    "dimension": "Country",
                    "dimension_id": "country",
                    "option": "England",
                    "option_id": "E"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Male",
                    "option_id": "0"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "1 or 2 siblings",
                    "option_id": "1-2"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "Country",
                    "dimension_id": "country",
                    "option": "England",
                    "option_id": "E"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Male",
                    "option_id": "0"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "3 or more siblings",
                    "option_id": "3+"
                }
            ],
            "observation": 1
        },
        {
            "dimensions": [
                {
                    "dimension": "Country",
                    "dimension_id": "country",
                    "option": "England",
                    "option_id": "E"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Female",
                    "option_id": "1"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "No siblings",
                    "option_id": "0"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "Country",
                    "dimension_id": "country",
                    "option": "England",
                    "option_id": "E"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Female",
                    "option_id": "1"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "1 or 2 siblings",
                    "option_id": "1-2"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "Country",
                    "dimension_id": "country",
                    "option": "England",
                    "option_id": "E"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Female",
                    "option_id": "1"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "3 or more siblings",
                    "option_id": "3+"
                }
            ],
            "observation": 1
        },
        {
            "dimensions": [
                {
                    "dimension": "Country",
                    "dimension_id": "country",
                    "option": "Northern Ireland",
                    "option_id": "N"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Male",
                    "option_id": "0"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "No siblings",
                    "option_id": "0"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "Country",
                    "dimension_id": "country",
                    "option": "Northern Ireland",
                    "option_id": "N"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Male",
                    "option_id": "0"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "1 or 2 siblings",
                    "option_id": "1-2"
                }
            ],
            "observation": 1
        },
        {
            "dimensions": [
                {
                    "dimension": "Country",
                    "dimension_id": "country",
                    "option": "Northern Ireland",
                    "option_id": "N"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Male",
                    "option_id": "0"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "3 or more siblings",
                    "option_id": "3+"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "Country",
                    "dimension_id": "country",
                    "option": "Northern Ireland",
                    "option_id": "N"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Female",
                    "option_id": "1"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "No siblings",
                    "option_id": "0"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "Country",
                    "dimension_id": "country",
                    "option": "Northern Ireland",
                    "option_id": "N"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Female",
                    "option_id": "1"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "1 or 2 siblings",
                    "option_id": "1-2"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "Country",
                    "dimension_id": "country",
                    "option": "Northern Ireland",
                    "option_id": "N"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Female",
                    "option_id": "1"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "3 or more siblings",
                    "option_id": "3+"
                }
            ],
            "observation": 2
        }
    ],
   "links": {
        "dataset_metadata": {
            "href": "https://api.example.com/v1/datasets/cantabular-flexible-table-test/editions/latest/versions/1/metadata"
        },
        "self": {
            "href": "https://api.example.com/v1/datasets/cantabular-flexible-table-test",
            "id": "cantabular-flexible-table-test"
        },
        "version": {
            "href": "https://api.example.com/v1/datasets/cantabular-flexible-table-test/editions/latest/versions/1",
            "id": "1"
        }
    },
    "total_observations": 12
    }
    """

    Scenario: Get the dataset observations for a specific area-type and area
    Given Cantabular returns this static dataset for the given request:
    """
    request:
    {
      "query": "query($dataset: String!, $variables: [String!]!, $filters: [Filter!]) {
        dataset(name: $dataset) {
          table(variables: $variables, filters: $filters) {
            rules {
              passed{
                count
              }
              evaluated
              {
                count
              }
              blocked {
                count
              }
            }
            dimensions {
              count
              variable {
                name
                label
              }
              categories {
                code
                label
              }
            }
            values
            error
          }
        }
      }",
      "variables": {
        "base": false,
        "category": "",
        "dataset": "Example",
        "filters": [
          {
            "codes": ["E"],
            "variable": "country"
          }
        ],
        "limit":20,
        "offset":0,
        "rule": false,
        "text": "",
        "variables": [
          "country",
          "sex",
          "siblings_3"
        ]
      }
    }
    response:
    {
      "data": {
        "dataset": {
          "table": {
            "dimensions": [
              {
                "categories": [
                  {
                    "code": "E",
                    "label": "England"
                  }
                ],
                "count": 2,
                "variable": {
                  "label": "Country",
                  "name": "country"
                }
              },
              {
                "categories": [
                  {
                    "code": "0",
                    "label": "Male"
                  },
                  {
                    "code": "1",
                    "label": "Female"
                  }
                ],
                "count": 2,
                "variable": {
                  "label": "Sex",
                  "name": "sex"
                }
              },
              {
                "categories": [
                  {
                    "code": "0",
                    "label": "No siblings"
                  },
                  {
                    "code": "1-2",
                    "label": "1 or 2 siblings"
                  },
                  {
                    "code": "3+",
                    "label": "3 or more siblings"
                  }
                ],
                "count": 3,
                "variable": {
                  "label": "Number of siblings (3 mappings)",
                  "name": "siblings_3"
                }
              }
            ],
            "error": null,
            "values": [
              1,
              0,
              1,
              0,
              0,
              1
            ]
          }
        }
      }
    }
    """

    And Cantabular returns this area for the given request:
    """
    request:
    {
      "query": "query ($dataset: String!, $text: String!, $category: String!) {
        dataset(name: $dataset) {
          variables(rule: true, names: [$text]) {
            edges {
              node {
                name
                label
                categories(codes: [$category]) {
                  edges {
                    node {
                      code
                      label
                    }
                  }
                }
              }
            }
          }
        }
      }",
      "variables": {
          "base": false,
          "category": "E",
          "dataset": "Example",
          "filters": null,
          "limit": 20,
          "offset": 0,
          "rule": false,
          "text": "country",
          "variables": null
      }
    }
    response:
    {
      "data": {
        "dataset": {
          "variables": {
            "edges": [
              {
                "node": {
                  "categories": {
                    "edges": [
                      {
                        "node": {
                          "code": "E",
                          "label": "England"
                        }
                      }
                    ]
                  },
                  "label": "Country",
                  "name": "country"
                }
              }
            ]
          }
        }
      }
    }
    """

    When I GET "/datasets/cantabular-flexible-table-test/editions/latest/versions/1/census-observations?area-type=country,E"

    Then the HTTP status code should be "200"

    And the getDatasetObservationResult result is:
    """
   {
    "observations": [
        {
            "dimensions": [
                {
                    "dimension": "Country",
                    "dimension_id": "country",
                    "option": "England",
                    "option_id": "E"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Male",
                    "option_id": "0"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "No siblings",
                    "option_id": "0"
                }
            ],
            "observation": 1
        },
        {
            "dimensions": [
                {
                    "dimension": "Country",
                    "dimension_id": "country",
                    "option": "England",
                    "option_id": "E"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Male",
                    "option_id": "0"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "1 or 2 siblings",
                    "option_id": "1-2"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "Country",
                    "dimension_id": "country",
                    "option": "England",
                    "option_id": "E"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Male",
                    "option_id": "0"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "3 or more siblings",
                    "option_id": "3+"
                }
            ],
            "observation": 1
        },
        {
            "dimensions": [
                {
                    "dimension": "Country",
                    "dimension_id": "country",
                    "option": "England",
                    "option_id": "E"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Female",
                    "option_id": "1"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "No siblings",
                    "option_id": "0"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "Country",
                    "dimension_id": "country",
                    "option": "England",
                    "option_id": "E"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Female",
                    "option_id": "1"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "1 or 2 siblings",
                    "option_id": "1-2"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "Country",
                    "dimension_id": "country",
                    "option": "England",
                    "option_id": "E"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Female",
                    "option_id": "1"
                },
                {
                    "dimension": "Number of siblings (3 mappings)",
                    "dimension_id": "siblings_3",
                    "option": "3 or more siblings",
                    "option_id": "3+"
                }
            ],
            "observation": 1
        }
    ],
    "links": {
        "dataset_metadata": {
            "href": "https://api.example.com/v1/datasets/cantabular-flexible-table-test/editions/latest/versions/1/metadata"
        },
        "self": {
            "href": "https://api.example.com/v1/datasets/cantabular-flexible-table-test",
            "id": "cantabular-flexible-table-test"
        },
        "version": {
            "href": "https://api.example.com/v1/datasets/cantabular-flexible-table-test/editions/latest/versions/1",
            "id": "1"
        }
    },
    "total_observations": 6
    }
    """

  Scenario: Get the dataset as observations asking for additional dimensions
    Given Cantabular returns this static dataset for the given request:
    """
    request:
    {
      "query": "query($dataset: String!, $variables: [String!]!, $filters: [Filter!]) {
        dataset(name: $dataset) {
          table(variables: $variables, filters: $filters) {
            rules {
              passed{
                count
              }
              evaluated
              {
                count
              }
              blocked {
                count
              }
            }
            dimensions {
              count
              variable { name label }
              categories { code label }
            }
            values
            error
          }
        }
      }",
      "variables": {
        "base": false,
        "category": "",
        "dataset": "Example",
        "filters": null,
        "limit": 20,
        "offset": 0,
        "rule": false,
        "text": "",
        "variables": [
          "city",
          "sex",
          "age_23_a"
        ]
      }
    }
    response:
    {
      "data": {
        "dataset": {
          "table": {
            "dimensions": [
              {
                "categories": [
                  {"code": "0", "label": "London"},
                  {"code": "1", "label": "Belfast"},
                  {"code": "2", "label": "Liverpool"}
                ],
                "count": 3,
                "variable": {"label": "City","name": "city"}
              },
              {
                "categories": [
                  {"code": "0","label": "Male"},
                   {"code": "1","label": "Female"}
                ],
                "count": 2,
                "variable": {"label": "Sex","name": "sex"}
              },
              {
                "categories": [
                  {"code": "0","label": "23"}
                ],
                "count": 1,
                "variable": {"label": "age_23_a","name": "Age 23"}
              }
            ],
            "error": null,
            "values": [1, 0, 0, 0, 0, 1]
          }
        }
      }
    }
    """

    And Cantabular returns this response for the given request:
    """
    request:
    {
      "query": "query($dataset: String!, $variables: [String!]!) {
         dataset(name: $dataset) {
          variables(names: $variables, rule: false) {
            edges {
              node {
                name
                mapFrom {
                  edges {
                    node {
                      label
                      name
                    }
                  }
                }
                description
                meta {
                  ONS_Variable {
                    Quality_Statement_Text
                    Quality_Summary_URL
                  }
                }
                label
                categories {
                  totalCount
                }
              }
            }
          }
        }
      }",
      "variables": {
          "base": false,
          "category": "",
          "dataset": "Example",
          "filters": null,
          "limit": 20,
          "offset": 0,
          "rule": false,
          "text": "",
          "variables": ["age_23_a"]
      }
    }
    response:
    {
      "data": {
        "dataset": {
          "variables": {
            "edges": [
              {
                "node": {
                  "categories": {
                    "totalCount": 7
                  },
                  "description": "",
                  "label": "age_23_a",
                  "mapFrom": [
                    {
                      "edges": [
                        {
                          "node": {
                            "label": "",
                            "name": "age_23_a"
                          }
                        }
                      ]
                    }
                  ],
                  "name": "age_23_a"
                }
              }
            ]
          }
        }
      }
    }
    """

    And Population Types API returns this GetCategorisations response for the given request:
    """
    {
      "request": {
        "dimension":       "city",
        "populationType":  "Example",  
        "limit":           99999,
        "serviceAuthToken": "testToken"
      },
      "response": {
        "dimensions": []
      }
    }
    """

    And Population Types API returns this GetCategorisations response for the given request:
    """
    {
      "request": {
        "dimension":       "sex",
        "populationType":  "Example",  
        "limit":           99999,
        "serviceAuthToken": "testToken"
      },
      "response": {
        "dimensions": []
      }
    }
    """

    When I GET "/datasets/cantabular-multivariate-table-test/editions/latest/versions/1/census-observations?dimensions=age_23_a"

    Then the HTTP status code should be "200"

    And the getDatasetObservationResult result is:
    """
    {
    "links": {
        "dataset_metadata": {
            "href": "https://api.example.com/v1/datasets/cantabular-multivariate-table-test/editions/latest/versions/1/metadata"
        },
        "self": {
            "href": "https://api.example.com/v1/datasets/cantabular-multivariate-table-test",
            "id": "cantabular-flexible-table-test"
        },
        "version": {
            "href": "https://api.example.com/v1/datasets/cantabular-multivariate-table-test/editions/latest/versions/1",
            "id": "1"
        }
    },
    "observations": [
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "London",
                    "option_id": "0"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Male",
                    "option_id": "0"
                },
                 {
                    "dimension": "age_23_a",
                    "dimension_id": "Age 23",
                    "option": "23",
                    "option_id": "0"
                }
            ],
            "observation": 1
        },
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "London",
                    "option_id": "0"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Female",
                    "option_id": "1"
                },
                 {
                    "dimension": "age_23_a",
                    "dimension_id": "Age 23",
                    "option": "23",
                    "option_id": "0"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "Belfast",
                    "option_id": "1"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Male",
                    "option_id": "0"
                },
                {
                    "dimension": "age_23_a",
                    "dimension_id": "Age 23",
                    "option": "23",
                    "option_id": "0"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "Belfast",
                    "option_id": "1"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Female",
                    "option_id": "1"
                },
                 {
                    "dimension": "age_23_a",
                    "dimension_id": "Age 23",
                    "option": "23",
                    "option_id": "0"
                }
            ],
            "observation": 0
        },
        {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "Liverpool",
                    "option_id": "2"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Male",
                    "option_id": "0"
                },
                 {
                    "dimension": "age_23_a",
                    "dimension_id": "Age 23",
                    "option": "23",
                    "option_id": "0"
                }
            ],
            "observation": 0
        },
         {
            "dimensions": [
                {
                    "dimension": "City",
                    "dimension_id": "city",
                    "option": "Liverpool",
                    "option_id": "2"
                },
                {
                    "dimension": "Sex",
                    "dimension_id": "sex",
                    "option": "Female",
                    "option_id": "1"
                },
                 {
                    "dimension": "age_23_a",
                    "dimension_id": "Age 23",
                    "option": "23",
                    "option_id": "0"
                }
            ],
            "observation": 1
        }
        
    ],
    "total_observations": 6
   }
    """

  Scenario: Get the dataset observations asking for additional dimensions from incorrect dataset type

    When I GET "/datasets/cantabular-flexible-table-test/editions/latest/versions/1/census-observations?dimensions=age_23_a"

    Then I should receive the following JSON response:
    """
    {
      "errors":["failed to get dataset params: invalid dataset type for custom dimensions"]
    }
    """

    And the HTTP status code should be "400"

  Scenario: Too many observations returned
    Given Cantabular returns this static dataset for the given request:
    """
    request:
    {
      "query":"query($dataset: String!, $variables: [String!]!, $filters: [Filter!]) {
        dataset(name: $dataset) {
          table(variables: $variables, filters: $filters) {
            rules {
              passed{
                count
              }
              evaluated
              {
                count
              }
              blocked {
                count
              }
            }
            dimensions {
              count
              variable { name label }
              categories { code label }
            }
            values
            error
          }
        }
      }",
      "variables": {"base":false,"category":"","dataset":"Example","filters":null,"limit":20,"offset":0,"rule":false,"text":"","variables":["city", "sex", "siblings_3"]}
    }
    response:
    {
      "data": {
        "dataset": {
          "table": {
            "dimensions": [
              {
                "categories": [
                  {"code": "0", "label": "London"},
                  {"code": "1","label": "Liverpool"},
                  {"code": "2","label": "Belfast" }
                ],
                "count": 3,
                "variable": {"label": "City","name": "city"}
              },
              {
                "categories": [
                  {"code": "0","label": "Male"},
                  {"code": "1","label": "Female"}
                ],
                "count": 2,
                "variable": {"label": "Sex","name": "sex"}
              },
              {
                "categories": [
                  {"code": "0","label": "No siblings"},
                  {"code": "1-2","label": "1 or 2 siblings"},
                  {"code": "3+","label": "3 or more siblings"
                  }
                ],
                "count": 3,
                "variable": {"label": "Number of siblings (3 mappings)", "name": "siblings_3"}
              }
            ],
            "error": null,
            "values": [1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 2]
          }
        }
      }
    }
    """
    And the maximum rows allowed to be returned is 3
   
    When I GET "/datasets/cantabular-flexible-table-test/editions/latest/versions/1/census-observations"
    Then the HTTP status code should be "403"
    Then I should receive the following JSON response:
    """
    {
    "errors": [
        "Too many rows returned, please refine your query by requesting specific areas or reducing the number of categories returned.  For further information please visit https://developer.ons.gov.uk/createyourowndataset/"
    ]
    }
    """