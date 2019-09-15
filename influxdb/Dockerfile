FROM quay.io/influxdb/influxdb:2.0.0-alpha

ENV INFLUX_USERNAME, "admin"
ENV INFLUX_PASSWORD, "admin"
ENV INFLUX_ORG_NAME, "my-org"
ENV INFLUX_BUCKET_NAME, "my-bucket"

RUN "influx setup --username=INFLUX_USERNAME"
