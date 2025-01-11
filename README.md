# BAK Backup System to backup files to a remote S3 bucket

## Configuration

First you need to set the configuration for the S3 bucket. You can do this by running the following command:  
`$ ./bak config set`

```shell
Enter S3 endpoint URL: https://your-s3-endpoint.com
Enter access key: AKIAXXXXXXXXXXXXXXXX
Enter secret key: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
Enter default bucket: your-bucket-name
```
The config file is located at `~/.bak/config.json`

## Upload files

To upload files to the S3 bucket, you can run the following command:  
`$ ./bak upload /path/to/file /remote/path/to/file`

You can encrypt the file before uploading it by using the `--encrypt` flag:  
`$ ./bak upload /path/to/file /remote/path/to/file --encrypt`

The base path for the remote path is the default bucket set in the configuration.

## Download files

To download files from the S3 bucket, you can run the following command:  
`$ ./bak download /remote/path/to/file /path/to/file`

It will automatically decrypt the file if it was encrypted before uploading.  
The base path for the remote path is the default bucket set in the configuration.

## List files

To list all files in the S3 bucket, you can run the following command:  
`$ ./bak ls`

The base path for the remote path is the default bucket set in the configuration.

## Delete files

To delete a file from the S3 bucket, you can run the following command:  
`$ ./bak delete /remote/path/to/file`

The base path for the remote path is the default bucket set in the configuration.


