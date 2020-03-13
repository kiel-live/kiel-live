FROM node:12.13-alpine as build-stage
ARG VERSION
ENV VUE_APP_VERSION=${VERSION}
WORKDIR /app
COPY spa/package.json ./spa/
COPY package.json yarn.lock ./
RUN yarn install --silent --frozen-lockfile
COPY spa/ ./spa/
RUN yarn run build

#####################################################

FROM node:12.13-alpine
ARG VERSION

WORKDIR /app
RUN chown -R node:node /app
COPY backend/package.json yarn.lock ./
USER node
RUN yarn install --silent --frozen-lockfile --production
COPY --chown=node:node docker/set-env-vars.sh docker/entrypoint.sh backend/ ./
COPY --chown=node:node --from=build-stage /app/spa/dist/ dist/
EXPOSE 8080
ENTRYPOINT [ "./entrypoint.sh" ]
CMD ["yarn", "run", "start"]
ENV VERSION=${VERSION}
