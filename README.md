# Feedmaker
Service allows to fetch data from various SQL databases, convert it to CSV, upload CSV files to FTP. Those files could be used to create ads in Google Ads automatically.
## Configuration
You can configure application in [config](/infrastructure/config/config.yml) file.
![config](/images/config.png?raw=true)
First of all, you have to create **.env** file fill variables with credentials for FTP, Redis(optional) and SQL for various generation types.
To add more generation types you can add a new key under **feeds** key. You should enter size limit, line limit of a single file, as well as SQL driver and connection string and paths to SQL queries.
## Running
```docker-compose up```
## API
All communications with this service made by API.
```
/list GET completed/active generations
/types GET list generation types
/types/{generation-type} POST start generation of feeds
/id/{generation-id} DELETE cancel generation
/id/{generation-id} POST restart generation
/ws/progress WS stream progress of active generations
/schedules GET list scheduled generations
/types/{generation-type}/schedules POST schedule generation
/types/{generation-type}/schedules DELETE unshedule generation
```
## UI
There is an SPA in React.js, but this interface is a part of bigger CRM system, so only screenshots could be attached:

![Generations](/images/generations.png?raw=true)
![Generations](/images/schedules.png?raw=true)
## Architecture
Project was designed using [Clean Architecture](http://cleancoder.com/files/cleanArchitectureCourse.md) and TDD principles. Fetching from database, converting data to CSV, uploading of files are done asynchronously with goroutines.
