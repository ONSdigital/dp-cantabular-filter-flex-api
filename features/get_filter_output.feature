Feature: Get Filter Private Endpoints Not Enabled

  Background:
    Given private endpoints are not enabled

    And I have these filter outputs:
    """
    [{
       "id": "94310d8d-72d6-492a-bc30-27584627edb1",
       "filter_id": "94310d8d-72d6-492a-bc30-27584627fil1",
       "instance_id": "94310d8d-72d6-492a-bc30-27584627inst1",
       "state": "published",
       "published": true,
       "downloads": {
         "xls": {
           "href": "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.xls",
           "private": "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.xls",
           "size": "6944",
           "public": "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.xls",
           "skipped": true
         },
         "csv": {
           "href": "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.csv",
           "private": "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.csv",
           "public": "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.csv",
           "size": "277",
           "skipped": true
         },
         "csvw": {
           "href": "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.csv-metadata.json",
           "private": "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.csvw",
           "public": "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.csvw",
           "size": "607",
           "skipped": true
         },
         "txt": {
           "href": "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.txt",
           "private": "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.txt",
           "public": "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.txt",
           "size": "530",
           "skipped": true
         }
       },
       "links": {
         "filter_blueprint": {
           "href": "test-href",
           "id": "test-id"
         },
         "self": {
           "href": "test-self-href",
           "id": "test-self-id"
         },
         "version": {
           "href": "test",
           "id": "test"
         }
       }
     },
     {
       "id": "94310d8d-72d6-492a-bc30-27584627edb2",
       "filter_id": "94310d8d-72d6-492a-bc30-27584627fil2",
       "instance_id": "94310d8d-72d6-492a-bc30-27584627inst2",
       "state": "completed",
       "downloads": {
         "xls": {
           "href": "string",
           "size": "string",
           "public": "string",
           "private": "string",
           "skipped": true
         },
         "csv": {
           "href": "string",
           "size": "string",
           "public": "string",
           "private": "string",
           "skipped": true
         },
         "csvw": {
           "href": "string",
           "size": "string",
           "public": "string",
           "private": "string",
           "skipped": true
         },
         "txt": {
           "href": "string",
           "size": "string",
           "public": "string",
           "private": "string",
           "skipped": true
         }
       }
  }]
    """

  Scenario: Get filter Output successfully
    When I GET "/filter-outputs/94310d8d-72d6-492a-bc30-27584627edb1"

    Then I should receive the following JSON response:
    """
        {
        "id": "94310d8d-72d6-492a-bc30-27584627edb1",
        "filter_id": "94310d8d-72d6-492a-bc30-27584627fil1",
        "instance_id": "94310d8d-72d6-492a-bc30-27584627inst1",
        "type": "",
        "downloads": {},
        "dimensions": null,
          "state": "published",
          "events": null,
          "published": true,
          "dataset": {
                  "id": "",
                  "edition": "",
                  "version": 0
           },
           "population_type": "",
          "downloads": {
            "xls": {
              "href": "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.xls",
              "private": "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.xls",
              "size": "6944",
              "public": "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.xls",
              "skipped": true
            },
            "csv": {
              "href": "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.csv",
              "private": "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.csv",
              "public": "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.csv",
              "size": "277",
              "skipped": true
            },
            "csvw": {
              "href": "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.csv-metadata.json",
              "private": "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.csvw",
              "public": "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.csvw",
              "size": "607",
              "skipped": true
            },
            "txt": {
              "href": "http://localhost:23600/downloads/datasets/cantabular-flexible-example/editions/2021/versions/1.txt",
              "private": "http://minio:9000/private-bucket/datasets/cantabular-flexible-example-2021-1.txt",
              "public": "https://csv-exported.s3.eu-west-1.amazonaws.com/full-datasets/cpih01-time-series-v1.csv-metadata.txt",
              "size": "530",
              "skipped": true
            }
          },
          "links": {
            "filter_blueprint": {
              "href": "test-href",
              "id": "test-id"
            },
            "self": {
              "href": "test-self-href",
              "id": "test-self-id"
            },
            "version": {
              "href": "test",
              "id": "test"
            }
          }
        }
    """
    And the HTTP status code should be "200"

  Scenario: Filter Output not found
    When I GET "/filter-outputs/not-exist"

    Then I should receive the following JSON response:
    """
    {
      "errors": ["failed to get filter output"]
    }
    """

    And the HTTP status code should be "404"

  Scenario: Mongo Connection is failing
    Given Mongo datastore is failing
    When I GET "/filter-outputs/not-exist"
    Then I should receive the following JSON response:
    """
    {
      "errors": ["failed to get filter output"]
    }
    """
    And the HTTP status code should be "500"
