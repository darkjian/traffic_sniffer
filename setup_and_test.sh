#!/bin/sh

# Настройка TUN-интерфейса
ip addr add 10.1.0.10/24 dev tun0
ip link set tun0 up

# Отправка 5 UDP пакетов
echo "Sending 5 UDP packets to 10.1.0.20"
for i in $(seq 1 5); do
  echo "UDP packet $i" | nc -v -u -w 3 10.1.0.20 80
  sleep 1
done

# Отправка 5 TCP пакетов
echo "Sending 5 TCP packets to 10.1.0.20"
for i in $(seq 1 5); do
  echo "TCP packet $i" | nc -v -4 -w 3 10.1.0.20 85
  sleep 1
done


