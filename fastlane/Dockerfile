FROM thyrlian/android-sdk@sha256:bb9ed3686968550d927228777bca787dd7913e679f1e73e85525ba0094ea170d

COPY ./Gemfile .

RUN apt update && \
    apt install -y ruby-dev build-essential curl && \
    gem install bundler && \
    bundle install

ENV LC_ALL=en_US.UTF-8
ENV LANG=en_US.UTF-8
