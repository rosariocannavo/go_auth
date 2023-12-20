# Use an existing Node.js LTS image as the base
FROM node:lts

# Set a working directory within the container
WORKDIR /usr/src/app

# Copy the package.json and package-lock.json files to the working directory
COPY package*.json ./

# Install Truffle globally
RUN npm install -g truffle

# Copy the entire Truffle project into the container
COPY . .
