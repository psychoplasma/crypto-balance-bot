FROM node:20-alpine AS deps

# Set the working directory
WORKDIR /app

# Install build dependencies for bcrypt
RUN apk add --no-cache python3 make g++ build-base

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
FROM node:20-alpine

# Set the working directory
WORKDIR /app

# Copy files to run application from builder stage
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/prisma ./prisma
COPY --from=builder /app/dist ./dist

# Expose the application port
EXPOSE 3000

ENTRYPOINT [ "/bin/sh", "-c" ]
