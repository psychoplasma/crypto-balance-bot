FROM node:20

# Set the working directory
WORKDIR /usr/src/app

# Copy yarn.lock, package.json and package-lock.json
COPY package*.json yarn.lock nest-cli.json ./

# Install dependencies
RUN yarn install

# Copy the rest of the application code
COPY . .

# Generate prisma client
RUN npx prisma generate

# Build the application
RUN npm run build

# Expose the application port
EXPOSE 3000

# Command to run the application
CMD ["npm", "start"]