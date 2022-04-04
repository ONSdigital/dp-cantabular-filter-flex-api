Feature: Get Filter Private Endpoints Not Enabled

  Background:
    Given private endpoints are not enabled

    And I have these filters:
    """
INPUT SOME FILTER OUTPUTS HERE
    """

  Scenario: Get filter Output successfully
    When I GET "/filters/ID HERE"

    Then I should receive the following JSON response:
    """
  OUTPUT GOES HERE
    """
    And the HTTP status code should be "200"

  Scenario: Filter Output not found
    When I GET "/filters/94310d8d-72d6-492a-03cb-27584627edb5"

    Then I should receive the following JSON response:
    """
    {
      "errors": ["failed to get filter"]
    }
    """

    And the HTTP status code should be "404"

  # Scenario: Unauthorized request on unpublished dataset
  #   Given I am not identified

  #   And I am not authorised

  #   When I GET "/filters/83210d8d-72d6-492a-bc30-27584627abc2"

  #   Then I should receive the following JSON response:
  #   """
  #   {
  #     "errors": ["failed to get filter"]
  #   }
  #   """

  #   And the HTTP status code should be "404"
