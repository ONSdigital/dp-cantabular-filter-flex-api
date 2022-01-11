Feature: Helloworld



  Scenario: Posting and checking a response
    When the service starts
    And I GET "/hello"
    Then I should receive the following JSON response with status "200":
      """
        {"message":"Hello, World!"}
      """