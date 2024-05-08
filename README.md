# Utility `devenvctl`

`devenvctl.sh` allows me to quickly build, start, stop and switch development environments using [Docker](https://www.docker.com/). 

<br>

## Example Use Case:

>   - Project 1 requires MongoDB + ElasticSearch + RabbitMQ + Redis 5.0
>   - Project 2 requires Cassandra + Kafka + Redis 3.0 + Consul
>   - Project 3 requires same stack as Project 2, but requires different initial condition (e.g. different database schema)
>
> I can quickly switch between projects using this uitility: 
> 
> ```
> devenvctl.sh stop project_1 
> devenvctl.sh start project_2
> ```

<br>

## Usage

### Prerequisites 

1. `Docker for Mac` is required for this utility (2.0.0.0+ recommended). After install, verify that `docker-compose` is on `$PATH`.

    ```
    $ docker-compose --version
    docker-compose version 1.23.2, build 1110ad01
    ```
    
    ![recommended Docker for Mac version](res/devenvctl-f1.png "About Docker for Mac")

2. In `Docker for Mac`'s menu, select `Preferrences` -> `File Sharing`, and make sure `/usr/local/var` is in the list of allowed volumes to bind.

    ![configure bindable folders](res/devenvctl-f2.png "File Sharing")

3. Optionally, install GNU version of `readlink` (`greadlink`). e.g. 

    ```
    brew install coreutils
    ```

4. Optionally, create a symbolic link of `devenvctl.sh` in `$PATH` and restart terminal. 

    ```
    ln -s /absolute/path/to/devenvctl.sh /usr/local/bin/devenvctl 
    ```

### Common Usage 

- To see the help `devenvctl -h` or `devenvctl --help`

- Syntax: `devenvctl <action> <env_name>`, where actions are typically `info`, `prepare`, `start`, `stop` and `restart`. 
    
- E.g. to start predefiend "NFV" environment: 

  ```
  devenvctl start nfv
  ```

- To add your own environment definition, I'm not going into details. But you need to understand Shell scripting, Docker Compose and Dockerfile syntax. 
All environment definitions are under `devenv` folder. Each environment are composed by following components: 
    
    - script file with name `devenv-<your-env-name>.script`, which is a shell script handles start/stop/restart actions.
    
      > **Note**: If a script file with name `devenv-<your-env-name>.script` present, it will be used and the script may choose to 
      ignore any other files/folders described here. The script has full control of how an environment is prepared or teared down.
    
    - alternatively, an environment description `devenv-<your-env-name>.yml`, which describe required servies and versions, and is consumed by `default.script`.
      The most of values of this YAML file ar
      
      > **Note**: If a script file with name `devenv-<your-env-name>.script` present, `default.script` is not used.
      
    - docker compose yml with name `docker-compose-<your-env-name>.yml`, which is the docker compose descriptor.
    
    - folder with name `res-<your-env-name>`, which contains all extra files you may need to build your customized docker images.
    
    - Folder `/usr/local/var/dev/<your-env-name>` is used for data volumes per environment. Data is isolated between environments/projects. 
    
      By default, only database, consul and vault are configured to use persistent store. You can modify the docker compose file to alter this behavior. 

<br>


