FROM node:10-alpine

RUN mkdir -p /home/node/app/node_modules && mkdir -p /home/node/app/spa/node_modules && chown -R node:node /home/node/app

WORKDIR /home/node/app

COPY package*.json ./
COPY spa/package*.json ./spa/

USER node

RUN npm install
RUN cd spa/ && npm install

COPY --chown=node:node . .

RUN cd spa/ && npm run build

EXPOSE 8080

CMD [ "npm", "run", "start" ]
