# first stage: compile ts and build with ncc
FROM node:18-alpine as builder
WORKDIR /usr/app
COPY ./typescript .
RUN npm ci 
RUN npm run build

# second stage: run app
FROM node:18-alpine
WORKDIR /usr/app
COPY --from=builder /usr/app/ncc .
COPY --from=builder /usr/app/.env.production .
COPY --from=builder /usr/app/healthcheck.js .
RUN apk --no-cache add curl
ENV NODE_ENV=production
EXPOSE 3001
CMD ["node", "index.js"]


