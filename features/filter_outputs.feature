Feature: Filter Outputs Private Endpoints not Enabled
  Background:
    Given private endpoints are not enabled
    
  Scenario: Creating a new filter outputs journey when not unauthorized

    When I POST "/filter-outputs"
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
        "caller not found"
      ]
    }
    """

    And the HTTP status code should be "404"

