### Quick start
> - start mongo service  
> - add `mongo uri` in `.env.default`
```bash
# go version => 1.11
> git clone git@github.com:BreezeFreak/self-help-api.git
> cd self-help-api/src/api
> go mod tidy
> go vendor

> go run
```
or

```bash
> docker-compose up -d
```