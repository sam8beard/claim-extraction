# REQUIREMENTS AND SETUP

## mc - https://github.com/minio/mc

 >MinIO CLI interface

### Setup

Download latest binary
```wget https://dl.min.io/client/mc/release/linux-amd64/mc```


Give execute permissions 
``` chmod +x mc ```


Move binary into `/usr/local/bin/`
```mv mc /usr/local/bin/``` 


Add an administrative user with the username and password that is set in the `docker-compose.yml` and link it to the IP and port where you service is hosted
>Default user is `admin` and default password is `password`

```mc alias set minio-admin http://127.0.0.1:9000 admin password ```

>NOTE: Port `9000` is for the service and port `9001` is for the GUI

Verify connection to the MinIO server with your new user
```mc admin info minio-admin```

# FINISH FOR ALL SERVICES 
