FROM node:23.10.0

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

# Copy the rest of the source code
COPY . .

# Expose the port that Vite uses
EXPOSE 5173

# Run Vite with host binding to make it accessible from outside container
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]