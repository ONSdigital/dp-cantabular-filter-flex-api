Feature: Put Filter Outputs Private Endpoints Not Enabled

  Background:
    Given private endpoints are enabled

    And I have these filter outputs:
    """
    [
      {
        "id": "94310d8d-72d6-492a-bc30-27584627edb1",
        "state": "published",
        "dataset": {
          "edition":"2021",
          "id":"cantabular-flexible-example",
          "version": 1
        },
        "dimensions": [
          {
            "name": "silbings",
            "id": "siblings_3",
            "label": "Number Of Siblings (3 Mappings)",
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
        "etag": "testEtag",
        "events": null,
        "filter_id": "74310d8d-72d6-492a-bc30-27584627edb3",
        "id": "94310d8d-72d6-492a-bc30-27584627edb1",
        "instance_id": "7de30d8d-72d62-412a-bc30-27584627ede4",
        "links":{
          "filter_blueprint":{
            "href":":27100/filters/74310d8d-72d6-492a-bc30-27584627edb3"
          },
          "self":{
            "href":":27100/filter-outputs/94310d8d-72d6-492a-bc30-27584627edb1"
          },
          "version":{
            "href":":27100/datasets/cantabular-flexible-example/editions/2021/versions/1"
          }
        },
        "population_type": "Example",
        "published": true,
        "type": "flexible",
        "downloads":
        {
          "csv":
          {
            "href":"http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.csv",
            "private":"http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.csv",
            "public":"https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.csv",
            "size":"277",
            "skipped":true
          },
          "csvw":
          {
            "href" : "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.csv-metadata.json",
            "private" : "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.csvw",
            "public" : "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.csvw",
            "size" : "607",
            "skipped": true
          },
          "txt":
          {
            "href" : "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.txt",
            "private" : "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.txt",
            "public" : "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.txt",
            "size" : "530",
            "skipped": true
          }
        }
      },
      {
        "id": "94310d8d-72d6-492a-bc30-27584627edb2",
        "state": "completed",
        "downloads":
        {
          "xls":
          {
            "href": "string",
            "size": "string",
            "public": "string",
            "private": "string",
            "skipped": true
          },
          "csv":
          {
            "href": "string",
            "size": "string",
            "public": "string",
            "private": "string",
            "skipped": true
          },
          "csvw":
          {
            "href": "string",
            "size": "string",
            "public": "string",
            "private": "string",
            "skipped": true
          },
          "txt":
          {
            "href": "string",
            "size": "string",
            "public": "string",
            "private": "string",
            "skipped": true
          }
        }
      }
    ]
    """

  Scenario: Update filter output successfully with existing ID in DB
    Given I am identified as "user@ons.gov.uk"

    And I am authorised

    When I PUT "/filter-outputs/94310d8d-72d6-492a-bc30-27584627edb1"
    """
    {
        "state": "published",
        "downloads":
        {
            "xls":
            {
              "href":"http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.xls",
              "private":"http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.xls",
              "size":"6944",
              "public":"https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.xls",
              "skipped":true
            }
        }
    }
    """

    Then the HTTP status code should be "200"

    And a document in collection "filterOutputs" with key "id" value "94310d8d-72d6-492a-bc30-27584627edb1" should match:
    """
    {
        "_id": "94310d8d-72d6-492a-bc30-27584627edb1",
        "dataset": {
          "edition": "2021",
          "id": "cantabular-flexible-example",
          "version": {
            "$numberInt": "1"
          }
        },
        "dimensions": [
          {
            "name": "silbings",
            "id": "siblings_3",
            "label": "Number Of Siblings (3 Mappings)",
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
        "etag": "",
        "state": "published",
        "events": null,
        "filter_id": "74310d8d-72d6-492a-bc30-27584627edb3",
        "id": "94310d8d-72d6-492a-bc30-27584627edb1",
        "instance_id": "7de30d8d-72d62-412a-bc30-27584627ede4",
        "last_updated": {
          "$date":{
            "$numberLong":"-62135596800000"
          }
        },
        "links":{
          "filterblueprint":{
            "href":":27100/filters/74310d8d-72d6-492a-bc30-27584627edb3"
          },
          "self":{
            "href":":27100/filter-outputs/94310d8d-72d6-492a-bc30-27584627edb1"
          },
          "version":{
            "href":":27100/datasets/cantabular-flexible-example/editions/2021/versions/1"
          }
        },
        "population_type": "Example",
        "published": true,
        "type": "flexible",
        "unique_timestamp": {
          "$timestamp":{
            "i":0,
            "t":0
          }
        },
        "downloads":
        {
            "xls":
            {
              "href":"http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.xls",
              "private":"http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.xls",
              "size":"6944",
              "public":"https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.xls",
              "skipped":true
            },
            "csv":
            {
              "href":"http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.csv",
              "private":"http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.csv",
              "public":"https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.csv",
              "size":"277",
              "skipped":true
            },
            "csvw":
            {
              "href" : "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.csv-metadata.json",
              "private" : "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.csvw",
              "public" : "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.csvw",
              "size" : "607",
              "skipped": true
            },
            "txt":
            {
              "href" : "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.txt",
              "private" : "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.txt",
              "public" : "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.txt",
              "size" : "530",
              "skipped": true
            }
        }
    }
    """

  Scenario: Update filter output successfully with new ID in DB
    Given I am identified as "user@ons.gov.uk"

    And I am authorised

    When I PUT "/filter-outputs/94310d8d-72d6-492a-bc30-27584627new1"
    """
    {
        "state": "published",
        "downloads":
        {
            "xls":
            {
              "href":"http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.xls",
              "private":"http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.xls",
              "size":"6944",
              "public":"https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.xls",
              "skipped":true
            },
            "csv":
            {
              "href":"http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.csv",
              "private":"http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.csv",
              "public":"https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.csv",
              "size":"277",
              "skipped":true
            },
            "csvw":
            {
              "href" : "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.csv-metadata.json",
              "private" : "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.csvw",
              "public" : "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.csvw",
              "size" : "607",
              "skipped": true
            },
            "txt":
            {
              "href" : "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.txt",
              "private" : "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.txt",
              "public" : "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.txt",
              "size" : "530",
              "skipped": true
            }
        }
    }
    """

    Then I should receive the following JSON response:
    """
    {
      "errors": ["failed to update filter output"]
    }
    """
    And the HTTP status code should be "404"


  Scenario: Update a filter output with broken mongo db
    Given I am identified as "user@ons.gov.uk"

    And I am authorised

    And Mongo datastore fails for update filter output

    When I PUT "/filter-outputs/94310d8d-72d6-492a-bc30-27584627edb1"
    """
    {
        "state": "published",
        "downloads":
        {
            "xls":
            {
              "href":"http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.xls",
              "private":"http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.xls",
              "size":"6944",
              "public":"https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.xls",
              "skipped":true
            },
            "csv":
            {
              "href":"http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.csv",
              "private":"http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.csv",
              "public":"https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.csv",
              "size":"277",
              "skipped":true
            },
            "csvw":
            {
              "href" : "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.csv-metadata.json",
              "private" : "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.csvw",
              "public" : "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.csvw",
              "size" : "607",
              "skipped": true
            },
            "txt":
            {
              "href" : "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.txt",
              "private" : "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.txt",
              "public" : "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.txt",
              "size" : "530",
              "skipped": true
            }
        }
    }
    """

    Then I should receive the following JSON response:
    """
    {
      "errors": [
        "failed to update filter output"
      ]
    }
    """

    And the HTTP status code should be "500"

    Scenario: Creating a new filter outputs bad request body
    Given I am identified as "user@ons.gov.uk"

    And I am authorised
    When I PUT "/filter-outputs/94310d8d-72d6-492a-bc30-27584627edb1"
    """
    {
      "ins
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

  Scenario: Updating filter output when the state is completed
    Given I am identified as "user@ons.gov.uk"

    And I am authorised

    When I PUT "/filter-outputs/94310d8d-72d6-492a-bc30-27584627edb2"
    """
    {
        "state": "published",
        "downloads":
        {
            "xls":
            {
              "href" : "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.xls",
              "private" : "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.xls",
              "size" : "6944",
              "public" : "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.xls",
              "skipped": true
            },
            "csv":
            {
              "href" : "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.csv",
              "private" : "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.csv",
              "public" : "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.csv",
              "size" : "277",
              "skipped": true
            },
            "csvw":
            {
              "href" : "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.csv-metadata.json",
              "private" : "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.csvw",
              "public" : "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.csvw",
              "size" : "607",
              "skipped": true
            },
            "txt":
            {
              "href" : "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.txt",
              "private" : "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.txt",
              "public" : "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.csvw",
              "size" : "530",
              "skipped": true
            }
        }
    }
    """

    Then the HTTP status code should be "403"

  Scenario: Updating a filter output with one empty field in request body
    Given I am identified as "user@ons.gov.uk"

    And I am authorised

    When I PUT "/filter-outputs/94310d8d-72d6-492a-bc30-27584627edb1"
    """
    {
        "state": "string",
        "downloads":
        {
            "xls":
            {
              "href": "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.xls",
              "private": "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.xls",
              "size": "6944",
              "public": "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.xls",
              "skipped": true
            },
            "csv":
            {
              "href": "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.csv",
              "private": "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.csv",
              "public": "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.csv",
              "size": "277",
              "skipped": true
            },
            "csvw":
            {
              "href": "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.csv-metadata.json",
              "private": "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.csvw",
              "public": "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.csvw",
              "size": "607",
              "skipped": true
            },
            "txt":
            {
              "href": "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.txt",
              "private": "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.txt",
              "public": "",
              "size": "530",
              "skipped": true
            }
        }
    }
    """

    Then I should receive the following JSON response:
    """
    {
      "errors": [
        "invalid request body: invalid request: 'txt' field not fully populated: \"public\" is empty in input"
      ]
    }
    """

    And the HTTP status code should be "400"
