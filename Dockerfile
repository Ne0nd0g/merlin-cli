# Image source is located at https://github.com/Ne0nd0g/merlin-docker/blob/main/Dockerfile
# Image repository is at https://hub.docker.com/r/ne0nd0g/merlin-base
FROM ne0nd0g/merlin-base:v1.5.0

# Build the Docker image first
#  > sudo docker build -t merlin-cli .

# To start the Merlin Server and interact with it, run:
#  > sudo docker run -it --network host merlin-cli:latest

# The '--network host' argument allows the Merlin CLI to connect to the Merlin Server when running on the same host

WORKDIR /opt/merlin-cli

ENTRYPOINT ["./merlinCLI-Linux-x64", "-addr", "127.0.0.1:50051"]