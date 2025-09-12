package builder

const BasicServerTemplate = `FROM alpine:latest

# Install necessary packages (openjdk for Minecraft, bash, curl, etc.)
RUN apk add --no-cache openjdk21-jre bash curl

# Set working directory
WORKDIR /server

# Copy the Minecraft server jar
COPY assets/executables/paper-1.21.8-58.jar /server/server.jar

# Expose default Minecraft port
EXPOSE 25565

# Accept EULA by default
RUN echo "eula=true" > eula.txt

# Default command to start the Minecraft server
CMD ["java", "-Xmx1024M", "-Xms1024M", "-jar", "server.jar", "nogui"]
`
