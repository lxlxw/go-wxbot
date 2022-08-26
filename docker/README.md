### Run

```shell
docker run -d \
    --name="go-wxbot" \
    --restart=always \
    -v $(pwd)/config.yaml:/app/config.yaml \
    -v $(pwd)/imgs:/app/imgs \
    lxlxw/wxbot
```