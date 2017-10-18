
# Concourse Pipelines Dashboard
The [Concourse CI](http://concourse.ci) tool currently does not have a nice way to view all your pipelines.
Pipelines are hidden behind a drawer and the order depends on when you've added the pipeline (and can't be sorted)

This project solves this issue by using the Concourse REST API to create a single page dashboard of all your pipelines that with their current statuses.