FROM centos:7

ENV WORK_DIR=/opt/app/kvservice

RUN mkdir -p $WORK_DIR/
ADD kvservice $WORK_DIR/

RUN useradd -ms /bin/bash kvservice
USER kvservice

CMD ["/opt/app/kvservice/kvservice", "-p", "8282"]
