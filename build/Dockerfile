FROM quay.io/app-sre/boilerplate:image-v2.1.0 AS builder

ENV OPERATOR=/usr/local/bin/hive-health-operator \
    USER_UID=1001 \
    USER_NAME=hive-health-operator

# install operator binary
COPY build/_output/bin/hive-health-operator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
