nohup ghz -c 10 -n 60000 --rps 75 --connections 2 \
--insecure --proto ../protos/processOrder/processOrder.proto --call processOrder.OrderService.CreateOrder \
-D trade0.json -m '{"symbol": "s0"}' localhost:32801 >log0.out 0<&1 2>&1 &

nohup ghz -c 10 -n 60000 --rps 75 --connections 2 \
--insecure --proto ../protos/processOrder/processOrder.proto --call processOrder.OrderService.CreateOrder \
-D trade1.json -m '{"symbol": "s1"}' localhost:32801 >log1.out 0<&1 2>&1 &

nohup ghz -c 10 -n 60000 --rps 75 --connections 2 \
--insecure --proto ../protos/processOrder/processOrder.proto --call processOrder.OrderService.CreateOrder \
-D trade2.json -m '{"symbol": "s2"}' localhost:32801 >log2.out 0<&1 2>&1 &

nohup ghz -c 10 -n 60000 --rps 75 --connections 2 \
--insecure --proto ../protos/processOrder/processOrder.proto --call processOrder.OrderService.CreateOrder \
-D trade3.json -m '{"symbol": "s3"}' localhost:32801 >log3.out 0<&1 2>&1 &

# 对四个交易对下单
