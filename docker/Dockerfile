FROM alpine
RUN apk add bash
WORKDIR /home
COPY run-generator.sh /home/run-generator.sh
COPY generator /home/generator
ENTRYPOINT ["/bin/bash","-c","/home/run-generator.sh"]