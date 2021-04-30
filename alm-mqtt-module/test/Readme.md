This is a manual test for testing registration, unregistration and timeout with the alm-mqtt-module.

Start environment:
```
git clone https://github.com/edgefarm/edgefarm-demos.git
cd edgefarm/edgefarm-demos/train-simulation/simulator
docker-compose -up -d
# enable simulation in UI
xdg-open http://localhost:1880/ui

docker run -p 4222:4222 -p 6222:6222 -p 8222:8222 --rm -d --name nats --network simulator_edgefarm-simulator nats
```

Build example and alm-mqtt-module & Start alm-mqtt-module:
```
cd alm-service-modules/alm-mqtt-module
make
make example
../bin/alm-mqtt-module
```

Start test:
```
cd alm-service-modules/alm-mqtt-module/test
./test.sh
```

Monitor the output of alm-mqtt-module
