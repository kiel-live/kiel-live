FROM thyrlian/android-sdk

COPY ./Gemfile .

RUN apt update && \
    apt install -y ruby-dev build-essential curl && \
    gem install bundler && \
    bundle install

ENV LC_ALL=en_US.UTF-8
ENV LANG=en_US.UTF-8
