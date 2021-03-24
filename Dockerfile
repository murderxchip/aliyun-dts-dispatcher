FROM golang:1.13.14

ENV dtsEnv dev

# Install modules
#RUN apt-get update && apt-get install -y \
#        git \
#        curl \
#        wget \
#        telnet \
#        vim \
#             --no-install-recommends
RUN mkdir -p /home/logs
WORKDIR /home/work/dts
COPY ./bin /home/work/dts/
EXPOSE 9900
CMD /home/work/dts/dts-dispatcher-$dtsEnv serve

#docker run --env dtsEnv=staging -itd dtstest:1