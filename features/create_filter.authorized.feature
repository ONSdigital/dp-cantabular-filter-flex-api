Feature: Filters Private Endpoints Enabled

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

  Scenario: Creating a new filter journey when authorized
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
      "population_type": "Example",
      "dimensions": [
        {
          "name": "siblings",
          "options": [
            "0-3",
            "4-7",
            "7+"
          ],
          "is_area_type": false
        },{
          "name": "geography",
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
      "instance_id":      "c733977d-a2ca-4596-9cb1-08a6e724858b",
      "dataset": {
        "id":      "cantabular-example-1",
        "edition": "2021",
        "version": 1
      },
      "population_type": "Example",
      "published": true,
      "type": "flexible"
    }
    """

    And the HTTP status code should be "201"

    And the document in the database for id "94310d8d-72d6-492a-bc30-27584627edb1" should match:
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
          "options": [
            "0-3",
            "4-7",
            "7+"
          ],
          "dimension_url": "http://dimension.url/siblings",
          "is_area_type":  false
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
          "dimension_url": "http://dimension.url/city",
          "is_area_type":  true
        }
      ],
      "dataset": {
        "id":      "c7b634c9-b4e9-4e7a-a0b8-d255d38db200",
        "edition": "2021",
        "version": 1
      },
      "published":       true,
      "population_type": "Example"
    }
    """

  Scenario: Creating a new filter journey when not authorized
    Given I am not identified

    And I am not authorised

    When I POST "/filters"
    """
    {"foo":"bar"}
    """

    Then the HTTP status code should be "401"
