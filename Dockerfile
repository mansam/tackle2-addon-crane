FROM registry.access.redhat.com/ubi9/go-toolset:latest as addon
ENV GOPATH=$APP_ROOT
COPY --chown=1001:0 . .
RUN make cmd
RUN wget https://github.com/migtools/crane/releases/download/v0.0.5/amd64-linux-crane-v0.0.5
RUN chmod +x amd64-linux-crane-v0.0.5
RUN wget https://github.com/eksctl-io/eksctl/releases/download/v0.176.0/eksctl_Linux_amd64.tar.gz
RUN tar xvf eksctl_Linux_amd64.tar.gz
RUN wget https://github.com/kubernetes-sigs/aws-iam-authenticator/releases/download/v0.6.14/aws-iam-authenticator_0.6.14_linux_amd64
RUN chmod +x aws-iam-authenticator_0.6.14_linux_amd64

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest
USER 1001
ENV HOME=/addon ADDON=/addon
WORKDIR /addon
ARG GOPATH=/opt/app-root
COPY --from=addon $GOPATH/src/bin/addon /usr/bin
COPY --from=addon $GOPATH/src/amd64-linux-crane-v0.0.5 /usr/bin/crane
COPY --from=addon $GOPATH/src/eksctl /usr/bin/eksctl
COPY --from=addon $GOPATH/src/aws-iam-authenticator_0.6.14_linux_amd64 /usr/bin/aws-iam-authenticator
ENTRYPOINT ["/usr/bin/addon"]