# Use a Node base image
FROM node:alpine

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application
COPY . .

# Build the application
RUN npm run build

# Install a simple server for serving static content
RUN npm install -g serve

# Serve the application on port 3000
CMD ["serve", "-s", "dist", "-p", "3000"]
