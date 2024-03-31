Feature: Version Command

  Background:
    Given I have "go" command installed
    And a build of sck
    Then the build should be present

  Scenario:
    Given a build of sck
    When I run `sck-int-test`
    Then the output should contain:
      """"
      usage: sck [<flags>] <command> [<args> ...]

      A command-line tool for modifying ssh config files.
      """"
