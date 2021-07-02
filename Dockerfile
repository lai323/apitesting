FROM ubuntu:18.04
COPY dist/apitesting /
COPY page/dist /page
# EXPOSE 8341
ENTRYPOINT ["/apitesting"]
