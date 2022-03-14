Feature: Filter Dimensions Private Endpoints Not Enabled

  Background:
    Given private endpoints are not enabled
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
            "name": "City"
          },
          {
            "label": "Number of siblings (3 mappings)",
            "links": {
              "code_list": {},
              "options": {},
              "version": {}
            },
            "href": "http://api.localhost:23200/v1/code-lists/siblings",
            "id": "siblings",
            "name": "Number of siblings (3 mappings)"
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
    And I have these filters:
    """
    [
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
            "name": "Number of siblings (3 mappings)",
            "options": [
              "0-3",
              "4-7",
              "7+"
            ],
            "dimension_url": "http://dimension.url/siblings",
            "is_area_type": false
          },
          {
            "name": "City",
            "options": [
              "Cardiff",
              "London",
              "Swansea"
            ],
            "dimension_url": "http://dimension.url/city",
            "is_area_type": true
          }
        ],
        "dataset": {
          "id": "cantabular-example-1",
          "edition": "2021",
          "version": 1
        },
        "published": true,
        "population_type": "Example",
        "type": "flexible"
      }
    ]
    """

  Scenario: Get filter dimensions successfully
    When I GET "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
    Then I should receive the following JSON response:
    """
    {
      "items": [
        {
          "name": "City",
          "options": [
            "Cardiff",
            "London",
            "Swansea"
          ],
          "dimension_url": "http://dimension.url/city",
          "is_area_type": true
        },
        {
          "name": "Number of siblings (3 mappings)",
          "options": [
            "0-3",
            "4-7",
            "7+"
          ],
          "dimension_url": "http://dimension.url/siblings",
          "is_area_type": false
        }
      ],
      "count": 2,
      "offset": 0,
      "limit": 20,
      "total_count": 2
    }
    """
    And the HTTP status code should be "200"

  Scenario: Add a filter dimension successfully
    When I POST "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
    """
    {
      "name": "Number of siblings (3 mappings)",
      "is_area_type": false,
      "options": [
        "4-7",
        "7+"
      ]
    }
    """
    Then I should receive the following JSON response:
    """
    {
      "name": "Number of siblings (3 mappings)",
      "is_area_type": false,
      "options": [
        "4-7",
        "7+"
      ]
    }
    """
    And the HTTP status code should be "201"
    And an ETag is returned

  Scenario: Get paginated filter dimensions successfully (limit)
    When I GET "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions?limit=1"
    Then I should receive the following JSON response:
    """
    {
      "items": [
        {
          "name": "City",
          "options": [
            "Cardiff",
            "London",
            "Swansea"
          ],
          "dimension_url": "http://dimension.url/city",
          "is_area_type": true
        }
      ],
      "count": 1,
      "offset": 0,
      "limit": 1,
      "total_count": 2
    }
    """
    And the HTTP status code should be "200"

  Scenario: Get paginated filter dimensions successfully (offset)
    When I GET "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions?offset=1"
    Then I should receive the following JSON response:
    """
    {
      "items": [
        {
          "name": "Number of siblings (3 mappings)",
          "options": [
            "0-3",
            "4-7",
            "7+"
          ],
          "dimension_url": "http://dimension.url/siblings",
          "is_area_type": false
        }
      ],
      "count": 1,
      "offset": 1,
      "limit": 20,
      "total_count": 2
    }
    """
    And the HTTP status code should be "200"

  Scenario: Get paginated filter dimensions successfully (0 limit)
    When I GET "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions?limit=0"
    Then I should receive the following JSON response:
    """
    {
      "items": [],
      "count": 0,
      "offset": 0,
      "limit": 0,
      "total_count": 2
    }
    """
    And the HTTP status code should be "200"

  Scenario: Get filter dimensions unsuccessfully
    When I GET "/filters/94310d8d-72d6-492a-03cb-27584627edb1/dimensions"
    Then I should receive the following JSON response:
    """
    {
      "errors": ["failed to get filter dimensions"]
    }
    """
    And the HTTP status code should be "404"

  Scenario: Get filter dimensions unsuccessfully (cannot parse offset)
    When I GET "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions?offset=dog"
    Then I should receive the following JSON response:
    """
    {
      "errors": ["invalid parameter offset"]
    }
    """
    And the HTTP status code should be "400"

  Scenario: Get filter dimensions unsuccessfully (cannot parse limit)
    When I GET "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions?limit=dog"
    Then I should receive the following JSON response:
    """
    {
      "errors": ["invalid parameter limit"]
    }
    """
    And the HTTP status code should be "400"

  Scenario: Get filter dimensions unsuccessfully (invalid limit, too large)
    Given the maximum pagination limit is set to 100
    When I GET "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions?limit=101"
    Then I should receive the following JSON response:
    """
    {
      "errors": ["limit cannot be larger than 100"]
    }
    """
    And the HTTP status code should be "400"

  Scenario: Get filter dimensions unsuccessfully (invalid limit, less than 0)
    When I GET "/filters/94310d8d-72d6-492a-03cb-27584627edb1/dimensions?limit=-1"
    Then I should receive the following JSON response:
    """
    {
      "errors": ["limit cannot be less than 0"]
    }
    """
    And the HTTP status code should be "400"

  Scenario: Get filter dimensions unsuccessfully (invalid offset, less than 0)
    When I GET "/filters/94310d8d-72d6-492a-03cb-27584627edb1/dimensions?offset=-1"
    Then I should receive the following JSON response:
    """
    {
      "errors": ["offset cannot be less than 0"]
    }
    """
    And the HTTP status code should be "400"

