FROM node:10-alpine

RUN mkdir -p /app/node_modules && chown -R node:node /app

WORKDIR /app

COPY package*.json ./

USER node

RUN npm install --silent

COPY --chown=node:node . .

EXPOSE 8080

ENTRYPOINT ["/app/start.sh"]
