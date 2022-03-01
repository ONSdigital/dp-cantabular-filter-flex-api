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
    """

    Then I should receive the following JSON response:
    """
    {
        "filter_output_id":"94310d8d-72d6-492a-bc30-27584627edb1",
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
                "href": "",
                "size": "string",
                "public": "string",
                "private": "string",
                "skipped": true
            },
            "csv":
            {
                "href": "string",
                "size": "string",
                "public": "",
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
    """

    Then I should receive the following JSON response:
    """
    {
      "errors": [
        "failed to validate request output: Public is empty in input"
      ]
    }
    """

    And the HTTP status code should be "400"