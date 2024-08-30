FROM mirror.gcr.io/node as build
WORKDIR /app

COPY package-lock.json package.json ./
RUN npm i

COPY ./public ./public
COPY ./src ./src
COPY ./styleguide ./styleguide 
COPY tsconfig.json styleguide.config.js webpack.config.js ./

RUN npm run build

FROM mirror.gcr.io/nginx
WORKDIR /root

COPY --from=build /app/build /static/
COPY ./nginx.conf /etc/nginx/nginx.conf