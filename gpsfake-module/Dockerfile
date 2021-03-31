FROM python:2.7.18-alpine3.11

EXPOSE 2947 8888
COPY simulator /opt/nmea/simulator/
WORKDIR /opt/nmea/simulator
RUN apk add --update --no-cache gpsd && chmod +x /opt/nmea/simulator/start.sh

CMD ["./start.sh"]
