# gobackup

Gobackup backs up a directory or file to S3.

### Usage

The env variables `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` must be set.  A .env file can be used to set them.

The following flags are _required_

```
-s3region=<s3 region>
-s3bucket=<s3 bucket>
-source=</path/to/source/directory/or/file>
```

The following flags are _optional_

```
-prefix=<prefix of the s3 path, ie a prefix of darren would make the s3 path <s3bucket>/darren/<file>
-env=</path/to/envfile>
```

Example usage

```bash
/usr/local/bin/gobackup -s3region=us-west-2 \
-s3bucket=replace-me -source=/home/directory \
-env=/path/to/env/file \
-prefix=test

```

woohoo!!!!!
