Feature: Multivariate Filters Private Endpoints Enabled

  Background:
    Given I am identified as "user@ons.gov.uk"
    And I am authorised
    Given private endpoints are enabled

    And the following version document with dataset id "cantabular-example-1", edition "2021" and version "1" is available from dp-dataset-api:
      """
      {
        "alerts": [],
        "collection_id": "dfb-38b11d6c4b69493a41028d10de503aabed3728828e17e64914832d91e1f493c6",
        "is_based_on":{"@type": "cantabular_multivariate_table"},
        "dimensions": [
           {
                        "href": "http://localhost:22400/code-lists/ladcd",
                        "id": "ladcd",
                        "is_hierarchy": false,
                        "name": "ladcd",
                        "is_cantabular_geography": true,
                        "is_cantabular_default_geography": true
                    },
                    {
                        "href": "http://localhost:22400/code-lists/hh_tenure_9a",
                        "id": "hh_tenure_9a",
                        "is_hierarchy": false,
                        "name": "hh_tenure_9a",
                        "is_cantabular_geography": false,
                         "is_cantabular_default_geography": false
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
  Scenario: Creating a new multivariate filter journey when authorized

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
        "name": "ethnic_group",
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
      "type":            "cantabular_multivariate_table"
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
      "etag":        "22d6ed269ce9691e27e7abaa424159017cd3609b",
      "instance_id": "c733977d-a2ca-4596-9cb1-08a6e724858b",
      "dimensions": [
            {

        "name": "ladcd",
        "id": "ladcd",
        "label": "",
        "is_area_type": true,
        "options": []
      },
      {

        "name": "ethnic_group",
        "is_area_type": false,
        "id": "",
        "label": "",
        "options": []

      },
    {

         "name": "hh_deprivation",
        "is_area_type": false,
        "id": "",
        "label": "",
        "options": []
      }
      ],
      "dataset": {
        "id":      "cantabular-example-1",
        "edition": "2021",
        "version": {
          "$numberInt":"1"
        }
      },
      "published":       true,
      "population_type": "dummy_data_households",
      "type":            "cantabular_multivariate_table",
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
