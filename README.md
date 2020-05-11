#Go MySQL driver Timestamp/Timezone validator

This repo was created with the idea to validate an improvement
to do-mysql driver, since we found some problems working with
`TIMESTAMP` fields in our project.

## Problem

We have found that the time difference found when validating a time difference
between the moment and action was taken from the user and now, and the data from 
the action was stored in the database.

When we were running locally, everything was working wonderfull but, after deploy
to staging there was an unexpected error. The systeam was telling that the difference
was 3 hours more.

After we analyzed the problem we have found that:

- Driver uses, when not set, location from UTC
- Driver does not read the TimeZone from the global or session from DB
- MySQL timestamp store data according to Timezone, based on session
- Combined `loc` argument and `timezone` argument from the connection string produced some different behaviors.

### Behavior

    ATTENTION: Database is en Europe/Sofia and images are considering America/Sao_Paulo
    Considering: both loc and db variable
                    Localtime: 2020-05-11 18:29:54 -0300 -03,
                    DB Time: 2020-05-11 18:29:54 -0300 -03,
                    Difference Local to DB: 0s

    Considering: only loc
                    Localtime: 2020-05-11 18:29:54 -0300 -03,
                    DB Time: 2020-05-12 01:29:54 -0300 -03,
                    Difference Local to DB: -7h0m0s

    Considering: only DB variable
                    Localtime: 2020-05-11 21:29:54 +0000 UTC,
                    DB Time: 2020-05-11 18:29:54 +0000 UTC,
                    Difference Local to DB: 3h0m0s

    Considering: using none configuration
                    Localtime: 2020-05-11 21:29:54 +0000 UTC,
                    DB Time: 2020-05-12 01:29:54 +0000 UTC,
                    Difference Local to DB: -4h0m0s

### How to reproduce

I've built a simple app that exposes a post endpoint that allows you to post data to DB and read automatically.

There is 4 instances in the `docker-compose.yml` file and a MySQL database in this repo.

Source code from the program is in `main.go`

To run, you just need to run `run.sh` (and problably give execution permissions to it before running `chmod +x run.sh`)

