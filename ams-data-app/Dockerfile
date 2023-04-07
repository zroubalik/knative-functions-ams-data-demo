# build environment
FROM node:16.20.0-alpine as build

WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH

# Copy package.json and package-lock.json
COPY frontend/package*.json ./

RUN npm ci --silent
RUN npm install react-scripts@ -g --silent

COPY frontend ./

RUN npm run build


# production environment
FROM nginxinc/nginx-unprivileged
USER nginx
COPY --from=build --chown=nginx:nginx /app/build /usr/share/nginx/html
EXPOSE 8080
CMD ["nginx", "-g", "daemon off;"]