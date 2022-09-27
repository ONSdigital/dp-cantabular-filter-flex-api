Feature: Get Dataset JSON

  Background:
    Given private endpoints are not enabled

    And the Cantabular service is a mocked extended Cantabular server

    And the following recipe is used to create a dataset based on the given cantabular dataset:
    """
    {
      "recipe": {
        "_id": "6cf112cb-87bd-41f5-9a70-e6abd67de4f2",
        "alias": "cantabular flexible table component test",
        "cantabular_blob": "Example",
        "format": "cantabular_flexible_table",
        "output_instances": [
          {
            "dataset_id": "cantabular-flexible-table-component-test",
            "editions": ["latest"],
            "title": "cantabular flexible table component test",
            "code_lists": [
              {
                "id": "city",
                "href": "http://hostname/code-lists/city",
                "name": "City",
                "is_hierarchy": false,
                "is_cantabular_geography": true,
                "is_cantabular_default_geography": true
              },
              {
                "id": "sex",
                "href": "http://hostname/code-lists/sex",
                "name": "Sex",
                "is_hierarchy": false,
                "is_cantabular_geography": false,
                "is_cantabular_default_geography": false
              },
              {
                "id": "siblings_3",
                "href": "http://hostname/code-lists/siblings_3",
                "name": "Number of siblings (3 mappings)",
                "is_hierarchy": false,
                "is_cantabular_geography": false,
                "is_cantabular_default_geography": false
              }
            ]
          }
        ]
      },
      "cantabular_dataset": {
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
                mapFrom {
                  edges {
                    node {
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
      "variables": {"base":false,"category":"","dataset":"Example","filters":null,"limit":20,"offset":0,"rule":false,"text":"","variables":null}
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

  Scenario: Get the dataset as JSON without asking for specific dimensions
    Given Cantabular returns this static dataset for the given request:
    """
    request:
    {
      "query":"query($dataset: String!, $variables: [String!]!, $filters: [Filter!]) {
        dataset(name: $dataset) {
          table(variables: $variables, filters: $filters) {
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

    When I GET "/datasets/cantabular-flexible-table-component-test/editions/latest/versions/1/json"

    Then the HTTP status code should be "200"

    And I should receive the following JSON response:
    """
    {
      "dimensions":[
        {
          "dimension_name":"city",
          "options":[
            {
              "href":"http://hostname/code-lists/city/codes/0",
              "id":"0"
            },
            {
              "href":"http://hostname/code-lists/city/codes/1",
              "id":"1"
            },
            {
              "href":"http://hostname/code-lists/city/codes/2",
              "id":"2"
            }
          ]
        },
        {
          "dimension_name":"sex",
          "options":[
            {
              "href":"http://hostname/code-lists/sex/codes/0",
              "id":"0"
            },
            {
              "href":"http://hostname/code-lists/sex/codes/1",
              "id":"1"
            }
          ]
        },
        {
          "dimension_name":"siblings_3",
          "options":[
            {
              "href":"http://hostname/code-lists/siblings_3/codes/0",
              "id":"0"
            },
            {
              "href":"http://hostname/code-lists/siblings_3/codes/1-2",
              "id":"1-2"
            },
            {
              "href":"http://hostname/code-lists/siblings_3/codes/3+",
              "id":"3+"
            }
          ]
        }
      ],
      "links":{
        "dataset_metadata": {
          "href":"http://hostname/datasets/cantabular-flexible-table-component-test/editions/latest/versions/1/metadata"
        },
        "self":{
          "href":"http://hostname/datasets/cantabular-flexible-table-component-test",
          "id":"cantabular-flexible-table-component-test"
        },
        "version":{
          "href":"http://hostname/datasets/cantabular-flexible-table-component-test/editions/latest/versions/1",
          "id":"1"
        }
      },
      "observations":[1,0, 0, 0, 0, 1, 0,  0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 2],
      "total_observations":18
    }
    """

  Scenario: Get the dataset as JSON asking for specific geography and non-geography dimensions
    Given Cantabular returns this static dataset for the given request:
    """
    request:
    {
      "query":"query($dataset: String!, $variables: [String!]!, $filters: [Filter!]) {
        dataset(name: $dataset) {
          table(variables: $variables, filters: $filters) {
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
      "variables": {"base":false,"category":"","dataset":"Example","filters":[{"codes": ["0", "1"], "variable": "city"},  {"codes": ["1"], "variable": "sex"}],"limit":20,"offset":0,"rule":false,"text":"","variables":["city", "sex", "siblings_3"]}
    }
    response:
    {
      "data": {
        "dataset": {
          "table": {
            "dimensions": [
              {
                "categories": [
                  {"code": "1","label": "Liverpool"},
                  {"code": "2","label": "Belfast" }
                ],
                "count": 2,
                "variable": {"label": "City","name": "city"}
              },
              {
                "categories": [
                  {"code": "1","label": "Female"}
                ],
                "count": 1,
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
            "values": [0, 0, 0, 0, 0, 2]
          }
        }
      }
    }
    """

    When I GET "/datasets/cantabular-flexible-table-component-test/editions/latest/versions/1/json?geography=city,London,Liverpool&dimension=sex&options=Female"

    Then the HTTP status code should be "200"

    And I should receive the following JSON response:
    """
    {
      "dimensions":[
        {
          "dimension_name":"city",
          "options":[
            {
              "href":"http://hostname/code-lists/city/codes/1",
              "id":"1"
            },
            {
              "href":"http://hostname/code-lists/city/codes/2",
              "id":"2"
            }
          ]
        },
        {
          "dimension_name":"sex",
          "options":[
            {
              "href":"http://hostname/code-lists/sex/codes/1",
              "id":"1"
            }
          ]
        },
        {
          "dimension_name":"siblings_3",
          "options":[
            {
              "href":"http://hostname/code-lists/siblings_3/codes/0",
              "id":"0"
            },
            {
              "href":"http://hostname/code-lists/siblings_3/codes/1-2",
              "id":"1-2"
            },
            {
              "href":"http://hostname/code-lists/siblings_3/codes/3+",
              "id":"3+"
            }
          ]
        }
      ],
      "links":{
        "dataset_metadata": {
          "href":"http://hostname/datasets/cantabular-flexible-table-component-test/editions/latest/versions/1/metadata"
        },
        "self":{
          "href":"http://hostname/datasets/cantabular-flexible-table-component-test",
          "id":"cantabular-flexible-table-component-test"
        },
        "version":{
          "href":"http://hostname/datasets/cantabular-flexible-table-component-test/editions/latest/versions/1",
          "id":"1"
        }
      },
      "observations":[0, 0, 0, 0, 0, 2],
      "total_observations":6
    }
    """
