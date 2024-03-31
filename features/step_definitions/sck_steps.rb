require 'fileutils'

Given(/^I have "([^"]*)" command installed$/) do |command|
  is_present = system("which #{ command} > /dev/null 2>&1")
  raise "Command #{command} is not present in the system" if not is_present
end

Then('the build should be present') do
  steps %Q(
    Then a file named "#{$bin_dir}/sck-int-test" should exist
  )
end

Given("a build of sck") do
  raise 'sck build failed' unless system("go build -o bin/sck-int-test main.go")
end