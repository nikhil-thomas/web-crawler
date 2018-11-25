# Web Crawler
An application to generate site map od a given domain

## Design
![alt text](/screenshots/web-crawler-design.png "web crawler design")

## Usage

### Simple Usage

 ```
    docker run --rm nikhilvep/webcrawl:0.1
 ```
#### Options
![alt text](/screenshots/web-crawler-options.png "commandline flags")
 ```
    docker run --rm nikhilvep/webcrawl:0.1 https://github.com
 ```
#### Sample Output
![alt text](/screenshots/web-crawler-sample-output.png "sample sitemap")

## Build docker image

### Build
```
docker build -t web-crawler:0.1 .
```
### List options
```
docker run --rm web-crawler:0.1
```

### Run with default params
```
docker run --rm web-crawler:0.1 https://github.com
```
## Effect of params

url input: https://github.com

### case 1: default params
 * maximum page limit: 250
 * maximum links stored per page: 100
 * number of concurrent workers (url fetch and parse): 10
```
time docker run --rm nikhilvep/webcrawl:0.1 https://github.com
```
> 0.06s user 0.22s system 1% cpu 22.937 total

### case 2: running with 5 workers
 * maximum page limit: 250
 * maximum links stored per page: 100
 * number of concurrent workers (url fetch and parse): 5
```
time docker run --rm nikhilvep/webcrawl:0.1 -w 5 https://github.com
```
> 0.07s user 0.18s system 0% cpu 47.206 total

### case 2: running with concurrency disabled
 * maximum page limit: 250
 * maximum links stored per page: 100
 * concurrency disabled
```
time docker run --rm nikhilvep/webcrawl:0.1 -con-off https://github.com
```
> 0.07s user 0.20s system 0% cpu 2:29.47 total
