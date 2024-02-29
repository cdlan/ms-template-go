FROM python:bookworm AS build-stage

WORKDIR /

RUN apt update && apt install -y plantuml
RUN pip install mkdocs
RUN cd / && mkdocs new ms-template-go
RUN cd /ms-template-go

# plugind
RUN pip install mkdocs-mermaid2-plugin
RUN pip install mkdocs_puml

WORKDIR /ms-template-go

# clean
RUN rm -f /ms-template-go/mkdocs.yml
RUN rm -rf /ms-template-go/docs

# Copy data
COPY ./mkdocs.yml /ms-template-go/

RUN mkdir docs
COPY ./docs /ms-template-go/docs

# build
RUN mkdocs build

#CMD ["sleep", "infinity"]

FROM nginx:stable-alpine-slim AS production-stage

RUN rm -rf /usr/share/nginx/html
RUN mkdir -p /usr/share/nginx/html

COPY --from=build-stage /ms-template-go/site /usr/share/nginx/html
