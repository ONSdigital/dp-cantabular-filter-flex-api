Feature: Post Filter Output event Private Endpoints Enabled

  Background:
    Given private endpoints are enabled

    And I have these filter outputs:
    """
    [
      {
        "id": "94310d8d-72d6-492a-bc30-27584627edb1",
        "filter_id": "94310d8d-72d6-492a-bc30-27584627fil1",
        "instance_id": "94310d8d-72d6-492a-bc30-27584627inst1",
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
        },
        "events":[]
      },
      {
        "id": "94310d8d-72d6-492a-bc30-27584627edb2",
        "filter_id": "94310d8d-72d6-492a-bc30-27584627fil2",
        "instance_id": "94310d8d-72d6-492a-bc30-27584627inst2",
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
        },
        "events":[]
      }
    ]
    """

  Scenario: Post filter output event successfully with existing ID in DB
    Given I am identified as "user@ons.gov.uk"
    
    And I am authorised

    When I POST "/filter-outputs/94310d8d-72d6-492a-bc30-27584627edb1/events"
    """
    {
      "timestamp": "2016-07-17T08:38:25.319+0000",
      "name": "cantabular-export-start"
    }
    """

    Then the HTTP status code should be "201"
  
  Scenario: Post filter output event successfully with new ID in DB
    Given I am identified as "user@ons.gov.uk"
    
    And I am authorised

    When I POST "/filter-outputs/94310d8d-72d6-492a-bc30-27584627new1/events"
    """
    {
      "timestamp": "2016-07-17T08:38:25.319+0000",
      "name": "cantabular-export-start"
    }   
    """

    Then the HTTP status code should be "404"

  
  Scenario: Creating a new filter outputs bad request body
    Given I am identified as "user@ons.gov.uk"
    
    And I am authorised
    When I POST "/filter-outputs/94310d8d-72d6-492a-bc30-27584627edb1/events"
    """
    {
      "timestamp": 2,
      "name": "cantabular-export-start"
    }
    """

    Then I should receive the following JSON response:
    """
    {
      "errors": [
        "badly formed request body: json: cannot unmarshal number into Go struct field addFilterOutputEventRequest.timestamp of type string"
      ]
    }
    """

    And the HTTP status code should be "400"
    
  Scenario: Creating a new filter output with one wrong field in request body
    Given I am identified as "user@ons.gov.uk"
    
    And I am authorised
    When I POST "/filter-outputs/94310d8d-72d6-492a-bc30-27584627edb1/events"
    """
    {
      "timestamp": 2,
      "name": "cantabular-export-start"
    }  
    """

    Then the HTTP status code should be "400"
     
  Scenario: Post a filter output event with broken mongo db
    Given I am identified as "user@ons.gov.uk"
    
    And I am authorised

    And Mongo datastore fails
    When I POST "/filter-outputs/94310d8d-72d6-492a-bc30-27584627edb1/events"
    """
    {
      "timestamp": "2016-07-17T08:38:25.319+0000",
      "name": "cantabular-export-start"
    }     
    """

    Then the HTTP status code should be "500"
  

    