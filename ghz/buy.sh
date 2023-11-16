ghz -n 1000000 --rps  500 \
    --concurrency-schedule=step --concurrency-start=5 \
    --concurrency-step=5 --concurrency-end=200 \
    --concurrency-step-duration=10s \
    --insecure --proto ../protos/processOrder/processOrder.proto \
    --call processOrder.OrderService.CreateOrder \
    -D createbuy.json \
    -o report/buy_2_details.txt --format influx-details \
    localhost:5432



# 压测命令  发送买单请求
