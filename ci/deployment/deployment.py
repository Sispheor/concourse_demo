#!/usr/bin/env python

import docker

# connect to docker daemon
client = docker.DockerClient(base_url="tcp://10.33.101.57:2376")

# pull images from the private registry
hello_image_from_private_registry = "10.33.101.57:5000/hello"
client.images.pull(hello_image_from_private_registry)

# kill existing
hello_container = client.containers.get("hello")
if hello_container:
    try:
        hello_container.kill()
        print("Old instance killed")
    except docker.errors.NotFound:
        pass
    except docker.errors.APIError:
        pass

    try:
        hello_container.remove()
        print("Old instance removed")
    except docker.errors.APIError:
        pass
    except docker.errors.APIError:
        pass

# run
ports = {'8000/tcp': 8000}
client.containers.run(hello_image_from_private_registry, ports=ports, detach=True, name="hello")
print("New instance created")
