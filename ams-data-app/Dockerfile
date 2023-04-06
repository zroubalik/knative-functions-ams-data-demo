FROM node:14-alpine

WORKDIR /app

# Copy the package.json and package-lock.json files to the working directory
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code to the working directory
COPY . .

RUN npm run build

# Expose port 3000 for the React app and port 3333 for the backend
EXPOSE 3000
EXPOSE 3333

# Start the application
CMD ["npm", "start"]
