FROM node:20-slim AS base

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

RUN corepack enable

WORKDIR /app

COPY ./package*.json pnpm-lock.yaml ./

RUN pnpm install

COPY . ./

RUN pnpm run build

CMD [ "pnpm", "run", "dev" ]