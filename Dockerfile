FROM registry.access.redhat.com/ubi9/go-toolset:latest as addon
ENV GOPATH=$APP_ROOT
COPY --chown=1001:0 . .
RUN make cmd
RUN wget https://github.com/migtools/crane/releases/download/v0.0.5/amd64-linux-crane-v0.0.5
RUN chmod +x amd64-linux-crane-v0.0.5

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest
USER root
RUN echo -e "[centos9]" \
 "\nname = centos9" \
 "\nbaseurl = http://mirror.stream.centos.org/9-stream/AppStream/\$basearch/os/" \
 "\nenabled = 1" \
 "\ngpgcheck = 0" > /etc/yum.repos.d/centos.repo
RUN microdnf -y install \
 openssh-clients \
 subversion \
 git \
 tar
USER 1001
ENV HOME=/addon ADDON=/addon
WORKDIR /addon
ARG GOPATH=/opt/app-root
COPY --from=addon $GOPATH/src/bin/addon /usr/bin
COPY --from=addon $GOPATH/src/amd64-linux-crane-v0.0.5 /usr/bin/crane
ENTRYPOINT ["/usr/bin/addon"]