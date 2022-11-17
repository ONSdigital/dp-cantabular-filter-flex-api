Feature: Filters Private Endpoints Enabled

  Background:
    Given private endpoints are enabled with permissions checking

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
            "name": "city",
            "is_area_type": true
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
            "name": "siblings_3",
            "is_area_type": false
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

    And Cantabular returns these dimensions for the dataset "Example" and search term "city":
    """
    {
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
                          "label": "LSOA",
                          "name": "lsoa"
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

  Scenario: Creating a new filter journey when authorized
    Given I use an X Florence user token "user token"

    And I am identified as "user@ons.gov.uk"

    And zebedee recognises the user token as valid

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
          "is_area_type": true,
          "filter_by_parent": "country"
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
        "version": 1,
        "lowest_geography": "lowest-geography",
        "release_date": "2021-11-19T00:00:00.000Z",
        "title": "cantabular-example-1"
      },
      "population_type": "Example",
      "published":       true,
      "type":            "flexible"
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
      "etag":        "ce22e9f70ec663fe65341ab4b439d7e5e3fa362f",
      "instance_id": "c733977d-a2ca-4596-9cb1-08a6e724858b",
      "dimensions": [
        {
          "name":  "city",
          "id":    "city",
          "label": "City",
          "options": [
            "Cardiff",
            "London",
            "Swansea"
          ],
          "is_area_type":  true,
          "filter_by_parent": "country"
        },
        {
          "name":  "siblings_3",
          "id":    "siblings_3",
          "label": "Number of siblings (3 mappings)",
          "options": [
            "0-3",
            "4-7",
            "7+"
          ],
          "is_area_type":  false
        }
      ],
      "dataset": {
        "id":      "cantabular-example-1",
        "edition": "2021",
        "version": {
          "$numberInt":"1"
        },
        "lowest_geography": "lowest-geography",
        "release_date": "2021-11-19T00:00:00.000Z",
        "title": "cantabular-example-1"
      },
      "published":       true,
      "type":            "flexible",
      "published":       true,
      "population_type": "Example",
      "type":            "flexible",
      "unique_timestamp":{
        "$timestamp":{
          "i": 1,
          "t": 1.643200024e+09
        }
      },
      "last_updated":{
        "$date":{
          "$numberLong": "1643200024783"
        }
      }
    }
    """

  Scenario: Creating a new filter journey when not authorized
    Given I use an X Florence user token "user token"

    And I am identified as "user@ons.gov.uk"

    But zebedee does not recognise the user token

    When I POST "/filters"
    """
    {"foo":"bar"}
    """

    Then the HTTP status code should be "401"

  Scenario: Creating a new filter journey when not authenticated
    Given I use an X Florence user token "bad user token"

    But I am not identified

    When I POST "/filters"
    """
    {"foo":"bar"}
    """

    Then the HTTP status code should be "401"


  Scenario: Creating a new filter journey when not authenticated (no token set)
    When I POST "/filters"
    """
    {"foo":"bar"}
    """

    Then the HTTP status code should be "401"
