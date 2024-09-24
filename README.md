# Golang Tech Test

## Introduction

This is a mock up of a "reactions" feature for the buffup platform.
This service allows users to react to a video stream by sending a reaction to the server.
This is then shown to all users watching the stream.

## The Test

Load has grown on our system and we need to scale the deployment of the service to multiple instances.
When we tried this, we noticed that all copies of the service try to process the full set of streams.
This means workload is not significantly reduced by adding more instances, and there are also issues with data consistency.

At the time this was reported, a dev has added a way to replicate this issue locally using `make run-multinode`

Like all start-ups sportbuff has grown quickly, and that has come at the cost of some technical debt.
Some services lack documentation and tests. This is one of those services.

You have a team member who was present when the service was built, but they did not write the code.
They may be able to answer some questions, but they are limited in their knowledge of the system.

Please update the code to allow for multiple instances of the service to run concurrently without data consistency issues.
