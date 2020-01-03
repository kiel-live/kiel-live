FROM node:12.13-alpine as build-stage
WORKDIR /app
COPY spa/package.json spa/yarn.lock ./
RUN yarn install --silent
COPY spa/ .
RUN yarn run build
ARG VERSION
ENV VUE_APP_VERSION=$VERSION

#####################################################

FROM node:12.13-alpine
WORKDIR /app
RUN chown -R node:node /app
COPY backend/package.json backend/yarn.lock ./
USER node
RUN yarn install --silent --frozen-lockfile --production
COPY --chown=node:node docker/set-env-vars.sh docker/entrypoint.sh backend/ ./
COPY --chown=node:node --from=build-stage /app/dist/ dist/
EXPOSE 8080
ENTRYPOINT [ "./entrypoint.sh" ]
CMD ["yarn", "run", "start"]
