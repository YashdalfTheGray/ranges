# ranges

A thing that yells about kitchen ranges on each port bound

## What is this?

Did you ever want range recommendations from GoodHousekeeping but delivered to you, in machine parsable format, ready for you to buy? Well, want no more. This will help.

Each port besides the first one, 8080, provides a new recommendation to a kitchen range in the format

```json
{
  "make": "Manufacturer",
  "model": "Range",
  "link": "buy link"
}
```

The first port, 8080, will return a status, that tells you whether things are okay and if we are ready to advertise our wonderful kitchen ranges.

## Usage

There are multiple ways to run this, they are documented below

1. `go run ./...` will just start up the server, easy peasy
1. `go build && ./ranges` will build the binary and start it up
1. `docker build -t ranges . && docker run -it -p 8080-8085:8080-8085 ranges` will build and run a docker image with the code and map ports 8080 through 8085 from this to ports 8080 through 8085 on the docker host
1. `docker run -d -p 8080-8085:8080-8085 public.ecr.aws/d8r6i5i8/ranges:latest` will pull it from ECR and run the container with the same port mappings as above.

You can also emit the host part of the port mappings to let Docker do dynamic port mapping by running the following command

`docker run -d -p 8080-8085 public.ecr.aws/d8r6i5i8/ranges:latest`
