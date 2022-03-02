Feature: Filter Outputs Private Endpoints Enabled
  Background:
    Given private endpoints are enabled
Scenario: Creating a new filter outputs journey when authorized
    Given I am identified as "user@ons.gov.uk"
    
    And I am authorised
    
    When I POST "/filter-outputs"
    """
    {
        "state": "string",
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
        "filter_output_id":"94310d8d-72d6-492a-bc30-27584627edb1",
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

    And the HTTP status code should be "201"

Scenario: Creating a new filter outputs bad request body
    Given I am identified as "user@ons.gov.uk"
    
    And I am authorised
    When I POST "/filter-outputs"
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
    
Scenario: Creating a new filter empty request body
    Given I am identified as "user@ons.gov.uk"
    
    And I am authorised
    When I POST "/filter-outputs"
    """
    {
        "state": "string",
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
              "public" : "",
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
        "failed to parse request: invalid request: \"Public\" is empty in input"
      ]
    }
    """

    And the HTTP status code should be "400"