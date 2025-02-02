FROM node:20-alpine AS builder
ENV NODE_ENV=development
WORKDIR /app
COPY package*.json ./
RUN yarn install
COPY . .
RUN yarn build
ENV NODE_ENV=production
RUN yarn install --production

FROM node:20-alpine
ENV NODE_ENV=production
WORKDIR /app
COPY --from=builder /app/build build/
COPY --from=builder /app/node_modules node_modules/
COPY package.json .
EXPOSE 3000
ENV NODE_ENV=production
CMD [ "node", "build" ]