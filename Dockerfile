FROM node:10-alpine

RUN mkdir -p /app/node_modules && chown -R node:node /app

WORKDIR /app

COPY package*.json ./

USER node

RUN npm install --silent

COPY --chown=node:node . .

RUN chmod +x /app/start.sh

EXPOSE 8080

CMD [ "/app/start.sh" ]
