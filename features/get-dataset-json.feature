Feature: Get Dataset JSON

  Background:
    Given private endpoints are not enabled

    And The Cantabular service is a real extended Cantabular server listening on the configured urls:

    And the following recipe is used to create a dataset via the cantabular server:
    """
    {
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
              "href": "http://localhost:22400/code-lists/city",
              "name": "City",
              "is_hierarchy": false,
              "is_cantabular_geography": true,
              "is_cantabular_default_geography": true
            },
            {
              "id": "sex",
              "href": "http://localhost:22400/code-lists/sex",
              "name": "Sex",
              "is_hierarchy": false,
              "is_cantabular_geography": false,
              "is_cantabular_default_geography": false
            },
            {
              "id": "siblings_3",
              "href": "http://localhost:22400/code-lists/siblings_3",
              "name": "Number of siblings (3 mappings)",
              "is_hierarchy": false,
              "is_cantabular_geography": false,
              "is_cantabular_default_geography": false
            }
          ]
        }
      ]
    }
    """

  Scenario: Get the dataset version as JSON
    When I GET "/datasets/cantabular-flexible-table-component-test/editions/latest/versions/1/json"

    Then the HTTP status code should be "200"

    And the getDatasetJSON result should be:
    """
    {
      "dimensions":[
        {
          "dimension_name":"city",
          "options":[
            {
              "href":"http://localhost:22400/code-lists/city/codes/0",
              "id":"0"
            },
            {
              "href":"http://localhost:22400/code-lists/city/codes/1",
              "id":"1"
            },
            {
              "href":"http://localhost:22400/code-lists/city/codes/2",
              "id":"2"
            }
          ]
        },
        {
          "dimension_name":"sex",
          "options":[
            {
              "href":"http://localhost:22400/code-lists/sex/codes/0",
              "id":"0"
            },
            {
              "href":"http://localhost:22400/code-lists/sex/codes/1",
              "id":"1"
            }
          ]
        },
        {
          "dimension_name":"siblings_3",
          "options":[
            {
              "href":"http://localhost:22400/code-lists/siblings_3/codes/0",
              "id":"0"
            },
            {
              "href":"http://localhost:22400/code-lists/siblings_3/codes/1-2",
              "id":"1-2"
            },
            {
              "href":"http://localhost:22400/code-lists/siblings_3/codes/3+",
              "id":"3+"
            }
          ]
        }
      ],
      "links":{
        "dataset_metadata": {},
        "self":{
          "href":"http://localhost:9999/datasets/cantabular-flexible-table-component-test",
          "id":"cantabular-flexible-table-component-test"
        },
        "version":{
          "href":"http://localhost:9999/datasets/cantabular-flexible-table-component-test/editions/latest/versions/1",
          "id":"1"
        }
      },
      "observations":[1,0, 0, 0, 0, 1, 0,  0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 2],
      "total_observations":18
    }
    """

#  metatdata should be as follows (update dp-api-clients to unmarshal the Metadata struct correctly)
#  "dataset_metadata": {
#          "href":"http://localhost:9999/datasets/cantabular-flexible-table-component-test/editions/latest/versions/1/metadata"
#  },
