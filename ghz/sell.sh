ghz -n 1000000 --rps  500 \
    --concurrency-schedule=step --concurrency-start=5 \
    --concurrency-step=5 --concurrency-end=200 \
    --concurrency-step-duration=10s \
    --insecure --proto ../protos/processOrder/processOrder.proto \
    --call processOrder.OrderService.CreateOrder \
    -D createsell.json \
    -o report/sell_2_details.txt --format influx-details \
    localhost:5432

# 撮合压测命令  发送卖单请求
