FROM gradle:7.2.0-jdk11

COPY . /home/work
WORKDIR /home/work

RUN gradle distTar
RUN mkdir /app
RUN cp build/distributions/*.tar /app/
WORKDIR /app
RUN tar -xvf /app/*.tar

ENTRYPOINT ["otel-workshop-1.0-SNAPSHOT/bin/otel-workshop"]