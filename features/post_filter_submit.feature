Feature: Post Filter Private Endpoints Not Enabled

  Background:
    Given private endpoints are not enabled

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
      },
      {
        "filter_id": "83210d8d-72d6-492a-bc30-27584627abc2",
        "links": {
          "version": {
            "href": "http://mockhost:9999/datasets/cantabular-example-unpublished/editions/2021/version/1",
            "id": "1"
          },
          "self": {
            "href": ":27100/filters/83210d8d-72d6-492a-bc30-27584627abc2"
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
          "id": "cantabular-example-unpublished",
          "edition": "2021",
          "version": 1
        },
        "published": false,
        "population_type": "Example",
        "type": "flexible"
      }
    ]
    """

  Scenario: POST filter successfully
    When I POST "/filters/94310d8d-72d6-492a-bc30-27584627edb1/submit"
    """
    """
    Then I should receive the following JSON response:
    """
        {

              "instance_id":"",
              "dimension_list_url":"",
              "filter_id":"94310d8d-72d6-492a-bc30-27584627edb1",
              "events":[{
                      "timestamp": "2016-07-17T08:38:25.316Z",
                      "name": "mock-export-event"
              }],

           "dataset":{
              "id":"mock-id",
              "edition":"mock-edition",
              "version":0
           },
           "links":{
              "version":{
                 "href":""
              },
              "self":{
                 "href":""
              }
           }
        }
    """
    And the HTTP status code should be "202"
