# Download and run

## Direct download
[pgtester](https://github.com/MannemSolutions/pgtester) is available for download for many platforms and architectures from the [Github Releases page](https://github.com/MannemSolutions/pgtester/releases).
It could be as simple as:
```bash
PGTESTER_VERSION=v0.3.0
cd $(mktemp -d)
curl -Lo "pgtester-${PGTESTER_VERSION}-linux-amd64.tar.gz" "https://github.com/MannemSolutions/pgtester/releases/download/${PGTESTER_VERSION}/pgtester-${PGTESTER_VERSION}-linux-amd64.tar.gz"
tar -xvf "./pgtester-${PGTESTER_VERSION}-linux-amd64.tar.gz"
mv pgtester /usr/local/bin
cd -
```
After that you can run pgtester directly from the prompt:
```bash
pgtester ./mytest1.yml mytest2.yml
```
Or using stdin:
```bash
cat ./mytests*.yml | pgtester
```

## Container image
For container environments [pgtester](https://github.com/MannemSolutions/pgtester) is also available on [dockerhub](https://hub.docker.com/repository/docker/mannemsolutions/pgtester).
You can easilly pull it with:
```bash
docker pull mannemsolutions/pgtester
```

Using it would be as easy as:
```bash
cat testdata/pgtester/tests.yaml | docker run -i mannemsolutions/pgtester pgtester
```

## docker-compose
You can use pgtester with docker compose.
The docker-compose.yml file could have contents like this:
```yaml
services:
  pgtester:
    image: mannemsolutions/pgtester
    command: pgtester /etc/pgtestdata/examples/tests1.yaml
  postgres:
    image: postgres:13
    environment:
      POSTGRES_HOST_AUTH_METHOD: 'md5'
      POSTGRES_PASSWORD: pgtester
```

it could be as easy as:
```bash
docker-compose up
```

or with only output for pgtester:
```bash
docker-compose up -d postgres
docker-compose up pgtester
```

or with tests defined locally:
```bash
docker-compose up -d postgres
cat ./mytests*.yml | docker-compose up pgtester
```

## Direct build

Although not advised, you can also directly build from source:
```bash
go install github.com/mannemsolutions/pgtester/cmd/pgtester@v0.3.0
```

After that you can run pgtester directly from the prompt:
```bash
pgtester ./mytest1.yml mytest2.yml
```

Or using stdin:
```bash
cat ./mytests*.yml | pgtester
```
