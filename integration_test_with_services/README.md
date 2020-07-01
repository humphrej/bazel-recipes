Integration test with services
------------------------------
### Introduction
This recipe packages up multiple services inside a container in order to run integration tests.

**Integration testing with services** requires that the test executes code in both the client 
and a server.  To ensure hermeticity, both client and server should be stood up on demand, and deleted
afterwards.

A common example is testing the data access component of an application process that queries a 
database.  The test should create the database, run a schema migration (using for example liquibase),
run a test workload and assert on the results.

### Context:
* See link for the importance of the test pyramid:
https://testing.googleblog.com/2015/04/just-say-no-to-more-end-to-end-tests.html
* See this video for a more advanced solution by dropbox:
https://www.youtube.com/watch?v=muvU1DYrY0w
* Other approaches: 
    * Integration tests that span multiple containers are harder to orchestrate, and harder to 
    distribute to a build farm.
    * Some approaches like testcontainers are java only
    https://www.testcontainers.org/
    * It is sometimes possible to emulate the database in the application, for example using H2 in
     PostgreSQL mode. This has the disadvantage that the tests run on different infrastructure to 
     deployment.
    http://www.h2database.com/html/features.html#compatibility

### Solution: 
Create a single container that contains all services in order to reduce complexity of the test.  
This container specifically does _not_ follow the 12 factor app principles.  The integration test 
container can be constructed at compile time, just like a regular service container and deployed 
to a container registry.  The execution of the integration test can then happen asynchronously, and
results reported.

In this example, the entrypoint.sh script calls the prestart.d scripts in order before executing 
the test workload.

#### Solution Limitations:
Not suitable for anything other than test use, as the PostgreSQL data is stored on transiently within
the container.

#### Consequences:
* Positive
    * easy to test
    * very reproducible
    * easy to run on a build farm
* Negative
    * Still need to consider the boundaries of unit testing, and higher level/system/end-to-end 
    testing boundaries to supplement the integration test.
