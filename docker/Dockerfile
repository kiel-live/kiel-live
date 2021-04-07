FROM node:12.22.1-alpine as build-stage
ARG VERSION
ENV VUE_APP_VERSION=${VERSION}
WORKDIR /app
COPY packages/app/package.json ./packages/app/
COPY package.json yarn.lock ./
RUN yarn install --silent --frozen-lockfile
COPY packages/app/ ./packages/app/
RUN yarn run build

#####################################################

FROM node:12.22.1-alpine
ARG VERSION

WORKDIR /app
RUN chown -R node:node /app
COPY packages/backend/package.json yarn.lock ./
USER node
RUN yarn install --silent --frozen-lockfile --production
COPY --chown=node:node docker/set-env-vars.sh docker/entrypoint.sh packages/backend/ ./
COPY --chown=node:node --from=build-stage /app/packages/app/dist/ dist/
EXPOSE 8080
ENTRYPOINT [ "./entrypoint.sh" ]
CMD ["yarn", "run", "start"]
ENV VERSION=${VERSION}
