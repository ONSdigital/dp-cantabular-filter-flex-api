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
        "type":"cantabular_multivariate_table",
        "is_based_on":{
          "@type":"cantabular_multivariate_table",
          "@id":"Example"
        }
      }
      """

    And the following version document with dataset id "cantabular-example-1", edition "2021" and version "1" is available from dp-dataset-api:
      """
      {
        "alerts": [],
        "collection_id": "dfb-38b11d6c4b69493a41028d10de503aabed3728828e17e64914832d91e1f493c6",
        "is_based_on":{"@type": "cantabular_multivariate_table"},
        "dimensions": [
          {
            "label": "ltla",
            "links": {
              "code_list": {},
              "options": {},
              "version": {}
            },
            "href": "http://api.localhost:23200/v1/code-lists/ltla",
            "id": "ltla",
            "name": "ltla",
            "is_area_type": true
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
        "is_based_on":{"@type": "cantabular_multivariate_table"},
        "dimensions": [
          {
            "label": "ltla",
            "links": {
              "code_list": {},
              "options": {},
              "version": {}
            },
            "href": "http://api.localhost:23200/v1/code-lists/ltla",
            "id": "ltla",
            "name": "ltla"
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

    And the following version document with dataset id "cantabular-example-non-multivariate", edition "2021" and version "1" is available from dp-dataset-api:
      """
      {
        "alerts": [],
        "collection_id": "dfb-38b11d6c4b69493a41028d10de503aabed3728828e17e64914832d91e1f493c6",
        "is_based_on":{"@type": "cantabular_flexible_table"},
        "dimensions": [
          {
            "label": "ltla",
            "links": {
              "code_list": {},
              "options": {},
              "version": {}
            },
            "href": "http://api.localhost:23200/v1/code-lists/ltla",
            "id": "ltla",
            "name": "ltla"
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
        "state": "published",
        "usage_notes": [],
        "version": 1
      }
      """
    And Cantabular returns dimensions for the dataset "dummy_data_households" for the following search terms:
      """
      {
        "responses": {
          "ltla": {
            "dataset": {
              "variables": {
                "edges": [
                  {
                    "node": {
                      "name": "ltla",
                      "label": "ltla",
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
    Given Population Types API returns this GetDefaultDatasetMetadata response for the given request:
    """
    {
        "population_type": "Example",
        "default_dataset_id": "cantabular-example-1",
        "edition": "2021",
        "version": 1
    }
    """

    When I POST "/filters/custom"
    """
    {
        "population_type": "Example" 
    }
    """

    Then I should receive the following JSON response:
    """
    {
      "filter_id": "94310d8d-72d6-492a-bc30-27584627edb1",
      "links": {
        "version": {
          "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1",
          "id": "1"
        },
        "self": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1"
        },
        "dimensions": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
        }
      },
      "instance_id": "c733977d-a2ca-4596-9cb1-08a6e724858b",
      "dataset": {
        "id": "cantabular-example-1",
        "edition": "2021",
        "version": 1,
        "lowest_geography": "",
        "release_date": "",
        "title": "custom"
      },
      "population_type": "Example",
      "published": true,
      "type": "multivariate",
      "custom": false
    }
    """

  Scenario: Creating a new filter unauthenticated on unpublished version
    Given Population Types API returns this GetDefaultDatasetMetadata response for the given request:
    """
    {
        "population_type": "Example",
        "default_dataset_id": "cantabular-example-unpublished",
        "edition": "2021",
        "version": 1
    }
    """

    When I POST "/filters/custom"
    """
    {
        "population_type": "Example" 
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


  Scenario: Creating a custom filter on non multivariate dataset
    Given Population Types API returns this GetDefaultDatasetMetadata response for the given request:
    """
    {
        "population_type": "Example",
        "default_dataset_id": "cantabular-example-non-multivariate",
        "edition": "2021",
        "version": 1
    }
    """
    
    When I POST "/filters/custom"
    """
    {
        "population_type": "Example" 
    }
    """

    Then I should receive the following JSON response:
    """
    {
      "errors": [
        "default dataset is not of type multivariate table"
      ]
    }
    """

    And the HTTP status code should be "400"

  Scenario: Creating a new filter bad request body

    When I POST "/filters/custom"
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
