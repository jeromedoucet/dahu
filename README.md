# Dahu

[![Go Report Card](https://goreportcard.com/badge/github.com/jeromedoucet/dahu)](https://goreportcard.com/report/github.com/jeromedoucet/dahu)

Dahu is a CI / CD plateform build with an obsession : simplicity. Why ? 

 - Because it is more stable
 - Because it allow build complex system easier
 - Because it lead to have less surprise


## Features (in progress ...)

 - may pipeline CI / CD if necessary
 - simple configuration import / export feature
 - REST Api configuration
 - cli configuration
 - full containered build : all jobs are run inside a container.
 - full containered clone : all scm operations are run inside a container 
 - full containered deployment : all deployments operation are run inside a container 
 - full containered notification : all notification are triggered inside a container 
 - full containered washing : even my washing is done inside a container !
 - full containered ... etc, you have the picture.
 

## Choices

Simplicity means less features. Here are some choices explanations

### No 'CI as code'

This feature is very common on modern CI / CD plateform. It allow you to locate the configuration of a build within the code, in the repository.

However, such feature do have some drawbacks :

 - add a dependency on a specific CI plateform
 - often, not all configuration can be stored on repository (credentials for instance)
 - Doesn't fit very well when many repositories are involve
 - Most of the time, it is not usable locally
 - add a syntax to learn, making the tool less simple to use


