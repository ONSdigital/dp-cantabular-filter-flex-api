Feature: Helloworld

  Background:
    Given I GET "/hello"

  Scenario: Posting and checking a response
    When the service starts
    Then I should receive the following JSON response with status "200":
      """
        {"message":"Hello, World!"}
      """
