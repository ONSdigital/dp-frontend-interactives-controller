Feature: Render embedded interactive

  As a Web User
  I am viewing a page with an iframe pointing to an embedded interactive

  Background:
    Given The browser loads the iframe

  Scenario: Viewing a valid embedded interactive
    Then I should see the "title" element