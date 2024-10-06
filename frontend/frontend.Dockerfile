FROM node:20.17.0-alpine3.19 as base-frontend

WORKDIR /frontend

COPY package.json yarn.lock ./

RUN yarn install

COPY . .

EXPOSE 3000

CMD ["yarn", "run", "dev"]

FROM base-frontend as builder-frontend

RUN yarn build

FROM nginx:1.25.3-alpine3.18 as production-frontend

COPY --from=builder-frontend /frontend/build /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
