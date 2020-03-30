Feature: render

  Background:
    Given I have a kubernetes cluster with name "stevedore" and version "1.17.0"
    And I have following helm repos:
      | Name   | URL                                              |
      | stable | https://kubernetes-charts.storage.googleapis.com |
    And I refresh helm local cache

  Scenario: Render simple stevedore manifest
    Given I have to install "stable/gocd" into my cluster as "gocd"
