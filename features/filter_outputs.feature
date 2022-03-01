Feature: Filter Outputs Private Endpoints not Enabled
  Background:
    Given private endpoints are not enabled
    
  Scenario: Creating a new filter outputs journey when not unauthorized

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
      "errors": [
        "caller not found"
      ]
    }
    """

    And the HTTP status code should be "404"

