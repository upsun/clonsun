ClonSun
=========

This tool makes a new copy of a project from a platform (Platform.sh, Upsun) to another platform (with some data) and convert at same-time for upsun.

> **WARNING : This tool does not clone the other services such as: ElasticSearch, mongoDB...**  

> **NOTE : This tool automatically converts to upsun, but the reverse is not supported.**

#### Syntax
```
Usage of clonsun:
      --src_provider string       Source provider CLI (default "platform")
      --src_project string        Source Project ID
      --src_env string            Source Environment (default "main")
      --dst_provider string       Destination provider CLI (default "upsun")
      --dst_project string        Destination Project ID
      --dst_env string            Destination Environment (default "main")
      --dst_organisation string   Destination Organisation
      --dst_region string         Destination Region (default "fr-4.platform.sh")
      --no_data                   Do not clone data (databases & mounts)
      --only_data                 Clone only data (databases & mounts)
      --no_users                  Do not clone user.
      --mount_type string         Change 'Local' mount to upsun compatible mode : 'storage' or 'instance'. (default "storage")
  -v, --verbose                   Enable verbose mode
```


#### CLI provider supported
| Provider | Style | Compatible | Tested | Note                |
|----------|-------|------------|--------|---------------------|
| platform | psh   | ✅         | ✅     | |
| upsun    | upsun | ✅         | ✅     | |
| ibexa    | psh   | ✅         | ⛔     | |
| shopware | psh   | ✅         | ⛔     | |
| pimcore-cloud  | psh   | ✅         | ⛔     | |
| webpaas  | psh   | ✅(partial)| ✅     | User is not apply |
| magento  | psh   | ⛔         | ⛔     | |


#### Samples of convertion

- Clone into new project (created by clonsun)  
`$ clonsun --src_provider=platform --src_project=xxxxxxxxxxxxx --src_env=master --dst_provider=upsun --dst_project="" --dst_env="" --dst_organisation=mick --dst_region=eu-3.platform.sh`
- Clone into existing project (and different environement)  
`$ clonsun --src_provider=platform --src_project=xxxxxxxxxxxxx --src_env=master --dst_provider=upsun --dst_project="yyyyyyyyyyyyy" --dst_env="main" --dst_organisation=mick --dst_region=eu-3.platform.sh`
- Clone into new project (created by clonsun) without copy of data.  
`$ clonsun --no_data --src_provider=platform --src_project=xxxxxxxxxxxxx --src_env=master --dst_provider=upsun --dst_project="" --dst_env="" --dst_organisation=mick --dst_region=eu-3.platform.sh`
- Copy only data.  
`$ clonsun --only_data --src_provider=platform --src_project=xxxxxxxxxxxxx --src_env=master --dst_provider=upsun --dst_project="" --dst_env="" --dst_organisation=mick --dst_region=eu-3.platform.sh`