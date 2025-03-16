FROM node:20 AS deps

# Set the working directory
WORKDIR /app

# Copy yarn.lock, package.json and package-lock.json
COPY package*.json yarn.lock nest-cli.json ./

# Install dependencies
RUN yarn install --frozen-lockfile

# Building stage
FROM node:20 AS builder

# Set the working directory
WORKDIR /app

# Copy dependencies from previous stage
COPY --from=deps /app/node_modules ./node_modules

# Copy the rest of the application code
COPY . .

# Generate prisma client
RUN npx prisma generate

# Build the application
RUN yarn build

# Runner stage
FROM node:20

# Set the working directory
WORKDIR /app

# Copy files to run application from builder stage
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/prisma ./prisma
COPY --from=builder /app/dist ./dist

# Expose the application port
EXPOSE 3000

CMD []
