Feature: Put Filter Outputs Private Endpoints Not Enabled

  Background:
    Given private endpoints are enabled

    And I have these filter outputs:
    """
    [
      {
        "id": "94310d8d-72d6-492a-bc30-27584627edb1",
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

    Then the HTTP status code should be "200"

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
      "errors": ["failed to update filter output: failed to find filter output: filter output not found for 94310d8d-72d6-492a-bc30-27584627new1"]
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
        "failed to update filter output: failed to upsert filter"
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
