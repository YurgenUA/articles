FROM node:16-buster-slim
RUN mkdir -p /app
COPY ./*.js /app
COPY ./*.json /app
WORKDIR /app
RUN npm install
CMD npm run start