# MinIO Client Setup

[https://github.com/minio/mc]
 >MinIO CLI interface

Download latest binary

```bash
wget https://dl.min.io/client/mc/release/linux-amd64/mc
```


Give execute permissions 

```bash
chmod +x mc 
```


Move binary into `/usr/local/bin/`

```bash
mv mc /usr/local/bin/
```


Add an admin with the username and password that is set in `docker-compose.yml` and link it to the IP and port where the service is hosted
>Default user is `admin` and default password is `password`

```bash
mc alias set minio-admin http://localhost:9000 admin password 
```

>NOTE: Port `9000` is for the service and port `9001` is for the GUI

Verify connection to the MinIO server with your new user

```bash
mc admin info minio-admin
```

Add user for database operations
```bash
mc admin user add minio-admin username password
```
>This can be any username or password you want 

Verify new user has been added 
```bash
mc admin user ls minio-admin
```

Attach policies needed for database operations to new user
```bash
mc admin policy attach minio-admin readwrite --user username
```

Connect to MinIO service using newly created user 
```bash
mc alias set s3 http://localhost:9000 username password
```

Make bucket for claimex
```bash
mc mb s3/claim-pipeline-docstore
```
