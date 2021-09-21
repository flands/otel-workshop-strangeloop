FROM ubuntu:20.04

RUN apt-get update -y && \
    apt-get install -y python3-pip python3-dev

# We copy just the requirements.txt first to leverage Docker cache
COPY ./requirements.txt /app/requirements.txt

WORKDIR /app

RUN pip install -r requirements.txt
# docker build will cache the following step as well if requirements.txt did not change
RUN opentelemetry-bootstrap -a=install

COPY . /app

CMD [ "opentelemetry-instrument", "python3", "app.py" ]
