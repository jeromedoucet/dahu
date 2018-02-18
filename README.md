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

to complete

## API endpoint

 - POST  /jobs create a new Job
 - GET   /jobs list all available jobs for the authenticated user
 - POST  /jobs/:jobId/run create a new run of a given job
 - GET   /jobs/:jobId get the details of a Job
 - PATCH /jobs/:jobId update a job
 - POST  /login authenticate a user

