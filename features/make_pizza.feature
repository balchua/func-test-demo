Feature: Making a hawaiian pizza
  I want to make a pizza

  Scenario: Make a perfect pizza
    Given the following thresholds
      | min | max | ingredient_type | crust size |
      | 10 | 30   | H | 12|
      | 10 | 30   | P | 12|
      | 0.5 | 1.0   | T | 12|
      | 5 | 15   | H | 10|
      | 10 | 15   | P | 10|
      | 0.25 | 0.55   | T | 10|
    When the crust size is 12 inches
    And the ingredients "<tomato>", "<pineapple>", "<hams>"
    Then it should be a "<status>" pizza
    Examples:
      | crust_size | tomato | pineapple | hams | status |
      | 12 | 0.25   | 5 | 20 | not perfect|
      | 10 | 0.50   | 10 | 10 | perfect|
      | 18 | 0.50   | 10 | 80 | not perfect|