#!/bin/bash

# Initialize the MiniPaaS environment using provided Compose files
minipaas code init dev  --namespace sample

# Expose the `example` service on `localhost`
minipaas code route dev / example:8080

# Define a job service that runs once and exits after migration
minipaas code job dev example-migration

# Define worker services that consume messages in the background
minipaas code worker dev example-worker example-consumer

# Define a cron-based service for periodic execution
minipaas code cron dev example-cron

# Create and register a hashed Docker secret for Postgres password
echo postgres | minipaas deploy secret dev --verbose --name postgres_password --for postgres --for example --for example-migration --for example-consumer --for example-worker --for example-cron

# Build the project images using Docker Compose and tag them
minipaas deploy build dev --verbose

# Deploy the stack to Docker Swarm using rollout strategy
minipaas deploy rollout dev --verbose

# Apply routing configuration (Caddy update)
minipaas deploy routing dev --verbose

# Run hurl tests
hurl --verbose api.http
