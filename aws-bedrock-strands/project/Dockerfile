FROM ghcr.io/astral-sh/uv:python3.12-bookworm-slim
WORKDIR /app

ENV UV_SYSTEM_PYTHON=1 UV_COMPILE_BYTECODE=1

COPY requirements.txt requirements.txt
RUN uv pip install -r requirements.txt

RUN uv pip install aws-opentelemetry-distro>=0.10.1


# Set AWS region environment variable
ENV AWS_REGION=us-west-2
ENV AWS_DEFAULT_REGION=us-west-2


# Signal that this is running in Docker for host binding logic
ENV DOCKER_CONTAINER=1

# Create non-root user
RUN useradd -m -u 1000 bedrock_agentcore
USER bedrock_agentcore

EXPOSE 8080
EXPOSE 8000

# Copy entire project (respecting .dockerignore)
COPY . .

CMD ["opentelemetry-instrument", "python", "-m", "strands_claude"]
