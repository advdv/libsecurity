#   (c) Copyright 2014 Hewlett-Packard Development Company, L.P.
#   All rights reserved. This program and the accompanying materials
#   are made available under the terms of the Apache License v2.0 which accompany this distribution.
#
#   The Apache License is available at
#   http://www.apache.org/licenses/LICENSE-2.0
#
####################################################
# Replace a Docker container with latest version.
#
# Inputs:
#   - container_name - ID of the container to be deleted
#   - docker_host - Docker machine host
#   - docker_username - Docker machine username
#   - docker_password - optional - Docker machine password
#   - private_key_file - optional - path to private key file
# Outputs:
#   - error_message - error message of the operation that failed
####################################################

namespace: io.cloudslang.docker.containers

imports:
 cmd: io.cloudslang.base.cmd

flow:
  name: clear_container
  inputs:
    - container_name
  workflow:
    - stop_container:
        do:
          cmd.run_command:
            - container_name
            - command: "'docker stop' + container_name"
        publish:
          - error_message

    - delete_container:
        do:
          cmd.run_command:
            - container_name
            - command: "'docker rm' + container_name"
        publish:
          - error_message

    - remove_old_image:
        do:
          cmd.run_command:
            - command: "'docker rmi jerbi/shellshock'"
        publish:
          - error_message

    - pull_new_image:
        do:
          cmd.run_command:
            - command: "'docker pull jerbi/shellshock'"
        publish:
          - error_message

    - run_new_image:
        do:
          cmd.run_command:
            - container_name
            - command: "'docker run -d -p 80:80 '+' --name ' + container_name + ' jerbi/shellshock:latest'"
        publish:
          - error_message
  outputs:
    - error_message
