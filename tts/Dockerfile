FROM alpine

RUN apk update && apk add flite
ADD tts /tts
ENTRYPOINT ["/tts"]
