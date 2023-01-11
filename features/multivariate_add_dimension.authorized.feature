Feature: Multivariate Feature Dimensions Private Endpoints
  Background:
    Given private endpoints are enabled
    And I am identified as "user@ons.gov.uk"

    And I am authorised

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
            "name": "geography"
          },
          {
            "label": "Number of siblings (3 mappings)",
            "links": {
              "code_list": {},
              "options": {},
              "version": {}
            },
            "href": "http://api.localhost:23200/v1/code-lists/siblings_3",
            "id": "siblings_3",
            "name": "siblings"
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
        "version": 1,
        "is_based_on": {
          "@id":"cantabular-example-1",
          "@type":"cantabular_multivariate_table"
        }
      }
      """
    And I have these filters:
    """
    [
      {
        "filter_id": "94310d8d-72d6-492a-bc30-27584627edb1",
        "links": {
          "version": {
            "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1",
            "id": "1"
          },
          "self": {
            "href": ":27100/filters/94310d8d-72d6-492a-bc30-27584627edb1"
          }
        },
        "instance_id": "c733977d-a2ca-4596-9cb1-08a6e724858b",
        "dimensions": [
        {
          "name": "siblings",
          "id": "siblings_3",
          "label": "Number of siblings (3 mappings)",
          "options": [],
          "is_area_type": false
        },
        {
          "name": "geography",
          "id": "country",
          "label": "Country",
          "options": [],
          "is_area_type": true
        }
        ],
        "dataset": {
          "id": "cantabular-example-1",
          "edition": "2021",
          "version": 1
        },
        "published": true,
        "population_type": "dummy_data_households",
        "type": "cantabular_multivariate_table"
      }
    ]
    """

    And Cantabular returns these dimensions for the dataset "dummy_data_households" and search term "hh_carers":
    """
    {
      "dataset": {
        "variables": {
          "edges": [
            {
              "node": {
                "categories": {
                  "totalCount": 32
                },
                "label": "Number of unpaid carers in household (32 categories)",
                "mapFrom": [],
                "name": "hh_carers"
              }
            }
          ]
        }
      }
    }
    """

 And Cantabular returns these categorisations for the dataset "dummy_data_households" and search term "hh_carers":
    """
    {
      "dataset": {
        "variables": {
          "edges": [
            {
              "node": {
                "isSourceOf": {
                  "edges": [
                    {
                      "node": {
                        "label": "Number of unpaid carers in household (6 categories)",
                        "name": "hh_carers_6a"
                      }
                    },
                    {
                      "node": {
                        "label": "Number of unpaid carers in household (32 categories)",
                        "name": "hh_carers_D"
                      }
                    }
                  ]
                },
                "mapFrom": []
              }
            }
          ]
        }
      }
    }
    """

  And Metadata api returns this response for the dataset "dummy_data_households" and search term "hh_deprivation":
  """
  {
    "data": {
      "dataset": {
        "vars": [
          {
            "meta": {
              "Default_Classification_Flag": "Y"
            },
            "name": "hh_carers_D"
          }
        ]
      }
    }
  }
   """

  Scenario: Adding a second dimension to a multivariate filter (mapFrom)
    Given I use an X Florence user token "user token"

    And I am identified as "user@ons.gov.uk"

    And zebedee recognises the user token as valid

    When I POST "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
    """
    {
      "name": "hh_carers",
      "is_area_type": false,
      "filter_by_parent": ""
    }
    """

    Then the HTTP status code should be "201"

    And I should receive the following JSON response:
    """
    {
    "name": "hh_carers_D",
    "id": "hh_carers_D",
    "default_categorisaion": "hh_carers_D",
    "label": "Number of unpaid carers in household (32 categories)",
    "links": {
      "filter": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1",
          "id": "94310d8d-72d6-492a-bc30-27584627edb1"
      },
      "options": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/hh_carers_D/options"
      },
      "self": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/hh_carers_D",
          "id": "hh_carers_D"
      }
      }
    }
    """


    And a document in collection "filters" with key "filter_id" value "94310d8d-72d6-492a-bc30-27584627edb1" should match:
    """
    {
      "_id": "94310d8d-72d6-492a-bc30-27584627edb1",
      "filter_id": "94310d8d-72d6-492a-bc30-27584627edb1",
      "links": {
        "version": {
          "href": "http://localhost:22000/datasets/cantabular-example-1/editions/2021/versions/1",
          "id": "1"
        },
        "self": {
          "href": ":27100/filters/94310d8d-72d6-492a-bc30-27584627edb1"
        },
        "dimensions": {
          "href": ""
        }
      },
      "etag": "",
      "instance_id": "c733977d-a2ca-4596-9cb1-08a6e724858b",
      "dataset": {
        "id": "cantabular-example-1",
        "edition": "2021",
        "lowest_geography": "",
        "release_date":  "",
        "title":"",
        "version": {
          "$numberInt":"1"
        }
      },
      "dimensions": [
        {
          "name": "siblings",
          "id": "siblings_3",
          "label": "Number of siblings (3 mappings)",
          "options": [],
          "is_area_type": false
        },
        {
          "name": "geography",
          "id": "country",
          "label": "Country",
          "options": [],
          "is_area_type": true
        },
        {
          "name": "hh_carers_D",
          "id": "hh_carers_D",
          "label": "Number of unpaid carers in household (32 categories)",
          "options":[],
          "is_area_type": false,
          "default_categorisation": "hh_carers_D"
        }
      ],
      "population_type": "dummy_data_households",
      "published": true,
      "type": "cantabular_multivariate_table",
      "unique_timestamp":{
        "$timestamp":{
          "i":0,
          "t":0
        }
      },
      "last_updated":{
        "$date":{
          "$numberLong":"-62135596800000"
        }
      }
    }
    """

  Scenario: Adding a second dimension to a multivariate filter sourceOf
    Given I use an X Florence user token "user token"

    And I am identified as "user@ons.gov.uk"

    And zebedee recognises the user token as valid

    And Cantabular returns these categorisations for the dataset "dummy_data_households" and search term "hh_carers":
    """
    {
      "dataset": {
        "variables": {
          "edges": [
            {
              "node": {
                "isSourceOf": {
                  "edges": [
                    {
                      "node": {
                        "label": "Number of unpaid carers in household (6 categories)",
                        "name": "SHOULD NOT USE"
                      }
                    },
                    {
                      "node": {
                        "label": "Number of unpaid carers in household (32 categories)",
                        "name": "SHOULT NOT USE"
                      }
                    }
                  ]
                },
                "mapFrom": [
                  {
                    "edges": [
                      {
                        "node": {
                          "label": "Number of unpaid carers in household (6 categories)",
                          "name": "SHOULD NOT USE",

                        "isSourceOf": {
                          "edges": [
                            {
                              "node": {
                                "label": "hh_carers_D",
                                "name": "hh_carers_D"
                              }
                            },
                            {
                              "node": {
                                "label": "Age (3 categories)",
                                "name": "SHOULD NOT USE"
                              }
                            }
                          ]
                        }

                        }}
                    ]
                  }
                ]
              }
            }
          ]
        }
      }
    }
    """
    When I POST "/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions"
    """
    {
      "name": "hh_carers",
      "is_area_type": false,
      "filter_by_parent": ""
    }
    """

    And I should receive the following JSON response:
    """
    {
    "name": "hh_carers_D",
    "id": "hh_carers_D",
    "label": "hh_carers_D",
    "default_categorisaion": "hh_carers_D",
    "links": {
      "filter": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1",
          "id": "94310d8d-72d6-492a-bc30-27584627edb1"
      },
      "options": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/hh_carers_D/options"
      },
      "self": {
          "href": "http://localhost:22100/filters/94310d8d-72d6-492a-bc30-27584627edb1/dimensions/hh_carers_D",
          "id": "hh_carers_D"
      }
      }
    }
    """
    Then the HTTP status code should be "201"
